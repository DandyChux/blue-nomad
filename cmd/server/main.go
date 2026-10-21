package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dandychux/blue-nomad/internal/handlers"
	"github.com/dandychux/blue-nomad/internal/middleware"
	"github.com/dandychux/blue-nomad/internal/services"

	"github.com/joho/godotenv"
)

// spaFileServer serves static files from the SvelteKit build output directory.
// If the requested file does not exist on disk, it falls back to serving
// index.html so that SvelteKit's client-side router can handle the route.
type spaFileServer struct {
	staticDir  string
	fileServer http.Handler
}

func newSPAFileServer(staticDir string) *spaFileServer {
	return &spaFileServer{
		staticDir:  staticDir,
		fileServer: http.FileServer(http.Dir(staticDir)),
	}
}

func (s *spaFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the URL path to prevent directory traversal
	urlPath := filepath.Clean(r.URL.Path)

	// Build the full file path on disk
	filePath := filepath.Join(s.staticDir, urlPath)

	// Check if the requested file exists on disk
	info, err := os.Stat(filePath)
	if err == nil && !info.IsDir() {
		// File exists — serve it directly (CSS, JS, images, fonts, etc.)
		s.fileServer.ServeHTTP(w, r)
		return
	}

	// If it's a directory, check for an index.html inside it
	if err == nil && info.IsDir() {
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			s.fileServer.ServeHTTP(w, r)
			return
		}
	}

	// File doesn't exist — serve the root index.html (SPA fallback)
	http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// ── Configure structured logging ────────────────────────────────────────────────
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// ── Initialize services ────────────────────────────────────────────────
	services.InitS3Client()

	pgPool, err := services.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgPool.Close()

	webhookQueue := services.NewWebhookQueue(pgPool)
	bookingStore := services.NewBookingRequestStore(pgPool)

	sanityClient := services.NewSanityClient(
		os.Getenv("SANITY_PROJECT_ID"),
		os.Getenv("SANITY_DATASET"),
		"2025-02-19",
		os.Getenv("SANITY_API_TOKEN"),
	)

	squareClient := services.NewSquareClient(
		os.Getenv("SQUARE_API_KEY"),
		"2026-01-22",
	)

	bookingFlow := services.NewBookingFlowService(bookingStore, squareClient)

	// ── Initialize handlers ────────────────────────────────────────────────
	newsletterHandler := handlers.NewNewsletterHandler()
	webhookHandler := handlers.NewWebhookHandler(
		os.Getenv("WEBHOOK_SECRET"),
		os.Getenv("SQUARE_WEBHOOK_SIGNATURE_KEY"),
		os.Getenv("SQUARE_WEBHOOK_URL"),
		sanityClient,
		squareClient,
		webhookQueue,
		bookingFlow,
	)
	postHandler := handlers.NewPostHandler(sanityClient)
	shopHandler := handlers.NewShopHandler(squareClient)
	bookingHandler := handlers.NewBookingHandler(squareClient, bookingFlow)
	diagnosticHandler := handlers.NewDiagnosticHandler(newsletterHandler)

	squareWorker := services.NewWebhookWorker(
		webhookQueue,
		services.WebhookProviderSquare,
		webhookHandler,
	)
	go squareWorker.Start(ctx)

	// ── API routes ────────────────────────────────────────────────
	api := http.NewServeMux()

	// Healthcheck
	api.HandleFunc("GET /healthcheck", webhookHandler.Health)

	// Newsletter
	api.HandleFunc("POST /subscribe", newsletterHandler.Subscribe)
	api.HandleFunc("POST /send-mail", newsletterHandler.SendMail)

	// Diagnostic
	api.HandleFunc("POST /diagnostic", diagnosticHandler.Submit)

	// Blog posts
	api.HandleFunc("GET /posts", postHandler.ListPosts)
	api.HandleFunc("GET /posts/recent", postHandler.RecentPosts)
	api.HandleFunc("GET /posts/{slug}", postHandler.GetPost)

	// Webhooks & cache
	api.HandleFunc("POST /revalidate", webhookHandler.Revalidate)
	api.HandleFunc("GET /cache/status", webhookHandler.LastInvalidation)
	api.HandleFunc("POST /webhooks/square", webhookHandler.HandleSquareWebhook)

	// Shop
	api.HandleFunc("GET /shop/catalog", shopHandler.GetCatalog)
	api.HandleFunc("GET /shop/images/{imageID}", shopHandler.GetProductImage)
	api.HandleFunc("POST /checkout", shopHandler.CreateCheckoutLink)

	// Booking
	api.HandleFunc("GET /booking/config", bookingHandler.GetPaymentConfig)
	api.HandleFunc("GET /booking/services", bookingHandler.GetServices)
	api.HandleFunc("POST /booking/availability", bookingHandler.GetAvailability)
	api.HandleFunc("POST /booking/request", bookingHandler.CreateRequest)
	api.HandleFunc("POST /booking/requests/{id}/store-card", bookingHandler.StoreCardAndBook)
	api.HandleFunc("POST /booking/requests/{id}/book-member", bookingHandler.BookFreeService)

	// ── Static frontend files ───────────────────────────────────────────
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./ui/build"
	}
	slog.Info("Serving static files", "dir", staticDir)

	// ── Root mux ───────────────────────────────────────────
	rootMux := http.NewServeMux()

	rootMux.HandleFunc("GET /feed.xml", shopHandler.GetProductFeed)

	rootMux.Handle("/api/", http.StripPrefix("/api", api))
	rootMux.Handle("/", middleware.StaticCacheMiddleware(newSPAFileServer(staticDir)))

	// Wrap the root mux with middleware
	handler := middleware.Logger(rootMux)
	handler = middleware.CorsMiddleware(handler)

	// ── Start server ─────────────────────────────────────────────────
	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	go func() {
		slog.Info("Server starting on", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shut down server", "error", err)
	}
}
