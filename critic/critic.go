package critic

import (
	"fmt"
	"log"
	"math"

	"github.com/samiam2013/passwordcritic/types"
)

// PassCandidate holds relevant information computed about the password
type PassCandidate struct {
	StringVal   string
	Cardinality int
	H           float64
	MinLength   int
	MinEntropy  float64
	MinRarity   types.Rarity
}

// MinLengthGlobal defines the lowest value for password lenght
const MinLengthGlobal = 8

// MinEntropy defines the lowest value for throwing a Homogeneity Error
const MinEntropyDefault = 3.0

// MinLength defines the shortest password allowed
const MinLengthDefault = 10

var MinRarityDefault types.Rarity = 10_000

const minCandidateCardinality = 6 // e.g. 8 char pass with 2 repititions

// NewPassCandidate creates a new PassCandidate taking optional (nil for default) minEntropy and minLength
func NewPassCandidate(candidate string, minEntropy *float64, minLength *int, minRarity *types.Rarity) (*PassCandidate, error) {
	var minE float64
	if minEntropy != nil {
		minE = *minEntropy
	} else {
		minE = MinEntropyDefault
	}
	var minL int
	if minLength != nil {
		minL = *minLength
		if minL < MinLengthGlobal {
			return nil, fmt.Errorf("minLength must be at least %d", MinLengthGlobal)
		}
	} else {
		minL = MinLengthDefault
	}
	var minR types.Rarity
	if minRarity != nil {
		minR = *minRarity
	} else {
		minR = MinRarityDefault
	}
	cand := &PassCandidate{
		StringVal:   candidate,
		Cardinality: 0,
		H:           0.0,
		MinLength:   minL,
		MinEntropy:  minE,
		MinRarity:   minR,
	}
	return cand, nil

}

func (p *PassCandidate) CheckALL() (err error) {
	if _, err := p.CheckLength(); err != nil {
		return err
	}
	h, err := p.CheckEntropy()
	if err != nil {
		if _, ok := err.(*types.HomogeneityError); !ok {
			err = fmt.Errorf("non 'homogeneity' type error encounter checking entropy"+
				" of candidate: %s", err.Error())
			return
		}
		hmgError := err.(*types.HomogeneityError)
		if hmgError.Cardinality < minCandidateCardinality {
			err = fmt.Errorf("high repetition of characters: minimum %f (percentage 0 to 1)",
				hmgError.LowestProbability)
			return
		}
		err = fmt.Errorf("low entropy for password: mix of low variety and length")
		return
	}
	p.H = h
	if _, err := p.IsInFilters(); err != nil {
		return err
	}
	return nil
}

// CheckLength checks the length of the password
func (p *PassCandidate) CheckLength() (int, error) {
	length := len(p.StringVal)
	if length < p.MinLength {
		err := fmt.Errorf("password too short (%d char), minimum %d characters", length, p.MinLength)
		p.H = -1.0
		return length, err
	}
	return length, nil
}

// Entropy returns a float calculated using variety and frequency of characters
func (p *PassCandidate) CheckEntropy() (float64, error) {
	occurrences := charOccurCount(p.StringVal)
	p.Cardinality = len(occurrences)
	probabilities := charProbabilites(p.StringVal, occurrences)
	h := entropy(probabilities)
	p.H = h
	if h < p.MinEntropy {
		return h,
			&types.HomogeneityError{
				Cardinality:       p.Cardinality,
				LowestProbability: minMap(probabilities),
			}
	}
	return h, nil
}

// CheckFrequency checks the frequency of password in the bloom filter
func (p *PassCandidate) IsInFilters() (bool, error) {
	filters, err := types.LoadFilters()
	if err != nil {
		log.Fatalf("error loading filters: %s", err.Error())
		return false, fmt.Errorf("error loading filters: %s", err.Error())
	}
	for elemsLen, bFilter := range filters {
		//log.Printf("checking filter with %d elements....", elemsLen)
		exists, err := bFilter.Test([]byte(p.StringVal))
		if err != nil {
			log.Fatalf("error checking candidate against %d passwords list: %s", elemsLen, err.Error())
		}
		if exists && (p.MinRarity >= elemsLen) {
			log.Fatalf("password too common, found in list of %d passwords, minimum rarity set to %d",
				elemsLen, p.MinRarity)
			return true, fmt.Errorf("password too common, found in list of %d passwords, ", elemsLen)
		} else if exists {
			log.Printf("password common, found in list with %d elements, but not more common than %d,"+
				" minimum set rarity.", elemsLen, p.MinRarity)
		}
	}
	return false, nil
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
