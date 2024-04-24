# HTTP-Cache-Go

HttpClient wrapper to cache http response Body in storage or memory.

## Example

```golang

client := httpcache.DefaultClient
resp, err := client.Get("https://example.com/content.json")

```

If it does not exist in storage, make an Http request, cache the result, and use it for the next time.

## Custom

```golang
client := &httpcache.Client{
    Client: &http.Client{
        // 
    },
    Cache: httpcache.NewMemoryCache(),
    Handler: httpcache.NewOnceLatestHandler(),
}

```

### Cache

You can choose local storage or memory for storage location.

```golang
NewStorageCache("/tmp")
httpcache.NewMemoryCache()
```

Or implements Cache interface.

```golang

type Cache interface {
 Put(Object) error
 Query(Object) (io.Reader, error)
}
```

### Handler

You may decide to use the cache or not.

```golang
// Always use the latest.
NewLatestHandler = func() *Handler
// Get the latest only once at the beginning.
NewOnceLatestHandler = func() *Handler
// If there is a cache, use it; if not, get the latest and save the cache.
NewDefaultHandler = func() *Handler
// Always use and save cache.
NewSimpleHandler = func() *Handler
```

Or your custom Handler

```golang
type Handler struct {
 mu   sync.Mutex
 keys []string
 Pre  func(Cache, Object) (io.Reader, error)
 Post func(Cache, Object) error
}
```
