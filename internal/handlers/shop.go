package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dandychux/blue-nomad/internal/services"
)

// ShopHandler serves catalog data and processes payments via Square.
type ShopHandler struct {
	square *services.SquareClient
}

type googleProductFeed struct {
	XMLName xml.Name          `xml:"rss"`
	Version string            `xml:"version,attr"`
	Channel googleFeedChannel `xml:"channel"`
}

type googleFeedChannel struct {
	Title       string              `xml:"title"`
	Link        string              `xml:"link"`
	Description string              `xml:"description"`
	Items       []googleFeedProduct `xml:"item"`
}

type googleFeedProduct struct {
	ID           string `xml:"id"`
	Title        string `xml:"title"`
	Description string `xml:"description"`
	Link         string `xml:"link"`
	ImageLink    string `xml:"image_link,omitempty"`
	Availability string `xml:"availability"`
	Price        string `xml:"price"`
	Condition    string `xml:"condition"`
}

func (p googleFeedProduct) MarshalXML(
	encoder *xml.Encoder,
	start xml.StartElement,
) error {
	start.Name.Local = "item"

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}

	const googleNamespace = "http://base.google.com/ns/1.0"

	fields := []struct {
		name  string
		value string
		omit  bool
	}{
		{"id", p.ID, false},
		{"title", p.Title, false},
		{"description", p.Description, false},
		{"link", p.Link, false},
		{"image_link", p.ImageLink, p.ImageLink == ""},
		{"availability", p.Availability, false},
		{"price", p.Price, false},
		{"condition", p.Condition, false},
	}

	for _, field := range fields {
		if field.omit {
			continue
		}

		element := xml.StartElement{
			Name: xml.Name{
				Space: googleNamespace,
				Local: field.name,
			},
		}

		if err := encoder.EncodeElement(field.value, element); err != nil {
			return err
		}
	}

	return encoder.EncodeToken(start.End())
}

func NewShopHandler(square *services.SquareClient) *ShopHandler {
	return &ShopHandler{square: square}
}

// GetCatalog returns all items from the Square catalog.
func (h *ShopHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	data, err := h.square.GetCatalogItems(r.Context())
	if err != nil {
		slog.Error("failed to fetch square catalog", "error", err)
		http.Error(w, "Failed to fetch catalog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *ShopHandler) GetProductImage(
	w http.ResponseWriter,
	r *http.Request,
) {
	imageID := r.PathValue("imageID")
	if imageID == "" {
		http.Error(w, "Missing image ID", http.StatusBadRequest)
		return
	}

	image, err := h.square.GetCatalogImage(r.Context(), imageID)
	if err != nil {
		slog.Warn(
			"failed to fetch Square catalog image",
			"image_id", imageID,
			"error", err,
		)
		http.NotFound(w, r)
		return
	}

	http.Redirect(
		w,
		r,
		image.ImageData.URL,
		http.StatusTemporaryRedirect,
	)
}

// CreateCheckoutLink handles incoming cart payloads and returns a Square hosted checkout URL.
func (h *ShopHandler) CreateCheckoutLink(w http.ResponseWriter, r *http.Request) {
	var req services.CheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid checkout request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Call the Square service to generate the link
	checkoutURL, err := h.square.CreatePaymentLink(r.Context(), req)
	if err != nil {
		var stockErr *services.InsufficientStockError
		if errors.As(err, &stockErr) {
			slog.Info("checkout blocked by inventory",
				"variation_id", stockErr.VariationID,
				"requested", stockErr.Requested,
				"available", stockErr.Available,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]any{
				"error":        "insufficient_stock",
				"message":      "One or more items in your bag are no longer available in the quantity requested.",
				"variation_id": stockErr.VariationID,
				"requested":    stockErr.Requested,
				"available":    stockErr.Available,
			})
			return
		}

		slog.Error("failed to create square payment link", "error", err)
		http.Error(w, "Checkout initialization failed", http.StatusBadGateway)
		return
	}

	// Return the URL as JSON so SvelteKit can redirect the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": checkoutURL,
	})
}

// GetProductFeed dynamically generates a Google Merchant Center-compatible
// RSS 2.0 product feed.
func (h *ShopHandler) GetProductFeed(w http.ResponseWriter, r *http.Request) {
	data, err := h.square.GetCatalogItems(r.Context())
	if err != nil {
		slog.Error(
			"failed to fetch square catalog for product feed",
			"error", err,
		)
		http.Error(
			w,
			"Failed to generate product feed",
			http.StatusBadGateway,
		)
		return
	}

	var catalog struct {
		Objects []struct {
			Type     string `json:"type"`
			ID       string `json:"id"`
			ItemData struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				ProductType string `json:"product_type"`
				ImageIDs    []string `json:"image_ids"`
				Variations  []struct {
					ItemVariationData struct {
						PriceMoney struct {
							Amount   int64  `json:"amount"`
							Currency string `json:"currency"`
						} `json:"price_money"`
					} `json:"item_variation_data"`
				} `json:"variations"`
			} `json:"item_data"`
		} `json:"objects"`
	}

	if err := json.Unmarshal(data, &catalog); err != nil {
		slog.Error(
			"failed to decode square catalog for product feed",
			"error", err,
		)
		http.Error(
			w,
			"Failed to generate product feed",
			http.StatusInternalServerError,
		)
		return
	}

	baseURL := os.Getenv("PUBLIC_SITE_URL")
	if baseURL == "" {
		baseURL = requestBaseURL(r)
	}
	baseURL = strings.TrimRight(baseURL, "/")

	feed := googleProductFeed{
		Version: "2.0",
		Channel: googleFeedChannel{
			Title:       "Blue Nomad Shop",
			Link:        baseURL + "/shop",
			Description: "Blue Nomad product feed",
			Items:       make([]googleFeedProduct, 0),
		},
	}

	for _, object := range catalog.Objects {
		if object.Type != "ITEM" ||
			object.ItemData.ProductType == "APPOINTMENTS_SERVICE" {
			continue
		}

		if object.ItemData.Name == "" {
			continue
		}

		product := googleFeedProduct{
			ID:           object.ID,
			Title:        object.ItemData.Name,
			Description:  object.ItemData.Description,
			Link:         baseURL + "/shop/" + url.PathEscape(object.ID),
			Availability: "in stock",
			Condition:    "new",
		}

		if len(object.ItemData.Variations) > 0 {
			priceMoney := object.ItemData.Variations[0].
				ItemVariationData.PriceMoney

			currency := strings.ToUpper(priceMoney.Currency)
			if currency == "" {
				currency = "USD"
			}

			product.Price = fmt.Sprintf(
				"%.2f %s",
				float64(priceMoney.Amount)/100,
				currency,
			)
		}

		if len(object.ItemData.ImageIDs) > 0 {
			product.ImageLink = baseURL + "/api/shop/images/" +
				url.PathEscape(object.ItemData.ImageIDs[0])
		}

		feed.Channel.Items = append(feed.Channel.Items, product)
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		slog.Error("failed to encode product feed", "error", err)
		http.Error(
			w,
			"Failed to generate product feed",
			http.StatusInternalServerError,
		)
		return
	}

	// encoding/xml does not add the g namespace declaration automatically
	// for custom namespaced elements, so prepend it explicitly.
	output = bytes.Replace(
		output,
		[]byte(`<rss version="2.0">`),
		[]byte(`<rss version="2.0" xmlns:g="http://base.google.com/ns/1.0">`),
		1,
	)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(output)
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.TrimSpace(
			strings.Split(forwardedProto, ",")[0],
		)
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = strings.TrimSpace(
			strings.Split(forwardedHost, ",")[0],
		)
	}

	return scheme + "://" + host
}
