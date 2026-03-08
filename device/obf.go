package device

import (
	"errors"
	"fmt"
	"strings"
)

type obfBuilder func(val string) (obf, error)

var obfBuilders = map[string]obfBuilder{
	"b":  newBytesObf,
	"c":  newCounterObf,
	"t":  newTimestampObf,
	"r":  newRandObf,
	"rc": newRandCharObf,
	"rd": newRandDigitsObf,
	"d":  newDataObf,
	"ds": newDataStringObf,
	"dz": newDataSizeObf,
}

type obf interface {
	Obfuscate(dst, src []byte)
	Deobfuscate(dst, src []byte) bool
	ObfuscatedLen(srcLen int) int
	DeobfuscatedLen(srcLen int) int
}

type obfChain struct {
	Spec string
	obfs []obf
}

func newObfChain(spec string) (*obfChain, error) {
	var (
		obfs []obf
		errs []error
	)

	remaining := spec[:]
	for {
		start := strings.IndexByte(remaining, '<')
		if start == -1 {
			if strings.TrimSpace(remaining) != "" {
				errs = append(errs, fmt.Errorf("unexpected text outside tags: %q", strings.TrimSpace(remaining)))
			}
			break
		}
		if strings.TrimSpace(remaining[:start]) != "" {
			errs = append(errs, fmt.Errorf("unexpected text outside tags: %q", strings.TrimSpace(remaining[:start])))
		}

		end := strings.IndexByte(remaining[start:], '>')
		if end == -1 {
			return nil, errors.New("missing enclosing >")
		}
		end += start

		tag := remaining[start+1 : end]
		parts := strings.Fields(tag)
		if len(parts) == 0 {
			errs = append(errs, errors.New("empty tag"))
			remaining = remaining[end+1:]
			continue
		}

		key := parts[0]
		builder, ok := obfBuilders[key]
		if !ok {
			errs = append(errs, fmt.Errorf("unknown tag <%s>", key))
			remaining = remaining[end+1:]
			continue
		}

		val := ""
		if len(parts) > 1 {
			val = parts[1]
		}

		o, err := builder(val)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to build <%s>: %w", key, err))
			remaining = remaining[end+1:]
			continue
		}

		obfs = append(obfs, o)
		remaining = remaining[end+1:]
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	if len(obfs) == 0 {
		return nil, errors.New("empty obfuscation chain")
	}

	return &obfChain{
		Spec: spec,
		obfs: obfs,
	}, nil
}

func (c *obfChain) Obfuscate(dst, src []byte) {
	written := 0
	for _, o := range c.obfs {
		obfLen := o.ObfuscatedLen(len(src))
		if obfLen < 0 || written+obfLen > len(dst) {
			return
		}
		o.Obfuscate(dst[written:written+obfLen], src)
		written += obfLen
	}
}

func (c *obfChain) Deobfuscate(dst, src []byte) bool {
	dynamicLen := len(src) - c.ObfuscatedLen(0)
	if dynamicLen < 0 {
		return false
	}

	written, read := 0, 0

	for _, o := range c.obfs {
		deobfLen := o.DeobfuscatedLen(dynamicLen)
		obfLen := o.ObfuscatedLen(deobfLen)
		if deobfLen < 0 || obfLen < 0 || written+deobfLen > len(dst) || read+obfLen > len(src) {
			return false
		}

		if !o.Deobfuscate(dst[written:written+deobfLen], src[read:read+obfLen]) {
			return false
		}

		written += deobfLen
		read += obfLen
	}

	return true
}

func (c *obfChain) ObfuscatedLen(n int) int {
	total := 0
	for _, o := range c.obfs {
		total += o.ObfuscatedLen(n)
	}
	return total
}

func (c *obfChain) DeobfuscatedLen(n int) int {
	dynamicLen := n - c.ObfuscatedLen(0)

	total := 0
	for _, o := range c.obfs {
		total += o.DeobfuscatedLen(dynamicLen)
	}
	return total
}
