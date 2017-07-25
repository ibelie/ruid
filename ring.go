// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"crypto/md5"
	"fmt"
	"sort"
)

const (
	VIRTUAL_NODES  = 50
	DEFAULT_WEIGHT = 1
)

type Ring struct {
	sorted  IDSlice
	weights map[string]int
	ring    map[ID]string
}

func NewRing(nodes ...string) *Ring {
	ring := &Ring{weights: make(map[string]int)}
	for _, node := range nodes {
		ring.weights[node] = DEFAULT_WEIGHT
	}
	ring.circle()
	return ring
}

func WeightedRing(weights map[string]int) *Ring {
	ring := &Ring{weights: weights}
	ring.circle()
	return ring
}

func (r *Ring) Update(weights map[string]int) {
	changed := false
	for node, weight := range weights {
		if w, ok := r.weights[node]; !ok || w != weight {
			r.weights[node] = weight
			changed = true
		}
	}
	if changed {
		r.circle()
	}
}

func (r *Ring) Append(nodes ...string) {
	for _, node := range nodes {
		r.weights[node] = DEFAULT_WEIGHT
	}
	r.circle()
}

func (r *Ring) Remove(nodes ...string) {
	for _, node := range nodes {
		delete(r.weights, node)
	}
	r.circle()
}

func (r *Ring) Get(key ID) (node string, ok bool) {
	if len(r.ring) <= 0 {
		return "", false
	}
	pos := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= key })
	if pos == len(r.sorted) {
		pos = 0
	}
	return r.ring[r.sorted[pos]], true
}

func (r *Ring) circle() {
	virtual := VIRTUAL_NODES
	total := 0
	for _, weight := range r.weights {
		total += weight
		if virtual < weight {
			virtual = weight
		}
	}

	r.sorted = nil
	r.ring = make(map[ID]string)
	for node, weight := range r.weights {
		factor := len(r.weights) * weight * virtual / total
		if factor < 1 {
			factor = 1
		}
		for i := 0; i < int(factor); i++ {
			d := md5.Sum([]byte(fmt.Sprintf("%s-%d", node, i)))
			for j := 0; j < 16; j += 8 {
				key := (ID(d[j+7]) << 56) |
					(ID(d[j+6]) << 48) |
					(ID(d[j+5]) << 40) |
					(ID(d[j+4]) << 32) |
					(ID(d[j+3]) << 24) |
					(ID(d[j+2]) << 16) |
					(ID(d[j+1]) << 8) |
					ID(d[j])
				r.ring[key] = node
				r.sorted = append(r.sorted, key)
			}
		}
	}
	r.sorted.Sort()
}
