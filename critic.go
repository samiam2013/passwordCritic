package passwordcritic

import (
	"math"
)

// PassCandidate holds relevant information computed about the password
type PassCandidate struct {
	StringVal   string
	cardinality int
	H           float32
	ErrorVal    error
}

// Load sets the values for the candidate password (entropy, error info)
func (p *PassCandidate) Load(s string) {
	p.StringVal = s
	p.H = p.Entropy()
}

// Entropy takes in a string and gives you a float32 entropy (w/ variety/freq of chars)
func (p *PassCandidate) Entropy() float32 {
	occurrences := charOccurrences(p.StringVal)
	p.cardinality = len(occurrences)
	probabilities := charProbabilites(p.StringVal, occurrences)
	return entropy(probabilities)
}

// charOccurrences maps the frequency of characters for the entropy calculation later
func charOccurrences(text string) map[rune]int {
	occurrences := map[rune]int{}
	for _, char := range text {
		if _, ok := occurrences[char]; !ok {
			occurrences[char] = 1
		} else {
			occurrences[char]++
		}
	}
	return occurrences
}

// Calculate the probability of occurrence of each character
func charProbabilites(text string, occurrences map[rune]int) map[rune]float32 {
	textLength := float32(len(text))
	probabilities := map[rune]float32{}
	for _, char := range text {
		probabilities[char] = float32(occurrences[char]) / textLength
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
