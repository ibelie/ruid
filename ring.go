// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"crypto/md5"
	"fmt"
	"log"
	"sort"
)

const (
	VIRTUAL_NODES  = 50
	DEFAULT_WEIGHT = 1
)

type Ring struct {
	ident   Ident
	sorted  IDSlice
	weights map[string]int
	ring    map[ID]string
}

func NewRing(ident Ident, nodes ...string) *Ring {
	ring := &Ring{ident: ident, weights: make(map[string]int)}
	for _, node := range nodes {
		ring.weights[node] = DEFAULT_WEIGHT
	}
	ring.circle()
	return ring
}

func WeightedRing(ident Ident, weights map[string]int) *Ring {
	ring := &Ring{ident: ident, weights: weights}
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
	hash := key.Hash()
	pos := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i].Ge(hash) })
	if pos == len(r.sorted) {
		pos = 0
	}
	return r.ring[r.sorted[pos]], true
}

func (r *Ring) Key(node string) ID {
	bytes := md5.Sum([]byte(node))
	for _, key := range r.ident.GetIDs(bytes[:]) {
		return key
	}
	return nil
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

	r.ring = make(map[ID]string)
	for node, weight := range r.weights {
		for i := 0; i < int(len(r.weights)*weight*virtual/total); i++ {
			bytes := md5.Sum([]byte(fmt.Sprintf("%s-%d", node, i)))
			for _, key := range r.ident.GetIDs(bytes[:]) {
				r.ring[key] = node
			}
		}
	}

	conflict := make(map[ID]string)
	for node, _ := range r.weights {
		hash := r.Key(node).Hash()
		if n, ok := conflict[hash]; ok {
			log.Fatalf("[RUID] Ring nodes conflict: %q %q", n, node)
		} else {
			conflict[hash] = node
			r.ring[hash] = node
		}
	}

	r.sorted = nil
	for key, _ := range r.ring {
		r.sorted = append(r.sorted, key)
	}
	r.sorted.Sort()
}
