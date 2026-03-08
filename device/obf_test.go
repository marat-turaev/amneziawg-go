package device

import "testing"

func TestNewObfChainRejectsTextOutsideTags(t *testing.T) {
	t.Parallel()

	cases := []string{
		"garbage",
		"<b 0xf6ab3267fa> trailing",
		"leading <b 0xf6ab3267fa>",
		"   ",
	}

	for _, spec := range cases {
		spec := spec
		t.Run(spec, func(t *testing.T) {
			t.Parallel()

			if _, err := newObfChain(spec); err == nil {
				t.Fatalf("expected %q to be rejected", spec)
			}
		})
	}
}

func TestNewObfChainAcceptsAdjacentTags(t *testing.T) {
	t.Parallel()

	chain, err := newObfChain("<b 0xf6ab3267fa><c><r 10>")
	if err != nil {
		t.Fatalf("expected valid chain, got error: %v", err)
	}
	if got := len(chain.obfs); got != 3 {
		t.Fatalf("expected 3 obfuscators, got %d", got)
	}
}
