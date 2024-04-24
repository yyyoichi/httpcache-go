package httpcache

import (
	"io"
	"slices"
	"sync"
)

type Handler struct {
	mu   sync.Mutex
	keys []string
	Pre  func(Cache, Object) (io.Reader, error)
	Post func(Cache, Object) error
}

var (
	// Always use the latest.
	NewLatestHandler = func() *Handler {
		return &Handler{
			Pre:  func(c Cache, o Object) (io.Reader, error) { return nil, nil },
			Post: func(c Cache, o Object) error { return c.Put(o) },
		}
	}
	// Get the latest only once at the beginning.
	NewOnceLatestHandler = func() *Handler {
		var hander Handler
		hander.keys = make([]string, 0, 1000) // latest
		hander.Pre = func(c Cache, o Object) (io.Reader, error) {
			hander.mu.Lock()
			defer hander.mu.Unlock()
			// already got
			if slices.Contains(hander.keys, o.Key()) {
				return c.Query(o)
			}
			return nil, nil
		}
		hander.Post = func(c Cache, o Object) error {
			hander.mu.Lock()
			defer hander.mu.Unlock()
			// already saved
			if slices.Contains(hander.keys, o.Key()) {
				return nil
			}
			hander.keys = append(hander.keys, o.Key())
			return c.Put(o)
		}
		return &hander
	}
	// If there is a cache, use it; if not, get the latest and save the cache.
	NewDefaultHandler = func() *Handler {
		var hander Handler
		hander.keys = make([]string, 0, 1000) // use cache
		hander.Pre = func(c Cache, o Object) (io.Reader, error) {
			hander.mu.Lock()
			defer hander.mu.Unlock()
			r, err := c.Query(o)
			if err != nil {
				return r, err
			}
			// use cache
			hander.keys = append(hander.keys, o.Key())
			return r, nil
		}
		hander.Post = func(c Cache, o Object) error {
			hander.mu.Lock()
			defer hander.mu.Unlock()
			// already cached data
			if slices.Contains(hander.keys, o.Key()) {
				return nil
			}
			return c.Put(o)
		}
		return &hander
	}
	// Always use and save cache.
	NewSimpleHandler = func() *Handler {
		var hander Handler
		hander.Pre = func(c Cache, o Object) (io.Reader, error) {
			return c.Query(o)
		}
		hander.Post = func(c Cache, o Object) error {
			return c.Put(o)
		}
		return &hander
	}
)
