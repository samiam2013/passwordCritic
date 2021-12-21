package types

import "fmt"

/* HomogeneityError creates a switchable type of error for problems
 * relating to the entropy of the password */
type HomogeneityError struct {
	Cardinality       int
	LowestProbability float32
}

func (h *HomogeneityError) Error() string {
	return fmt.Sprintf("password is homogenous with cardinality %d "+
		"and lowest probability of occurence for a letter (%f)[value 0-1]",
		h.Cardinality, h.LowestProbability)
}

/* TooCommonError creates a switchable type of error for problems
 * caused by a hit in the bloom filter for common passwords */
type TooCommonError struct {
	Occurrences int
}

func (t *TooCommonError) Error() string {
	return fmt.Sprintf("password appears at least %d times in the 'too commmon' "+
		"passwords list.", t.Occurrences)
}

// tie these types to the error interface
var _, _ error = &HomogeneityError{}, &TooCommonError{}
