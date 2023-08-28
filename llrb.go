// Copyright 2023 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package llrb implements LLRB 2-3 trees, based on this paper:
//   - https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
//
// A left-leaning red–black (LLRB) tree is a type of self-balancing binary search tree.
// It is a variant of the red–black tree and guarantees the same asymptotic complexity
// for operations, but is designed to be easier to implement.
package llrb

const (
	_red   = true
	_black = false
)

type (
	CompareFunc[T any] func(a, b T) int
	IterFunc[T any]    func(a T) bool
)

// LLRBTree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees.
type LLRBTree[T any] struct {
	root    *node[T]
	compare CompareFunc[T]
	len     int
}

type node[T any] struct {
	item        T
	left, right *node[T]
	color       bool
}

// NewLLRBTree creates a new LLRB-Tree with the given compare function.
func NewLLRBTree[T any](compare CompareFunc[T]) *LLRBTree[T] {
	return &LLRBTree[T]{
		compare: compare,
	}
}

// ReplaceOrInsert adds the given item to the tree.  If an item in the tree
// already equals the given one, it is removed from the tree and returned,
// and the second return value is true.  Otherwise, (zeroValue, false)
//
// nil cannot be added to the tree (will panic).
func (t *LLRBTree[T]) ReplaceOrInsert(item T) (prev T, exist bool) {
	t.root, prev, exist = t.insert(t.root, item)
	t.root.color = _black
	if !exist {
		t.len++
	}
	return prev, exist
}

// Get looks for the key item in the tree, returning it.  It returns
// (zeroValue, false) if unable to find that item.
func (t *LLRBTree[T]) Get(item T) (T, bool) {
	x := t.root
	for x != nil {
		cmp := t.compare(item, x.item)
		if cmp == 0 {
			return x.item, true
		} else if cmp < 0 {
			x = x.left
		} else {
			x = x.right
		}
	}

	return zero[T](), false
}

// Has returns true if the given key is in the tree.
func (t *LLRBTree[T]) Has(item T) bool {
	_, ok := t.Get(item)
	return ok
}

// DeleteMin removes the smallest item in the tree and returns it.
// If no such item exists, returns nil.
func (t *LLRBTree[T]) DeleteMin() (deleted T, ok bool) {
	t.root, deleted, ok = t.deleteMin(t.root)
	if t.root != nil {
		t.root.color = _black
	}
	if ok {
		t.len--
	}
	return deleted, ok
}

// DeleteMax removes the largest item in the tree and returns it.
// If no such item exists, returns nil.
func (t *LLRBTree[T]) DeleteMax() (deleted T, ok bool) {
	t.root, deleted, ok = t.deleteMax(t.root)
	if t.root != nil {
		t.root.color = _black
	}
	if ok {
		t.len--
	}
	return deleted, ok
}

// Delete removes an item equal to the passed in item from the tree, returning
// it.  If no such item exists, returns nil.
func (t *LLRBTree[T]) Delete(item T) (deleted T, ok bool) {
	t.root, deleted, ok = t.delete(t.root, item)
	if t.root != nil {
		t.root.color = _black
	}
	if ok {
		t.len--
	}
	return deleted, ok
}

// Clear removes all items from the LLRB-Tree.
func (t *LLRBTree[T]) Clear() {
	t.root = nil
	t.len = 0
}

// Len returns the number of items currently in the tree.
func (t *LLRBTree[T]) Len() int {
	return t.len
}

// AscendRange calls the iterator for every value in the tree within the range
// [greaterOrEqual, lessThan), until iterator returns false.
func (t *LLRBTree[T]) AscendRange(greaterOrEqual, lessThan T, iter IterFunc[T]) {
	t.iterate(t.root, false,
		nullItem[T]{item: greaterOrEqual, valid: true},
		nullItem[T]{item: lessThan, valid: true},
		iter)
}

// AscendLessThan calls the iterator for every value in the tree within the range
// [first, pivot), until iterator returns false.
func (t *LLRBTree[T]) AscendLessThan(pivot T, iter IterFunc[T]) {
	t.iterate(t.root, false,
		nullItem[T]{valid: false},
		nullItem[T]{item: pivot, valid: true},
		iter)
}

// AscendGreaterOrEqual calls the iterator for every value in the tree within
// the range [pivot, last], until iterator returns false.
func (t *LLRBTree[T]) AscendGreaterOrEqual(pivot T, iter IterFunc[T]) {
	t.iterate(t.root, false,
		nullItem[T]{item: pivot, valid: true},
		nullItem[T]{valid: false},
		iter)
}

// Ascend calls the iterator for every value in the tree within the range
// [first, last], until iterator returns false.
func (t *LLRBTree[T]) Ascend(iter IterFunc[T]) {
	t.iterate(t.root, false,
		nullItem[T]{valid: false},
		nullItem[T]{valid: false},
		iter)
}

// DescendRange calls the iterator for every value in the tree within the range
// [lessOrEqual, greaterThan), until iterator returns false.
func (t *LLRBTree[T]) DescendRange(lessOrEqual, greaterThan T, iter IterFunc[T]) {
	t.iterate(t.root, true,
		nullItem[T]{item: lessOrEqual, valid: true},
		nullItem[T]{item: greaterThan, valid: true},
		iter)
}

// DescendLessOrEqual calls the iterator for every value in the tree within the range
// [pivot, first], until iterator returns false.
func (t *LLRBTree[T]) DescendLessOrEqual(pivot T, iter IterFunc[T]) {
	t.iterate(t.root, true,
		nullItem[T]{item: pivot, valid: true},
		nullItem[T]{valid: false},
		iter)
}

// DescendGreaterThan calls the iterator for every value in the tree within
// the range [last, pivot), until iterator returns false.
func (t *LLRBTree[T]) DescendGreaterThan(pivot T, iter IterFunc[T]) {
	t.iterate(t.root, true,
		nullItem[T]{valid: false},
		nullItem[T]{item: pivot, valid: true},
		iter)
}

// Descend calls the iterator for every value in the tree within the range
// [last, first], until iterator returns false.
func (t *LLRBTree[T]) Descend(iter IterFunc[T]) {
	t.iterate(t.root, true,
		nullItem[T]{valid: false},
		nullItem[T]{valid: false},
		iter)
}

func (t *LLRBTree[T]) deleteMin(h *node[T]) (_ *node[T], deleted T, ok bool) {
	if h == nil {
		return nil, zero[T](), false
	}

	if h.left == nil {
		return nil, h.item, true
	}

	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}

	h.left, deleted, ok = t.deleteMin(h.left)

	return fixUp(h), deleted, ok
}

func (t *LLRBTree[T]) deleteMax(h *node[T]) (_ *node[T], deleted T, ok bool) {
	if h == nil {
		return nil, zero[T](), false
	}

	if isRed(h.left) {
		h = rotateRight(h)
	}

	if h.right == nil {
		return nil, h.item, true
	}

	if !isRed(h.right) && !isRed(h.right.left) {
		h = moveRedRight(h)
	}

	h.right, deleted, ok = t.deleteMax(h.right)

	return fixUp(h), deleted, ok
}

func (t *LLRBTree[T]) delete(h *node[T], item T) (_ *node[T], deleted T, ok bool) {
	if h == nil {
		return nil, zero[T](), false
	}
	if t.compare(item, h.item) < 0 {
		if h.left == nil {
			return h, zero[T](), false
		}
		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, deleted, ok = t.delete(h.left, item)
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
		}
		if t.compare(item, h.item) == 0 && h.right == nil {
			return nil, h.item, true
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
		}
		if t.compare(item, h.item) == 0 {
			var rightMin T
			h.right, rightMin, _ = t.deleteMin(h.right)
			deleted, h.item = h.item, rightMin
			ok = true
		} else {
			h.right, deleted, ok = t.delete(h.right, item)
		}
	}

	return fixUp(h), deleted, ok
}

func (t *LLRBTree[T]) insert(h *node[T], item T) (_ *node[T], prev T, exist bool) {
	if h == nil {
		return newNode(item), zero[T](), false
	}

	cmp := t.compare(item, h.item)
	if cmp == 0 {
		prev = h.item
		exist = true
		h.item = item
	} else if cmp < 0 {
		h.left, prev, exist = t.insert(h.left, item)
	} else {
		h.right, prev, exist = t.insert(h.right, item)
	}

	return fixUp(h), prev, exist
}

type nullItem[T any] struct {
	item  T
	valid bool
}

func (t *LLRBTree[T]) iterate(
	h *node[T],
	desc bool,
	start, end nullItem[T],
	iter IterFunc[T],
) bool {
	if h == nil {
		return true
	}

	if !desc {
		if end.valid && t.compare(h.item, end.item) >= 0 {
			return t.iterate(h.left, desc, start, end, iter)
		}
		if start.valid && t.compare(h.item, start.item) < 0 {
			return t.iterate(h.right, desc, start, end, iter)
		}
		if !t.iterate(h.left, desc, start, end, iter) {
			return false
		}
		if !iter(h.item) {
			return false
		}
		return t.iterate(h.right, desc, start, end, iter)
	} else {
		if end.valid && t.compare(h.item, end.item) <= 0 {
			return t.iterate(h.right, desc, start, end, iter)
		}
		if start.valid && t.compare(h.item, start.item) > 0 {
			return t.iterate(h.left, desc, start, end, iter)
		}
		if !t.iterate(h.right, desc, start, end, iter) {
			return false
		}
		if !iter(h.item) {
			return false
		}
		return t.iterate(h.left, desc, start, end, iter)
	}
}

func newNode[T any](item T) *node[T] {
	return &node[T]{
		item:  item,
		color: _red,
	}
}

func rotateLeft[T any](h *node[T]) *node[T] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = _red
	return x
}

func rotateRight[T any](h *node[T]) *node[T] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = _red
	return x
}

func colorFlip[T any](h *node[T]) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func isRed[T any](h *node[T]) bool {
	if h == nil {
		return false
	}
	return h.color
}

func fixUp[T any](h *node[T]) *node[T] {
	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		colorFlip(h)
	}
	return h
}

func moveRedLeft[T any](h *node[T]) *node[T] {
	colorFlip(h)
	if isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
		colorFlip(h)
	}
	return h
}

func moveRedRight[T any](h *node[T]) *node[T] {
	colorFlip(h)
	if isRed(h.left.left) {
		h = rotateRight(h)
		colorFlip(h)
	}
	return h
}

func zero[T any]() T {
	var zero T
	return zero
}
