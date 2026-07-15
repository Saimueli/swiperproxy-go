package proxy

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	Value      []byte
	Header     http.Header
	StatusCode int
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

type Cache struct {
	items    map[string]*list.Element
	order    *list.List
	mu       sync.RWMutex
	maxSize  int
	ttl      int
	enabled  bool
}

type CacheEntry struct {
	key        string
	value      []byte
	header     http.Header
	statusCode int
	createdAt  time.Time
	expiresAt  time.Time
}

func NewCache(maxSize int, ttl int, enabled bool) *Cache {
	return &Cache{
		items:    make(map[string]*list.Element),
		order:    list.New(),
		maxSize:  maxSize,
		ttl:      ttl,
		enabled:  enabled,
	}
}

func (c *Cache) Get(key string) ([]byte, http.Header, int, bool) {
	if !c.enabled {
		return nil, nil, 0, false
	}

	c.mu.RLock()
	elem, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		return nil, nil, 0, false
	}

	entry := elem.Value.(*CacheEntry)

	if time.Now().After(entry.expiresAt) {
		c.Delete(key)
		return nil, nil, 0, false
	}

	c.mu.Lock()
	c.order.MoveToFront(elem)
	c.mu.Unlock()

	return entry.value, entry.header, entry.statusCode, true
}

func (c *Cache) Set(key string, value []byte, header http.Header, statusCode int) {
	if !c.enabled || c.maxSize <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.order.Remove(elem)
		delete(c.items, key)
	}

	if c.order.Len() >= c.maxSize {
		c.evict()
	}

	entry := &CacheEntry{
		key:        key,
		value:      value,
		header:     header,
		statusCode: statusCode,
		createdAt:  time.Now(),
		expiresAt:  time.Now().Add(time.Duration(c.ttl) * time.Second),
	}

	elem := c.order.PushFront(entry)
	c.items[key] = elem
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

func (c *Cache) evict() {
	elem := c.order.Back()
	if elem != nil {
		c.order.Remove(elem)
		entry := elem.Value.(*CacheEntry)
		delete(c.items, entry.key)
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.order.Init()
}

func GenerateCacheKey(r *http.Request) string {
	hash := sha256.Sum256([]byte(r.Method + r.URL.String()))
	return hex.EncodeToString(hash[:])
}

func (c *Cache) IsCacheable(r *http.Request) bool {
	if !c.enabled {
		return false
	}

	if r.Method != http.MethodGet {
		return false
	}

	if r.Header.Get("Cache-Control") == "no-cache" {
		return false
	}

	return true
}

func (c *Cache) ShouldCacheResponse(resp *http.Response) bool {
	if resp.StatusCode != http.StatusOK {
		return false
	}

	if resp.Header.Get("Cache-Control") == "no-store" {
		return false
	}

	contentLength := resp.ContentLength
	if contentLength > 10*1024*1024 {
		return false
	}

	return true
}

func CopyResponse(original *http.Response) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(original.Body)
	if err != nil {
		return nil, err
	}
	original.Body.Close()

	newResp := &http.Response{
		Status:        original.Status,
		StatusCode:    original.StatusCode,
		Proto:         original.Proto,
		ProtoMajor:    original.ProtoMajor,
		ProtoMinor:    original.ProtoMinor,
		Header:        original.Header.Clone(),
		Body:          io.NopCloser(bytes.NewReader(bodyBytes)),
		ContentLength: int64(len(bodyBytes)),
	}

	original.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return newResp, nil
}