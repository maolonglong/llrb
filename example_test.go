// Copyright 2024 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llrb_test

import (
	"fmt"

	"github.com/maolonglong/llrb"
)

func ExampleLLRBTree() {
	// Create a new LLRBTree
	tree := llrb.NewOrdered[int]()

	// Insert some items into the tree
	tree.ReplaceOrInsert(5)
	tree.ReplaceOrInsert(2)
	tree.ReplaceOrInsert(7)
	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(4)

	// Iterate over the tree in ascending order
	tree.Ascend(func(item int) bool {
		fmt.Println(item)
		return true
	})
	// Output:
	// 1
	// 2
	// 4
	// 5
	// 7
}

func ExampleLLRBMap() {
	m := llrb.NewMap[int, string]()

	m.Set(1, "apple")
	m.Set(2, "banana")
	m.Set(3, "cherry")

	v, ok := m.Get(2)
	if ok {
		fmt.Println(v)
	}

	m.Delete(3)

	m.Range(func(key int, value string) bool {
		fmt.Println(key, value)
		return true
	})
	// Output:
	// banana
	// 1 apple
	// 2 banana
}

func ExampleLLRBSet() {
	s := llrb.NewSet[int]()

	s.Insert(1)
	s.Insert(2)
	s.Insert(2)
	s.Insert(3)

	fmt.Println(s.Len())

	fmt.Println(s.Has(2))

	s.Delete(2)
	fmt.Println(s.Len())

	s.Range(func(x int) bool {
		fmt.Println(x)
		return true
	})
	// Output:
	// 3
	// true
	// 2
	// 1
	// 3
}
