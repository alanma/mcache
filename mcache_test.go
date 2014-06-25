package mcache_test

import (
	"appengine"
	"appengine/aetest"
	"appengine/memcache"
	"bytes"
	"github.com/qedus/mcache"
	"testing"
)

func TestMcache(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	item, err := mcache.Get(c, "one")
	if err != memcache.ErrCacheMiss {
		t.Fatalf("incorrect error: %s", err)
	}
	if item != nil {
		t.Fatal("item should be nil")
	}

	err = mcache.Add(c, &memcache.Item{
		Key:   "one",
		Value: []byte("oneval")})
	if err != nil {
		t.Fatal(err)
	}
	err = mcache.Add(c, &memcache.Item{
		Key:   "two",
		Value: []byte("twoval")})
	if err != nil {
		t.Fatal(err)
	}

	keys := []string{"one", "two"}
	items, err := mcache.GetMulti(c, keys)
	if _, ok := err.(appengine.MultiError); !ok {
		t.Fatal("appengine.MultiError expected")
	}
	for i, key := range keys {
		item := items[i]
		if key != item.Key {
			t.Fatal("incorrect key")
		}
		val := []byte(key + "val")
		if !bytes.Equal(item.Value, val) {
			t.Fatal("incorrect val")
		}
	}

	// Test cache miss error.
	keys = []string{"one", "onehalf", "two"}
	items, err = mcache.GetMulti(c, keys)
	me, ok := err.(appengine.MultiError)
	if !ok {
		t.Fatal("not appengine.MultiError")
	}
	for i, key := range keys {
		if i == 1 {
			if me[i] != memcache.ErrCacheMiss {
				t.Fatal("incorrect error")
			}
			continue
		}
		if me[i] != nil {
			t.Fatal("not nil error")
		}
		item := items[i]
		if key != item.Key {
			t.Fatal("incorrect key")
		}
		val := []byte(key + "val")
		if !bytes.Equal(item.Value, val) {
			t.Fatal("incorrect val")
		}
	}

}

func TestMcacheCodec(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	type Entity struct {
		Value int
	}

	keys := []string{"0", "1"}
	for i, key := range keys {
		err = mcache.Gob.Add(c, &memcache.Item{
			Key:    key,
			Object: &Entity{i},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	entities := make([]*Entity, len(keys))
	for i := range entities {
		entities[i] = &Entity{}
	}
	items, err := mcache.Gob.GetMulti(c, keys, entities)
	if _, ok := err.(appengine.MultiError); !ok {
		t.Fatal("appengine.MultiError expected")
	}

	for i, key := range keys {
		if entities[i].Value != i {
			t.Fatal("incorrect value")
		}
		if items[i].Key != key {
			t.Fatal("incorrect key")
		}
		if items[i].Object != entities[i] {
			t.Fatal("incorrect object")
		}
	}
}
