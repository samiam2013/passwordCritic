package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test items if items may exist into the set
func TestBloomFilter(t *testing.T) {
	assert := require.New(t)
	bf := NewBloom(1024)
	inputs := []string{"hello", "world", "hi", "fizz", "buzz", "foo", "bar"}
	controls := []string{"shalom", "gaia", "konnichiwa", "kuohua", "zumbido"}

	for _, input := range inputs {
		err := bf.Add([]byte(input))
		if err != nil {
			t.Error(err)
		}
	}

	for _, input := range inputs {
		gotResult, err := bf.Test([]byte(input))
		if err != nil {
			t.Errorf("failed to test candidate '%s' against filter: %s", input, err.Error())
		}
		assert.True(gotResult)
	}

	for _, control := range controls {
		gotResult, err := bf.Test([]byte(control))
		if err != nil {
			t.Error(err)
		}
		assert.False(gotResult)
	}
}
