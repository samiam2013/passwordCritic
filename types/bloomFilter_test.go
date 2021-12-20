package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test items if items may exist into the set
func TestBloomFilter(t *testing.T) {
	assert := require.New(t)
	bf := NewBloom(1024)
	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))
	bf.Add([]byte("sir"))
	bf.Add([]byte("madam"))
	bf.Add([]byte("io"))

	cases := map[string]bool{
		"hello": true,
		"world": true,
		"hi":    false,
	}
	for candidate, expected := range cases {
		gotResult, err := bf.Test([]byte(candidate))
		if err != nil {
			t.Errorf("failed to test candidate '%s' against filter: %s", candidate, err.Error())
		}
		assert.Equal(gotResult, expected)
	}
}
