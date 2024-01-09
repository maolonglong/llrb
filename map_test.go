// Copyright 2024 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLLRBMap(t *testing.T) {
	assert := assert.New(t)

	m := NewMap[string, string]()

	m.Set("foo", "bar")
	assert.Equal(1, m.Len())
	v, ok := m.Get("foo")
	assert.Equal("bar", v)
	assert.True(ok)
	v, ok = m.Get("baz")
	assert.Zero(v)
	assert.False(ok)
	assert.Equal(1, m.Len())
	assert.True(m.Has("foo"))
	assert.False(m.Has("baz"))

	prev, exist := m.Set("foo", "baz")
	assert.True(exist)
	assert.Equal("bar", prev)
	v, ok = m.Get("foo")
	assert.Equal("baz", v)
	assert.True(ok)
	assert.Equal(1, m.Len())

	v, ok = m.Delete("foo")
	assert.Equal("baz", v)
	assert.True(ok)
	assert.Equal(0, m.Len())

	v, ok = m.Delete("foo")
	assert.Zero(v)
	assert.False(ok)
	assert.Equal(0, m.Len())

	m2 := NewMap[int, struct{}]()
	m2.Set(3, struct{}{})
	m2.Set(5, struct{}{})
	m2.Set(1, struct{}{})
	m2.Set(4, struct{}{})
	m2.Set(2, struct{}{})
	var a []int
	m2.Range(func(key int, _ struct{}) bool {
		a = append(a, key)
		return true
	})
	assert.Equal([]int{1, 2, 3, 4, 5}, a)

	m2.Clear()
	assert.Equal(0, m.Len())
}
