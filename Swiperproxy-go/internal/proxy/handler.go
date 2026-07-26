package proxy

import (
	"io"
	"net/http"
	"time"
)

type Handler struct {
	cache   *Cache
	timeout time.Duration
	client  *http.Client
}

func NewHandler(cache *Cache, timeout time.Duration) *Handler {
	return &Handler{
		cache:   cache,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
		return
	}

	if cached, ok := h.cache.Get(targetURL); ok {
		w.Write(cached.Body)
		for k, v := range cached.Headers {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, "invalid target URL", http.StatusBadRequest)
		return
	}

	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, "upstream request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	entry := CacheEntry{
		Body:    body,
		Headers: resp.Header,
	}
	h.cache.Set(targetURL, entry)

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}