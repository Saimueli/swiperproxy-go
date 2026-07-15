package proxy

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/saimueli/swiperproxy-go/internal/config"
)

type Handler struct {
	cfg          *config.Config
	transport    *http.Transport
	reverseProxy *httputil.ReverseProxy
	cache        *Cache
}

func NewHandler(cfg *config.Config) *Handler {
	targetURL, _ := url.Parse(cfg.Proxy.Target)

	transport := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    false,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = transport

	cache := NewCache(
		cfg.Cache.MaxSize,
		cfg.Cache.TTL,
		cfg.Cache.Enabled,
	)

	return &Handler{
		cfg:          cfg,
		transport:    transport,
		reverseProxy: proxy,
		cache:        cache,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Cache enabled: %v, IsCacheable: %v", h.cfg.Cache.Enabled, h.cache.IsCacheable(r))

	if h.cfg.Cache.Enabled && h.cache.IsCacheable(r) {
		cacheKey := GenerateCacheKey(r)

		if cachedData, header, statusCode, ok := h.cache.Get(cacheKey); ok {
			log.Printf("CACHE HIT: %s", cacheKey)
			for k, v := range header {
				w.Header()[k] = v
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(statusCode)
			w.Write(cachedData)
			return
		}
		log.Printf("CACHE MISS: %s", cacheKey)
	}

	recorder := &responseRecorder{
		ResponseWriter: w,
		HeaderMap:      make(http.Header),
		Body:           &bytes.Buffer{},
		StatusCode:     http.StatusOK,
	}

	h.reverseProxy.ServeHTTP(recorder, r)

	if h.cfg.Cache.Enabled && h.cache.IsCacheable(r) {
		if recorder.StatusCode == http.StatusOK {
			cacheKey := GenerateCacheKey(r)
			bodyBytes := recorder.Body.Bytes()

			copyHeader := recorder.HeaderMap.Clone()
			copyHeader.Set("X-Cache", "MISS")

			h.cache.Set(cacheKey, bodyBytes, copyHeader, recorder.StatusCode)

			w.Header().Set("X-Cache", "MISS")
		}
	}

	for k, v := range recorder.HeaderMap {
		w.Header()[k] = v
	}
	w.WriteHeader(recorder.StatusCode)
	w.Write(recorder.Body.Bytes())
}

type responseRecorder struct {
	http.ResponseWriter
	HeaderMap  http.Header
	Body       *bytes.Buffer
	StatusCode int
	written    bool
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	if !r.written {
		r.StatusCode = statusCode
		r.written = true
	}
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if !r.written {
		r.WriteHeader(http.StatusOK)
	}
	return r.Body.Write(data)
}

func (r *responseRecorder) Header() http.Header {
	return r.HeaderMap
}