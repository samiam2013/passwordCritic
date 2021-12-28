package critic

import (
	"reflect"
	"testing"

	"github.com/samiam2013/passwordcritic/types"
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
	}

	for pwCase, hExpected := range cases {
		pwCand.StringVal = pwCase
		entropy, err := pwCand.Entropy()
		if err != nil && entropy >= 4.0 {
			t.Errorf("case '%s' expected no error, got '%s'", pwCase, err.Error())
		} else if entropy < 4.0 {
			if reflect.TypeOf(err) != reflect.TypeOf(&types.HomogeneityError{}) {
				t.Errorf("case '%s' entropy < 4.0 (%f), error type '%s'", pwCase, entropy,
					reflect.TypeOf(err))
			}
		}
		if entropy != hExpected {
			t.Errorf("case '%s' expected entropy %+v; got %+v", pwCase, entropy, hExpected)
		}
	}
}
