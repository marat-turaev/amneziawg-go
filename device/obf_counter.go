package device

import (
	"encoding/binary"
	"errors"
	"sync/atomic"
)

func newCounterObf(val string) (obf, error) {
	if val != "" {
		return nil, errors.New("counter tag takes no arguments")
	}
	return &counterObf{}, nil
}

type counterObf struct {
	counter atomic.Uint32
}

func (o *counterObf) Obfuscate(dst, src []byte) {
	n := o.counter.Add(1)
	binary.BigEndian.PutUint32(dst, n)
}

func (o *counterObf) Deobfuscate(dst, src []byte) bool {
	return len(src) >= 4
}

func (o *counterObf) ObfuscatedLen(n int) int {
	return 4
}

func (o *counterObf) DeobfuscatedLen(n int) int {
	return 0
}
