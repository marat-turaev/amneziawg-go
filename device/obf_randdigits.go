package device

import (
	"crypto/rand"
	"errors"
	"strconv"
	"unicode"
)

const digits10 = "0123456789"

func newRandDigitsObf(val string) (obf, error) {
	length, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("length must be non-negative")
	}

	return &randDigitObf{
		length: length,
	}, nil
}

type randDigitObf struct {
	length int
}

func (o *randDigitObf) Obfuscate(dst, src []byte) {
	if o.length > len(dst) {
		return
	}
	_, _ = rand.Read(dst[:o.length])
	for i := range dst[:o.length] {
		dst[i] = digits10[dst[i]%10]
	}
}

func (o *randDigitObf) Deobfuscate(dst, src []byte) bool {
	if o.length > len(src) {
		return false
	}
	for _, b := range src[:o.length] {
		if !unicode.IsDigit(rune(b)) {
			return false
		}
	}
	return true
}

func (o *randDigitObf) ObfuscatedLen(n int) int {
	return o.length
}

func (o *randDigitObf) DeobfuscatedLen(n int) int {
	return 0
}
