// Copyright 2023 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb

import (
	"cmp"
)

type LLRBSet[T cmp.Ordered] struct {
	tr *LLRBTree[T]
}

// NewSet creates a new LLRBSet.
func NewSet[T cmp.Ordered]() *LLRBSet[T] {
	return &LLRBSet[T]{
		tr: NewLLRBTree[T](cmp.Compare[T]),
	}
}

// Insert inserts a value into the set.
// It returns true if the value already exists in the set, false otherwise.
func (s *LLRBSet[T]) Insert(item T) (exist bool) {
	_, exist = s.tr.ReplaceOrInsert(item)
	return exist
}

// Delete removes a value from the set.
// It returns true if the value existed in the set, false otherwise.
func (s *LLRBSet[T]) Delete(item T) (exist bool) {
	_, exist = s.tr.Delete(item)
	return exist
}

// Range iterates over the values in the set in ascending order.
// The provided callback function is called for each value.
// Iteration stops if the callback function returns false.
func (s *LLRBSet[T]) Range(iter IterFunc[T]) {
	s.tr.Ascend(iter)
}

// Has checks if the set contains the specified value.
// It returns true if the value exists in the set, false otherwise.
func (s *LLRBSet[T]) Has(item T) bool {
	return s.tr.Has(item)
}

// Len returns the number of values in the set.
func (s *LLRBSet[T]) Len() int {
	return s.tr.Len()
}

// Clear removes all values from the set, resulting in an empty set.
func (s *LLRBSet[T]) Clear() {
	s.tr.Clear()
}
