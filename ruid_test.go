// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"fmt"
	"sort"
	"testing"
)

func TestString(t *testing.T) {
	r1 := New()
	s := r1.String()
	fmt.Println(r1, uint64(r1))
	r2, err := FromString(s)
	if err != err {
		t.Error(err)
	} else if r1 != r2 {
		t.Errorf("FromString error: %v %v", r1, r2)
	}
}

func TestBytes(t *testing.T) {
	r1 := New()
	b := r1.Bytes()
	fmt.Println(r1, b)
	r2 := FromBytes(b)
	if r1 != r2 {
		t.Errorf("FromBytes error: %v %v %v", r1, r2, b)
	}
}

func randomRing(ring *Ring, n int) {
	count := make(map[string]int)
	for i := 0; i < n; i++ {
		id := New()
		node, _ := ring.Get(id)
		if _, exist := count[node]; !exist {
			count[node] = 0
		}
		count[node]++
	}
	var sorted []int
	for _, c := range count {
		sorted = append(sorted, c)
	}
	sort.Ints(sorted)
	for i, c := range sorted {
		fmt.Printf("%d, %d,\n", i, c)
	}
}

func TestRing(t *testing.T) {
	var nodes []string
	weights := make(map[string]int)
	for i := 0; i < 100; i++ {
		node := fmt.Sprintf("node_%d", i)
		nodes = append(nodes, node)
		weights[node] = i
	}
	randomRing(WeightedRing(weights), 100000)
	randomRing(NewRing(nodes...), 10000)
}
