// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	HARDWARE_MASK    = 0x7FF
	HARDWARE_BITS    = 11
	SEQUENCE_MASK    = 0xFFF
	SEQUENCE_BITS    = 12
	HARDWARE_OFFSET  = SEQUENCE_BITS
	TIMESTAMP_OFFSET = SEQUENCE_BITS + HARDWARE_OFFSET
)

var (
	initial  sync.Once
	sequence uint64
	hardware uint64
	lastTime uint64
)

// RUID: Recently Unique Identifier
// <-             timestamp            -> <- hardware -> <- sequence ->
// 00000000 00000000 00000000 00000000 0 - 000 00000000 - 0000 00000000

type RUID uint64

func New() RUID {
	initial.Do(func() {
		bytes := make([]byte, 2)
		rand.Read(bytes)
		sequence = uint64(binary.BigEndian.Uint16(bytes))

		if interfaces, err := net.Interfaces(); err == nil {
			for _, iface := range interfaces {
				if len(iface.HardwareAddr) >= 2 {
					copy(bytes, iface.HardwareAddr)
					break
				}
			}
		} else {
			rand.Read(bytes)
		}
		hardware = uint64(binary.BigEndian.Uint16(bytes)&HARDWARE_MASK) << HARDWARE_OFFSET
	})

	currTime := uint64(time.Now().UnixNano() / 1e6)
	if atomic.CompareAndSwapUint64(&lastTime, lastTime, currTime) {
		atomic.AddUint64(&sequence, 1)
	}

	return RUID(hardware | (sequence & SEQUENCE_MASK) | (currTime << TIMESTAMP_OFFSET))
}

func (r RUID) String() string {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(r))
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func FromString(s string) (RUID, error) {
	if bytes, err := base64.RawURLEncoding.DecodeString(s); err != nil {
		return 0, err
	} else {
		return RUID(binary.BigEndian.Uint64(bytes)), nil
	}
}

func (r RUID) Bytes() []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(r))
	return bytes
}

func FromBytes(b []byte) RUID {
	return RUID(binary.BigEndian.Uint64(b))
}
