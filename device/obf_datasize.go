package device

import (
	"errors"
	"strconv"
)

func newDataSizeObf(val string) (obf, error) {
	length, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("length must be non-negative")
	}

	return &dataSizeObf{
		length: length,
	}, nil
}

type dataSizeObf struct {
	length int
}

func (o *dataSizeObf) Obfuscate(dst, src []byte) {
	if o.length > len(dst) {
		return
	}
	srcLen := len(src)
	for i := o.length - 1; i >= 0; i-- {
		dst[i] = byte(srcLen & 0xFF)
		srcLen >>= 8
	}
}

func (o *dataSizeObf) Deobfuscate(dst, src []byte) bool {
	return true
}

func (o *dataSizeObf) ObfuscatedLen(srcLen int) int {
	return o.length
}

func (o *dataSizeObf) DeobfuscatedLen(srcLen int) int {
	return 0
}
