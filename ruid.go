// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	TIMESTAMP_MASK  = 0x1FFFFFFFFFF
	TIMESTAMP_BITS  = 41
	HARDWARE_MASK   = 0x7FF
	HARDWARE_BITS   = 11
	HARDWARE_OFFSET = TIMESTAMP_BITS
	SEQUENCE_OFFSET = HARDWARE_BITS + HARDWARE_OFFSET
)

var (
	initial  sync.Once
	sequence uint64
	hardware uint64
)

// RUID: Recently Unique Identifier
// <- sequence -> <- hardware -> <-                 timestamp                 ->
// 00000000 0000 - 0000 0000000 - 0 00000000 00000000 00000000 00000000 00000000

type RUID uint64

func New() RUID {
	initial.Do(func() {
		bytes := make([]byte, 2)
		rand.Read(bytes)
		sequence = uint64(binary.LittleEndian.Uint16(bytes))

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
		hardware = uint64(binary.LittleEndian.Uint16(bytes)&HARDWARE_MASK) << HARDWARE_OFFSET
	})

	sequence++
	timestamp := uint64(time.Now().UnixNano() / 1e6)
	return RUID(hardware | (sequence << SEQUENCE_OFFSET) | (timestamp & TIMESTAMP_MASK))
}

func (r RUID) String() string {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(r))
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func FromString(s string) (RUID, error) {
	if bytes, err := base64.RawURLEncoding.DecodeString(s); err != nil {
		return 0, err
	} else {
		return RUID(binary.LittleEndian.Uint64(bytes)), nil
	}
}

func (r RUID) Bytes() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(r))
	return bytes
}

func FromBytes(b []byte) RUID {
	return RUID(binary.LittleEndian.Uint64(b))
}

type RUIDSlice []RUID

func (s RUIDSlice) Len() int           { return len(s) }
func (s RUIDSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s RUIDSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RUIDSlice) Sort()              { sort.Sort(s) }
