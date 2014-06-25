package mcache

import (
	"appengine"
	"appengine/memcache"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

func Add(c appengine.Context, item *memcache.Item) error {
	return memcache.Add(c, item)
}

func AddMulti(c appengine.Context, items []*memcache.Item) error {
	err := memcache.AddMulti(c, items)
	if err == nil {
		return make(appengine.MultiError, len(items))
	}
	return err
}

func CompareAndSwap(c appengine.Context, item *memcache.Item) error {
	return memcache.CompareAndSwap(c, item)
}

func CompareAndSwapMulti(c appengine.Context, items []*memcache.Item) error {
	err := memcache.CompareAndSwapMulti(c, items)
	if err == nil {
		return make(appengine.MultiError, len(items))
	}
	return err
}

func Delete(c appengine.Context, key string) error {
	return memcache.Delete(c, key)
}

func DeleteMulti(c appengine.Context, keys []string) error {
	err := memcache.DeleteMulti(c, keys)
	if err == nil {
		return make(appengine.MultiError, len(keys))
	}
	return err
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
	errs := make(appengine.MultiError, len(keys))
	for i, key := range keys {
		if item, ok := itemsMap[key]; ok {
			items[i] = item
		} else {
			errs[i] = memcache.ErrCacheMiss
		}
	}

	return items, errs
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
	err := memcache.SetMulti(c, items)
	if err == nil {
		return make(appengine.MultiError, len(items))
	}
	return err
}

type Codec struct {
	Marshal   func(interface{}) ([]byte, error)
	Unmarshal func([]byte, interface{}) error
}

func (cd Codec) GetMulti(c appengine.Context, keys []string,
	dst interface{}) ([]*memcache.Item, error) {

	itemsMap, err := memcache.GetMulti(c, keys)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(dst)
	items := make([]*memcache.Item, len(keys))
	multiErr := make(appengine.MultiError, len(keys))
	for i, key := range keys {
		if item, ok := itemsMap[key]; ok {
			err := cd.Unmarshal(item.Value, v.Index(i).Interface())
			if err != nil {
				multiErr[i] = err
			} else {
				item.Object = v.Index(i).Interface()
				items[i] = item
			}
		} else {
			multiErr[i] = memcache.ErrCacheMiss
		}
	}
	return items, multiErr
}

func (cd Codec) Get(c appengine.Context,
	key string, dst interface{}) (*memcache.Item, error) {

	items, err := cd.GetMulti(c, []string{key}, []interface{}{dst})
	if err == nil {
		return items[0], err
	} else if me, ok := err.(appengine.MultiError); ok {
		return nil, me[0]
	}
	return nil, err
}

func (cd Codec) marshalItems(items []*memcache.Item) error {
	for _, item := range items {

		v, err := cd.Marshal(item.Object)
		if err != nil {
			return err
		}
		item.Value = v
	}
	return nil
}

func (cd Codec) Set(c appengine.Context, item *memcache.Item) error {
	err := cd.SetMulti(c, []*memcache.Item{item})
	if me, ok := err.(appengine.MultiError); ok {
		return me[0]
	}
	return err
}

func (cd Codec) SetMulti(c appengine.Context, items []*memcache.Item) error {
	if err := cd.marshalItems(items); err != nil {
		return err
	}
	return SetMulti(c, items)
}

func (cd Codec) Add(c appengine.Context, item *memcache.Item) error {
	err := cd.AddMulti(c, []*memcache.Item{item})
	if me, ok := err.(appengine.MultiError); ok {
		return me[0]
	}
	return err
}

func (cd Codec) AddMulti(c appengine.Context, items []*memcache.Item) error {
	if err := cd.marshalItems(items); err != nil {
		return err
	}
	return AddMulti(c, items)
}

func (cd Codec) CompareAndSwap(c appengine.Context, item *memcache.Item) error {
	return cd.CompareAndSwapMulti(c, []*memcache.Item{item})
}

func (cd Codec) CompareAndSwapMulti(c appengine.Context,
	items []*memcache.Item) error {
	if err := cd.marshalItems(items); err != nil {
		return err
	}
	return CompareAndSwapMulti(c, items)
}

var (
	Gob  = Codec{gobMarshal, gobUnmarshal}
	JSON = Codec{json.Marshal, json.Unmarshal}
)

func gobMarshal(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobUnmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
