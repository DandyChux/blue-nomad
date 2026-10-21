package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/dandychux/blue-nomad/internal/services"
)

const (
	postsCacheTTL = 5 * time.Minute
	postCacheTTL  = 10 * time.Minute
)

// GROQ queries — direct translations from the old Next.js queries.ts
const (
	groqAllPosts = `*[_type == "post"] | order(_createdAt desc) {
		_id,
		_createdAt,
		title,
		description,
		"slug": slug.current,
		"imageUrl": mainImage.asset->url,
		"mainImage": mainImage {
			asset-> { _ref, url, metadata },
			alt
		},
		categories[]->{ title, "slug": slug.current },
		"author": author->{ name, bio }
	}`

	groqPostBySlug = `*[_type == "post" && slug.current == $slug][0] {
		_id,
		title,
		description,
		"slug": slug.current,
		"mainImage": mainImage {
			asset-> { url, metadata },
			alt
		},
		body[] {
			...,
			_type == "image" => {
				...,
				asset-> { url, metadata }
			}
		},
		categories[]->{ title },
		"date": _createdAt,
		"authorName": author->name,
		"authorImage": author->image
	}`

	groqRecentPosts = `*[_type == "post"] | order(_createdAt desc) [0...4] {
		"slug": slug.current,
		title,
		description,
		"date": _createdAt,
		"image": mainImage.asset->url
	}`
)

// PostHandler serves blog post data from Sanity.
type PostHandler struct {
	sanity *services.SanityClient
}

func NewPostHandler(sanity *services.SanityClient) *PostHandler {
	return &PostHandler{sanity: sanity}
}

// ListPosts returns all published posts.
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	data, err := h.sanity.QueryCached(r.Context(), "posts:all", groqAllPosts, nil, postsCacheTTL)
	if err != nil {
		slog.Error("failed to fetch posts", "error", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetPost returns a single post by slug.
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	cacheKey := "post:" + slug
	params := map[string]string{"slug": slug}

	data, err := h.sanity.QueryCached(r.Context(), cacheKey, groqPostBySlug, params, postCacheTTL)
	if err != nil {
		slog.Error("failed to fetch post", "slug", slug, "error", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	if string(data) == "null" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// RecentPosts returns the 4 most recent posts.
func (h *PostHandler) RecentPosts(w http.ResponseWriter, r *http.Request) {
	data, err := h.sanity.QueryCached(r.Context(), "posts:recent", groqRecentPosts, nil, postsCacheTTL)
	if err != nil {
		slog.Error("failed to fetch recent posts", "error", err)
		http.Error(w, "Failed to fetch recent posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
