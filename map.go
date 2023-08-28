// Copyright 2023 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb

import "cmp"

type entry[K cmp.Ordered, V any] struct {
	key   K
	value V
}

type LLRBMap[K cmp.Ordered, V any] struct {
	tr *LLRBTree[*entry[K, V]]
}

func compareMapEntry[K cmp.Ordered, V any](e1, e2 *entry[K, V]) int {
	return cmp.Compare(e1.key, e2.key)
}

// NewMap creates a new LLRBMap.
func NewMap[K cmp.Ordered, V any]() *LLRBMap[K, V] {
	return &LLRBMap[K, V]{
		tr: New[*entry[K, V]](compareMapEntry[K, V]),
	}
}

// Set inserts or replaces a key-value pair in the map.
// It returns the previous value associated with the key
// and a boolean indicating if the key existed.
func (m *LLRBMap[K, V]) Set(key K, value V) (V, bool) {
	prev, exist := m.tr.ReplaceOrInsert(&entry[K, V]{
		key:   key,
		value: value,
	})
	if exist {
		return prev.value, true
	}
	return zero[V](), false
}

// Get retrieves the value associated with the specified key from the map.
// It returns the value and a boolean indicating if the key exists in the map.
func (m *LLRBMap[K, V]) Get(key K) (V, bool) {
	ent, ok := m.tr.Get(&entry[K, V]{key: key})
	if ok {
		return ent.value, true
	}
	return zero[V](), false
}

// Delete removes the key-value pair with the specified key from the map.
// It returns the value associated with the key and a boolean indicating
// if the key existed.
func (m *LLRBMap[K, V]) Delete(key K) (V, bool) {
	ent, ok := m.tr.Delete(&entry[K, V]{key: key})
	if ok {
		return ent.value, true
	}
	return zero[V](), false
}

// Range iterates over the key-value pairs in the map in ascending order of the keys.
// The provided callback function is called for each key-value pair.
// Iteration stops if the callback function returns false.
func (m *LLRBMap[K, V]) Range(iter func(key K, value V) bool) {
	m.tr.Ascend(func(ent *entry[K, V]) bool {
		return iter(ent.key, ent.value)
	})
}

// Has checks if the map contains the specified key.
// It returns true if the key exists in the map, false otherwise.
func (m *LLRBMap[K, V]) Has(key K) bool {
	return m.tr.Has(&entry[K, V]{key: key})
}

// Len returns the number of key-value pairs in the map.
func (m *LLRBMap[K, V]) Len() int {
	return m.tr.Len()
}

// Clear removes all key-value pairs from the map, resulting in an empty map.
func (m *LLRBMap[K, V]) Clear() {
	m.tr.Clear()
}
