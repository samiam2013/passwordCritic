package critic

import (
	"fmt"
	"math"

	"github.com/samiam2013/passwordcritic/types"
)

// PassCandidate holds relevant information computed about the password
type PassCandidate struct {
	StringVal   string
	Cardinality int
	H           float64
}

// MinEntropy defines the lowest value for throwing a Homogeneity Error
const MinEntropy = 3.0

// MinLength defines the shortest password allowed
const MinLength = 8

// Entropy returns a float calculated using variety and frequency of characters
func (p *PassCandidate) Entropy() (float64, error) {
	if len(p.StringVal) < MinLength {
		err := fmt.Errorf("password too short, minimum %d characters", MinLength)
		p.H = -1.0
		return p.H, err
	}
	occurrences := charOccurCount(p.StringVal)
	p.Cardinality = len(occurrences)
	probabilities := charProbabilites(p.StringVal, occurrences)
	h := entropy(probabilities)
	p.H = h
	if h < MinEntropy {
		return h,
			&types.HomogeneityError{
				Cardinality:       p.Cardinality,
				LowestProbability: minMap(probabilities),
			}
	}
	return h, nil
}

// minMap gives the min number map for custom errors
func minMap(runeMap map[rune]float64) float64 {
	minFloat := math.MaxFloat64
	for _, candidate := range runeMap {
		if candidate < minFloat {
			minFloat = candidate
		}
	}
	return minFloat
}

// charOccurCount maps the frequency of characters for the entropy calculation later
func charOccurCount(text string) map[rune]int {
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
func charProbabilites(text string, occurrences map[rune]int) map[rune]float64 {
	textLength := float64(len(text))
	probabilities := map[rune]float64{}
	for _, char := range text {
		probabilities[char] = float64(occurrences[char]) / textLength
	}
	return probabilities
}

// Calculate the entropy using the equation
// h = Σ prob(i) * log₂(1/prob(i))
func entropy(probabilities map[rune]float64) float64 {
	h := 0.0
	for _, probability := range probabilities {
		h += float64(probability) * math.Log2(1.0/probability)
	}
	return h
}
