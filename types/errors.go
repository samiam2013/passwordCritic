package types

import (
	"fmt"
	"log"
	"strings"
)

// HomogeneityError creates a switchable type of error for low entropy
type HomogeneityError struct {
	Cardinality       int
	LowestProbability float32
}

func (h *HomogeneityError) Error() string {
	return fmt.Sprintf("password is homogenous with cardinality %d "+
		"and lowest probability of occurrence for a letter (%f)[value 0-1]",
		h.Cardinality, h.LowestProbability)
}

// TooCommonError creates a switchable type of error for common passwords
type TooCommonError struct {
	Occurrences int
	Level       string
}

func (t *TooCommonError) Error() string {
	switch t.Level {
	case "info":
		// continue
	case "warn":
		// continue
	case "error":
		// continue
	default:
		log.Fatal("TooCommonError.Level must be 'info', 'warn', or 'error'")
	}
	return fmt.Sprintf("%s: password appears at least %d times in the 'too commmon' "+
		"passwords list.", strings.ToUpper(t.Level), t.Occurrences)
}

// tie these types to the error interface
var _, _ error = &HomogeneityError{}, &TooCommonError{}
