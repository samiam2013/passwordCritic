package main

import (
	"fmt"
	"math"
)

// PassCandidate holds relevant information computed about the password
type PassCandidate struct {
	StringVal   string
	cardinality int
	H           float32
	ErrorVal    error
}

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

// Load sets the values for the candidate password (entropy, error info)
func (p *PassCandidate) Load(s string) {
	p.StringVal = s
	p.H = p.Entropy()
}

/* Entropy takes in a string and gives you a float32 entropy
 * value based on variety of characters */
func (p *PassCandidate) Entropy() float32 {
	occurences := charOccurences(p.StringVal)
	p.cardinality = len(occurences)
	probabilities := charProbabilites(p.StringVal, occurences)
	return entropy(probabilities)
}

// charOccurences maps the frequency of characters for the entropy calculation later
func charOccurences(text string) map[rune]int {
	occurences := map[rune]int{}
	for _, char := range text {
		if _, ok := occurences[char]; !ok {
			occurences[char] = 1
		} else {
			occurences[char]++
		}
	}
	return occurences
}

// Calculate the probability of occurrence of each character
func charProbabilites(text string, occurences map[rune]int) map[rune]float32 {
	textLength := float32(len(text))
	probabilities := map[rune]float32{}
	for _, char := range text {
		probabilities[char] = float32(occurences[char]) / textLength
	}
	return probabilities
}

/* Calculate the entropy using the equation
 * H = Σp(i)log_2(1/p(i)) */
func entropy(probabilities map[rune]float32) float32 {
	var h float64 = 0.0
	for _, probability := range probabilities {
		h += float64(probability) * float64(math.Log2(float64(1.0/probability)))
	}
	return float32(h)
}
