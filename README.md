mcache
========

Package [mcache](http://godoc.org/github.com/qedus/mcache) is a wrapper API for App Engine's memcache. It is consistent with the datastore API which returns slices instead of maps for all the `Multi` functions.

*Important:* Each `Multi` function always returns an error. The error is likely to be a `appengine.MultiError` if the function was executed successfully. This is not conventional for Go however it helps make the code neater and easier to read in [https://github.com/qedus/nds](https://github.com/qedus/nds).

There are [good arguments](https://groups.google.com/forum/#!topic/google-appengine-go/kiuvTHf32zw/discussion) for the current official [appengine/memcache](https://developers.google.com/appengine/docs/go/memcache/reference) API. However I prefer what is created here in some cases.

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
    items, err := nds.GetMulti(c, keys)
    me, ok := err.(appengine.MultiError)
    if !ok {
        return err
    }
    // Only some keys have returned an item.
    for i, item := range items {
        if me[i] == nil {
	    // Valid item.
        } else if me[i] == memcache.ErrCacheMiss {
	    // Cache miss. Note that item will also be nil.
        } else {
	    // This should never be reached.
        }
    }
```

Having the `GetMulti` method as above means `Codec.GetMulti` can and is easily implemented.

