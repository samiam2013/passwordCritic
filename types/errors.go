package types

import "fmt"

type HomogeneityError struct {
	Cardinality       int
	LowestProbability float32
}

func (h *HomogeneityError) Error() string {
	return fmt.Sprintf("password is homogenous with cardinality %d "+
		"and lowest probability of occurence for a letter (%f)[value 0-1]",
		h.Cardinality, h.LowestProbability)
}

type TooCommonError struct {
	Occurrences int
}

func (t *TooCommonError) Error() string {
	return fmt.Sprintf("password appears at least %d times in the 'too commmon' "+
		"passwords list.", t.Occurrences)
}

// TODO hit match in bloom filter error