package store

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type elem[T any] struct {
    key string
    val *T
    next *elem[T]
}

type bucket[T any] struct {
    first *elem[T]
    last *elem[T]
    sync.RWMutex
}

type hashmap[T any] struct {
    buckets []bucket[T]
}

func (b *bucket[T]) Write(key string, val T) {
    b.Lock()
    defer b.Unlock()

    if b.first == nil {
        b.first = &elem[T]{
            key: key,
            val: &val,
        }
        b.last = b.first

        return
    }

    b.last.next = &elem[T]{
        key: key,
        val: &val,
    }
    b.last = b.last.next
}

func (b *bucket[T]) Read(key string) (*T, bool) {
    b.RLock()
    defer b.RUnlock()

    if b.first == nil {
        return nil, false
    }
        
    cur := b.first
    for ; cur != nil && cur.key != key; cur = cur.next {}

    if cur == nil {
        return nil, false
    }

    return cur.val, true
}

func (b *bucket[T]) Delete(key string) {
    b.Lock()
    defer b.Unlock()

    if b.first == nil {
        return
    }

    cur := b.first
    for ; cur.next != nil && cur.next.key != key; cur = cur.next {}

    if cur.next == nil {
        return
    }

    cur.next = cur.next.next
}

func NewHashmap[T any](bsize int) *hashmap[T] {
    return &hashmap[T] {
        buckets: make([]bucket[T], bsize),
    }
}

func (hm *hashmap[T]) Set(key string, value T) {
    hm.buckets[murmur3.Sum64([]byte(key)) % uint64(len(hm.buckets))].Write(key, value)
}

func (hm *hashmap[T]) Get(key string) (*T, bool) {
    return hm.buckets[murmur3.Sum64([]byte(key)) % uint64(len(hm.buckets))].Read(key)
}

func (hm *hashmap[T]) Delete(key string) {
    hm.buckets[murmur3.Sum64([]byte(key)) % uint64(len(hm.buckets))].Delete(key)
}
