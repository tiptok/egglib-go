package model

import (
	"encoding/json"
	"time"
)

type Item struct {
	Object     interface{} `json:"object"`     // object
	TTL        int         `json:"ttl"`        // key ttl, in second
	Expiration int64       `json:"expiration"` // expired keys will be deleted from redis.
	MarshData  []byte      `json:"-"`
}

func (item Item) Expire() bool {
	if item.Expiration == 0 {
		return false
	}

	if item.Expiration < time.Now().UnixNano() {
		return true
	}
	return false
}

func (item Item) Data() []byte {
	return item.MarshData
}

func (item Item) String() string {
	d, _ := json.Marshal(item)
	return string(d)
}

func NewItem(v interface{}, d int) *Item {
	ttl := d
	var e int64
	if d > 0 {
		e = time.Now().Add(time.Duration(d*1) * time.Second).UnixNano() //lazyFactor
	}

	return &Item{
		Object:     v,
		TTL:        ttl,
		Expiration: e,
	}
}
