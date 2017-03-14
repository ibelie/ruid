// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"fmt"
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
