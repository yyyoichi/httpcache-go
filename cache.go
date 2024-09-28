package httpcache

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Cache interface {
	Put(Object) error
	Query(Object) (io.Reader, error)
}

var ErrNoCache = errors.New("no cache")

var DefaultStorageCache = &StorageCache{}

type StorageCache struct {
	dir string
}

func NewStorageCache(dir string) *StorageCache {
	c := &StorageCache{dir: dir}
	c.init()
	if _, err := os.Stat(dir); err != nil {
		_ = os.Mkdir(dir, 0755)
	}
	return c
}

func (c *StorageCache) Put(o Object) error {
	c.init()
	filename := c.dir + o.Key()
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := io.Copy(f, o.NewReader()); err != nil {
		return err
	}
	return nil
}

func (c *StorageCache) Query(o Object) (io.Reader, error) {
	c.init()
	filename := c.dir + o.Key()
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNoCache, err)
	}
	defer f.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}
	return &buf, nil
}

func (c *StorageCache) init() {
	if c.dir == "" {
		c.dir = "./"
	}
	if !strings.HasSuffix(c.dir, "/") {
		c.dir += "/"
	}
}

var DefaultMemoryCache = &MemoryCache{}

type MemoryCache struct {
	store map[string][]byte
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{store: make(map[string][]byte)}
}

func (c *MemoryCache) Put(o Object) error {
	var buf bytes.Buffer
	buf.Grow(int(o.Length()))
	if _, err := io.Copy(&buf, o.NewReader()); err != nil {
		return err
	}

	key := o.Key()
	c.store[key] = buf.Bytes()
	return nil
}

func (c *MemoryCache) Query(o Object) (io.Reader, error) {
	key := o.Key()
	b, ok := c.store[key]
	if !ok {
		return nil, ErrNoCache
	}
	return bytes.NewReader(b), nil
}

func (c *MemoryCache) init() {
	if c.store == nil {
		c.store = make(map[string][]byte)
	}
}
