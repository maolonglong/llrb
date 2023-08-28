// Copyright 2023 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLLRBSet(t *testing.T) {
	assert := assert.New(t)

	s := NewSet[int]()
	{
		a := rnd(1000, 10)

		for _, x := range a {
			s.Insert(x)
		}
		assert.LessOrEqual(s.Len(), 10)
	}

	{
		s.Clear()
		assert.Equal(0, s.Len())

		a := []int{3, 5, 1, 4, 2}

		for _, x := range a {
			assert.False(s.Insert(x))
		}
		for _, x := range a {
			assert.True(s.Insert(x))
		}
		assert.Equal(5, s.Len())

		for _, x := range a {
			assert.True(s.Has(x))
		}

		var aa []int
		s.Range(func(x int) bool {
			aa = append(aa, x)
			return true
		})
		assert.Equal([]int{1, 2, 3, 4, 5}, aa)

		for _, x := range a {
			assert.True(s.Delete(x))
		}
		assert.Equal(0, s.Len())

		for _, x := range a {
			assert.False(s.Delete(x))
		}
	}
}
