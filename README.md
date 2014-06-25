mcache
========

Package mcache is an API for App Engine's memcache. It is more consistent with the datastore API.

There are [good arguments](https://groups.google.com/forum/#!topic/google-appengine-go/kiuvTHf32zw/discussion) for the current official [appengine/memcache](https://developers.google.com/appengine/docs/go/memcache/reference) API. However I prefer what is created here.

The main difference is with `GetMulti`:

Package `appengine/memcache.GetMulti`:
```go
    // GetMulti is a batch version of Get. The returned map from keys to items may
    // have fewer elements than the input slice, due to memcache cache misses.
    // Each key must be at most 250 bytes in length.
    func GetMulti(c appengine.Context, key []string) (map[string]*Item, error) {
        ...
    }
```

Package `github.com/qedus/mcache.GetMulti`:
```go
    // GetMulti is a batch version of Get. The returned error can be an
    // appengine.MultiError. In which case nil error indicates a cache hit
    // and memcache.ErrCacheMiss indicates a cache miss.
    func GetMulti(c appengine.Context, keys []string) ([]*Item, error) {
        ...
    }
    
    // Use as follows:
    if items, err := GetMulti(c, keys); err == nil {
        // All your keys have returned an item.
        for _, item := range items {
            // Each item is valid here.
        }
    } else if me, ok := err.(appengine.MultiError); ok {
        // Only some keys have returned an item.
        for i, item := range items {
            if me[i] == nil {
                // Valid item.
            } else if me[i] == memcache.ErrCacheMiss {
                // Cache miss. Note that item will also be nil.
                continue
            } else {
                // This should never be reached.
            }
        }
    } else {
        // Some other error with the underlying memcache.
    }
```

Having the `GetMulti` method as above means `Codec.GetMulti` can and is easily implemented.

