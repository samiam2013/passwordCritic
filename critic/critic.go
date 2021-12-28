package critic

import (
	"math"

	"github.com/samiam2013/passwordcritic/types"
)

// PassCandidate holds relevant information computed about the password
type PassCandidate struct {
	StringVal   string
	Cardinality int
	H           float32
}

// Entropy returns a float32 calculated using variety and frequency of characters
func (p *PassCandidate) Entropy() (float32, error) {
	occurrences := charOccurrences(p.StringVal)
	p.Cardinality = len(occurrences)
	probabilities := charProbabilites(p.StringVal, occurrences)
	h := entropy(probabilities)
	p.H = h
	if h < 4.0 {
		return h,
			&types.HomogeneityError{
				Cardinality:       p.Cardinality,
				LowestProbability: minMapRuneFloat32(probabilities),
			}
	}
	return h, nil
}

// minMapRuneFloat32 gives the min number map[rune]float32 for custom errors
func minMapRuneFloat32(runeMap map[rune]float32) float32 {
	var minFloat float32 = math.MaxFloat32
	for _, candidate := range runeMap {
		if candidate < minFloat {
			minFloat = candidate
		}
	}
	return minFloat
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

// Calculate the entropy using the equation
// h = Σ prob(i) * log₂(1/prob(i))
func entropy(probabilities map[rune]float32) float32 {
	var h float64 = 0.0
	for _, probability := range probabilities {
		h += float64(probability) * float64(math.Log2(float64(1.0/probability)))
	}
	return float32(h)
}
