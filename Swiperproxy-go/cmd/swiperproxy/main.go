package main

import (
	"log"
	"net/http"
	"os"
	"swiperproxy-go/internal/config"
	"swiperproxy-go/internal/middleware"
	"swiperproxy-go/internal/proxy"
)

func main() {
	cfgPath := "configs/config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cache := proxy.NewCache(cfg.Cache.TTL, cfg.Cache.MaxEntries)
	handler := proxy.NewHandler(cache, cfg.Proxy.Timeout)

	var h http.Handler = handler
	h = middleware.Logger(h)
	h = middleware.RateLimit(h, cfg.RateLimit.Requests, cfg.RateLimit.Window)
	h = middleware.SecurityHeaders(h)

	addr := cfg.Server.Listen
	log.Printf("swiperproxy-go listening on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatalf("server: %v", err)
	}
}