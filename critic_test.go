package passwordcritic

import (
	"testing"
)

// TestEntropy _
func TestEntropy(t *testing.T) {
	cases := map[string]float32{
		"aaaaaa":                               0.0,
		"password":                             2.75,
		"p455W0rD!":                            2.947703,
		"correcthorsebatterystaple":            3.363856,
		"thequickbrownfoxjumpedoverthelazydog": 4.447703,
	}

	// create an instance for use of .Entropy()
	pwCand := PassCandidate{
		StringVal:   "",
		cardinality: 0,
		H:           0.0,
		ErrorVal:    nil,
	}

	for pwCase, hExpected := range cases {
		pwCand.Load(pwCase)
		entropy := pwCand.Entropy()

		if entropy != hExpected {
			t.Errorf("case '%s' expected entropy %+v; got %+v", pwCase, entropy, hExpected)
		}
	}
}
