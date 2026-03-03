package device

import (
	"crypto/rand"
	"errors"
	"strconv"
	"unicode"
)

const chars52 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func newRandCharObf(val string) (obf, error) {
	length, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("length must be non-negative")
	}

	return &randCharObf{
		length: length,
	}, nil
}

type randCharObf struct {
	length int
}

func (o *randCharObf) Obfuscate(dst, src []byte) {
	if o.length > len(dst) {
		return
	}
	_, _ = rand.Read(dst[:o.length])
	for i := range dst[:o.length] {
		dst[i] = chars52[dst[i]%52]
	}
}

func (o *randCharObf) Deobfuscate(dst, src []byte) bool {
	if o.length > len(src) {
		return false
	}
	for _, b := range src[:o.length] {
		if !unicode.IsLetter(rune(b)) {
			return false
		}
	}
	return true
}

func (o *randCharObf) ObfuscatedLen(n int) int {
	return o.length
}

func (o *randCharObf) DeobfuscatedLen(n int) int {
	return 0
}
