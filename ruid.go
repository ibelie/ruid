// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package ruid

import (
	"net"
	"sync"
	"time"

	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"github.com/ibelie/tygo"
)

type ID interface {
	Lt(ID) bool
	Ge(ID) bool
	Hash() ID
	String() string
	Nonzero() bool
	ByteSize() int
	Serialize(*tygo.ProtoBuf)
}

type Ident interface {
	New() ID
	Zero() ID
	Deserialize(*tygo.ProtoBuf) (ID, error)
	GetIDs([]byte) []ID
}

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

const ZERO RUID = 0

const EncodeRUID = "-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

var RUIDEncoding = base64.NewEncoding(EncodeRUID).WithPadding(base64.NoPadding)

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
	return RUIDEncoding.EncodeToString(bytes)
}

func FromString(s string) (RUID, error) {
	if bytes, err := RUIDEncoding.DecodeString(s); err != nil {
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

func (r RUID) Hash() ID {
	return r
}

func (r RUID) Lt(o ID) bool {
	return r < o.(RUID)
}

func (r RUID) Ge(o ID) bool {
	return r >= o.(RUID)
}

func (r RUID) Nonzero() bool {
	return r != ZERO
}

func (r RUID) ByteSize() (size int) {
	return 8
}

func (r RUID) Serialize(output *tygo.ProtoBuf) {
	output.WriteFixed64(uint64(r))
}

func (r *RUID) Deserialize(input *tygo.ProtoBuf) (err error) {
	x, err := input.ReadFixed64()
	*r = RUID(x)
	return
}

type RUIdentity int

var RUIdent RUIdentity = 0

func (_ RUIdentity) New() ID {
	return New()
}

func (_ RUIdentity) Zero() ID {
	return ZERO
}

func (_ RUIdentity) Deserialize(input *tygo.ProtoBuf) (r ID, err error) {
	x, err := input.ReadFixed64()
	r = RUID(x)
	return
}

func (_ RUIdentity) GetIDs(bytes []byte) (ids []ID) {
	for i, n := 0, len(bytes)-8; i <= n; i += 8 {
		ids = append(ids, RUID(binary.LittleEndian.Uint64(bytes[i:])))
	}
	return
}
