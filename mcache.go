package mcache

import (
	"appengine"
	"appengine/memcache"
	"reflect"
)

var (
	ErrCacheMiss = memcache.ErrCacheMiss
)

func Add(c appengine.Context, item *memcache.Item) error {
	return memcache.Add(c, item)
}

func AddMulti(c appengine.Context, items []*memcache.Item) error {
	return memcache.AddMulti(c, items)
}

func CompareAndSwap(c appengine.Context, item *memcache.Item) error {
	return memcache.CompareAndSwap(c, item)
}

func CompareAndSwapMulti(c appengine.Context, items []*memcache.Item) error {
	return memcache.CompareAndSwapMulti(c, items)
}

func Delete(c appengine.Context, key string) error {
	return memcache.Delete(c, key)
}

func DeleteMulti(c appengine.Context, keys []string) error {
	return memcache.DeleteMulti(c, keys)
}

func Flush(c appengine.Context) error {
	return memcache.Flush(c)
}

func Get(c appengine.Context, key string) (*memcache.Item, error) {
	return memcache.Get(c, key)
}

func GetMulti(c appengine.Context, keys []string) ([]*memcache.Item, error) {
	itemsMap, err := memcache.GetMulti(c, keys)
	if err != nil {
		return nil, err
	}

	items := make([]*memcache.Item, len(keys))
	multiErr, errsNil := make(appengine.MultiError, len(keys)), true
	for i, key := range keys {
		if item, ok := itemsMap[key]; ok {
			items[i] = item
		} else {
			multiErr[i] = memcache.ErrCacheMiss
			errsNil = false
		}
	}
	if errsNil {
		return items, nil
	}
	return items, multiErr
}

func Increment(c appengine.Context,
	key string, delta int64, initialValue uint64) (newValue uint64, err error) {
	return memcache.Increment(c, key, delta, initialValue)
}

func IncrementExisting(c appengine.Context,
	key string, delta int64) (newValue uint64, err error) {
	return IncrementExisting(c, key, delta)
}

func Set(c appengine.Context, item *memcache.Item) error {
	return memcache.Set(c, item)
}

func SetMulti(c appengine.Context, items []*memcache.Item) error {
	return memcache.SetMulti(c, items)
}

type Codec struct {
	memcache.Codec
}

func (cd Codec) GetMulti(c appengine.Context, keys []string,
	dst interface{}) ([]*memcache.Item, error) {

	itemsMap, err := memcache.GetMulti(c, keys)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(dst)
	items := make([]*memcache.Item, len(keys))
	multiErr, errsNil := make(appengine.MultiError, len(keys)), true
	for i, key := range keys {
		if item, ok := itemsMap[key]; ok {
			err := cd.Unmarshal(item.Value, v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			item.Object = v.Index(i).Interface()
			items[i] = item
		} else {
			multiErr[i] = memcache.ErrCacheMiss
			errsNil = false
		}
	}
	if errsNil {
		return items, nil
	}
	return items, multiErr
}

func (cd Codec) Get(c appengine.Context,
	key string, dst interface{}) (*memcache.Item, error) {
	return cd.Codec.Get(c, key, dst)
}

func (cd Codec) Set(c appengine.Context, item *memcache.Item) error {
	return cd.Codec.Set(c, item)
}

func (cd Codec) SetMulti(c appengine.Context, items []*memcache.Item) error {
	return cd.Codec.SetMulti(c, items)
}

func (cd Codec) Add(c appengine.Context, item *memcache.Item) error {
	return cd.Codec.Add(c, item)
}

func (cd Codec) AddMulti(c appengine.Context, items []*memcache.Item) error {
	return cd.Codec.AddMulti(c, items)
}

func (cd Codec) CompareAndSwap(c appengine.Context, item *memcache.Item) error {
	return cd.Codec.CompareAndSwap(c, item)
}

func (cd Codec) CompareAndSwapMulti(c appengine.Context,
	items []*memcache.Item) error {
	return cd.Codec.CompareAndSwapMulti(c, items)
}

var (
	Gob  = Codec{Codec: memcache.Gob}
	JSON = Codec{Codec: memcache.JSON}
)
