// Copyright 2023 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb

import (
	"cmp"
	"math"
	"math/rand"
	"runtime"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLLRBTree(t *testing.T) {
	assert.PanicsWithValue(t, "nil compare", func() {
		_ = NewLLRBTree[int](nil)
	})
}

func TestLLRBTree_insert(t *testing.T) {
	const N = 1000
	const LOOP = 100

	assert := assert.New(t)
	tree := NewLLRBTree[int](cmp.Compare[int])

	a := seq(N)

	for i := 0; i < LOOP; i++ {
		shuffle(a)
		for _, x := range a {
			prev, exist := tree.ReplaceOrInsert(x)
			assert.Zero(prev)
			assert.False(exist)
		}
		assertMaxDepth(t, tree)
		tree.Clear()
	}
}

func TestLLRBTree_random_insert_delete(t *testing.T) {
	assert := assert.New(t)
	tree := NewLLRBTree[int](cmp.Compare[int])

	const N = 10000

	a := rnd(N, 10*N)
	uniqNums := uniq(a)

	insert := func() {
		for _, x := range a {
			_, _ = tree.ReplaceOrInsert(x)
		}
		assert.Equal(uniqNums, tree.Len())
		assertMaxDepth(t, tree)
	}

	insert()
	for _, x := range a {
		item, ok := tree.Get(x)
		assert.Equal(item, x)
		assert.True(ok)
	}

	for _, x := range a {
		_, _ = tree.Delete(x)
	}
	assert.Equal(0, tree.Len())
	assert.Nil(tree.root)

	insert()
	for range a {
		_, _ = tree.DeleteMin()
	}
	assert.Equal(0, tree.Len())
	assert.Nil(tree.root)

	insert()
	for range a {
		_, _ = tree.DeleteMax()
	}
	assert.Equal(0, tree.Len())
	assert.Nil(tree.root)
}

func TestLLRBTree_iterator(t *testing.T) {
	assert := assert.New(t)

	a := seq(100)
	tree := NewLLRBTree[int](cmp.Compare[int])
	for _, x := range a {
		tree.ReplaceOrInsert(x)
	}
	assertMaxDepth(t, tree)

	var collect []int
	tree.Ascend(func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal(a, collect)

	collect = collect[:0]
	tree.AscendGreaterOrEqual(95, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{95, 96, 97, 98, 99, 100}, collect)

	collect = collect[:0]
	tree.AscendLessThan(5, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{1, 2, 3, 4}, collect)

	collect = collect[:0]
	tree.AscendRange(48, 52, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{48, 49, 50, 51}, collect)

	collect = collect[:0]
	tree.Descend(func(x int) bool {
		collect = append(collect, x)
		return true
	})
	slices.Reverse(a)
	assert.Equal(a, collect)

	collect = collect[:0]
	tree.DescendGreaterThan(98, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{100, 99}, collect)

	collect = collect[:0]
	tree.DescendLessOrEqual(5, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{5, 4, 3, 2, 1}, collect)

	collect = collect[:0]
	tree.DescendRange(52, 48, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{52, 51, 50, 49}, collect)

	collect = collect[:0]
	tree.AscendGreaterOrEqual(101, func(x int) bool {
		collect = append(collect, x)
		return true
	})
	assert.Equal([]int{}, collect)
}

func TestLLRBTree_iterator_break(t *testing.T) {
	assert := assert.New(t)

	a := seq(100)
	tree := NewLLRBTree[int](cmp.Compare[int])
	for _, x := range a {
		tree.ReplaceOrInsert(x)
	}

	var collect []int
	tree.Ascend(func(x int) bool {
		collect = append(collect, x)
		return false
	})
	assert.Equal([]int{1}, collect)

	collect = collect[:0]
	tree.Descend(func(x int) bool {
		collect = append(collect, x)
		return false
	})
	assert.Equal([]int{100}, collect)

	collect = collect[:0]
	tree.Ascend(func(x int) bool {
		collect = append(collect, x)
		return x < 5
	})
	assert.Equal([]int{1, 2, 3, 4, 5}, collect)
}

func BenchmarkLLRBTree_insert_random(b *testing.B) {
	const L = 50000
	assert := assert.New(b)

	a := seq(L)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		shuffle(a)
		runtime.GC()
		b.StartTimer()

		t := NewLLRBTree[int](cmp.Compare[int])
		for _, x := range a {
			prev, exist := t.ReplaceOrInsert(x)

			assert.Equal(0, prev)
			assert.False(exist)
		}
		assert.Equal(L, t.Len())
	}
}

func BenchmarkLLRBTree_insert_ascending(b *testing.B) {
	const L = 50000
	assert := assert.New(b)

	a := seq(L)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		runtime.GC()
		b.StartTimer()

		t := NewLLRBTree[int](cmp.Compare[int])
		for _, x := range a {
			prev, exist := t.ReplaceOrInsert(x)

			assert.Equal(0, prev)
			assert.False(exist)
		}
		assert.Equal(L, t.Len())
	}
}

func BenchmarkLLRBTree_get_random(b *testing.B) {
	const L = 50000
	assert := assert.New(b)

	a := seq(L)
	getOps := seq(2 * L)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		shuffle(a)
		shuffle(getOps)
		t := NewLLRBTree[int](cmp.Compare[int])
		for _, x := range a {
			_, _ = t.ReplaceOrInsert(x)
		}
		runtime.GC()
		b.StartTimer()

		for _, x := range getOps {
			item, ok := t.Get(x)
			if x <= L {
				assert.Equal(x, item)
				assert.True(ok)
			} else {
				assert.Equal(0, item)
				assert.False(ok)
			}
		}
	}
}

func BenchmarkLLRBTree_get_ascending(b *testing.B) {
	const L = 50000
	assert := assert.New(b)

	a := seq(L)
	getOps := seq(2 * L)

	t := NewLLRBTree[int](cmp.Compare[int])
	for _, x := range a {
		_, _ = t.ReplaceOrInsert(x)
	}
	runtime.GC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, x := range getOps {
			item, ok := t.Get(x)
			if x <= L {
				assert.Equal(x, item)
				assert.True(ok)
			} else {
				assert.Equal(0, item)
				assert.False(ok)
			}
		}
	}
}

func uniq(a []int) int {
	set := make(map[int]struct{})
	for _, x := range a {
		set[x] = struct{}{}
	}
	return len(set)
}

func assertMaxDepth[T any](tb testing.TB, tree *LLRBTree[T]) {
	tb.Helper()

	assert.LessOrEqual(tb, maxDepth(tree.root), int(2*math.Log2(float64(tree.len)+1)))
}

func maxDepth[T any](h *node[T]) int {
	if h == nil {
		return 0
	}
	return 1 + max(maxDepth(h.left), maxDepth(h.right))
}

func rnd(n, maxValue int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(maxValue)
	}
	return a
}

func seq(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i + 1
	}
	return a
}

func shuffle(a []int) []int {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}
