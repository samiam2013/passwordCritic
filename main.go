package main

import (
	"fmt"
	"math"
)

func main() {
	inputText := "rthedfff" //"correcthorsebatterystaple"

	occurences := charOccurences(inputText)
	fmt.Println("occurrences", occurences)

	probabilities := charProbabilites(inputText, occurences)
	fmt.Println("probabilities", probabilities)

	entropy := entropy(probabilities)
	fmt.Println("entropy", entropy)
}

// Calculates the # occurrences of each character
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

// Calculate the entropy using the equation
// H = Σp(i)log_2(1/p(i))
func entropy(probabilities map[rune]float32) float32 {
	var h float64 = 0.0
	for _, probability := range probabilities {
		h += float64(probability) * float64(math.Log2(float64(1.0/probability)))
	}
	return float32(h)
}
