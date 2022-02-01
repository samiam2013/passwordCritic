package critic

import (
	"strings"
	"testing"
)

// TestEntropy _
func TestEntropy(t *testing.T) {
	cases := map[string]float32{
		"aaaaaa":                               -1.0,
		"password":                             2.75,
		"p455W0rD!":                            2.947703,
		"correcthorsebatterystaple":            3.363856,
		"thequickbrownfoxjumpedoverthelazydog": 4.447703,
	}

	// create an instance for use of .Entropy()
	pwCand := PassCandidate{
		StringVal:   "",
		Cardinality: 0,
		H:           0.0,
	}

	for pwCase, hExpected := range cases {
		pwCand.StringVal = pwCase
		entropy, err := pwCand.Entropy()
		if err != nil {
			if len(pwCase) < MinLength {
				if !strings.HasPrefix(err.Error(), "password too short") {
					t.Errorf("case '%s' expected password too short, got '%s'", pwCase, err.Error())
				}
			} else if !strings.HasPrefix(err.Error(), "password is homogenous") {
				t.Errorf("case '%s' expected homogeneity error with prefix, got '%s'",
					pwCase, err.Error())
			} else if entropy >= MinEntropy {
				t.Errorf("case '%s' expected no error, got '%s'", pwCase, err.Error())
			}
		}
		if entropy != hExpected {
			t.Errorf("case '%s' expected entropy %+v; got %+v", pwCase, hExpected, entropy)
		}
	}
}
