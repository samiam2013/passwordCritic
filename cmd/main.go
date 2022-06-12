package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/samiam2013/passwordcritic/critic"
	"github.com/samiam2013/passwordcritic/types"
)

// maybe add (-v|--verbose) or (-q|--quiet) flag?

const charProbFloor = 0.1
const minCandidateCardinality = 6 // e.g. 8 char pass with 2 repititions

func main() {
	const minRarity = 1_000 // TODO make this a parameter

	rbFlag := flag.Bool("r", false, "http get & rebuild filters first")

	pwCandPtr := flag.String("p", "", "password to check")
	flag.Parse()

	// check for (-r|--rebuild-filter) and rebuild the filter if needed
	var filters map[int]types.BloomFilter
	var err error
	if *rbFlag {
		filters, err = types.RebuildFilters()
		if err != nil {
			log.Fatalf("error rebuilding filters on flag -r: %s", err.Error())
		}
	} else {
		// else try reading in filter from serialized data (json?)
		log.Print("loading filters..")
		filters, err = types.LoadFilters()
		if err != nil {
			log.Fatalf("error loading filters: %s", err.Error())
		}
	}

	// check the command line for -p="string" -p=string or -p string
	if len(*pwCandPtr) == 0 {
		*pwCandPtr = getStdIn(os.Stdin)
	}

	entropyCandidate, err := checkEntropy(pwCandPtr)
	if err != nil {
		log.Println(err)
	}

	for elemsLen, bFilter := range filters {
		//log.Printf("checking filter with %d elements....", elemsLen)
		exists, err := bFilter.Test([]byte(*pwCandPtr))
		if err != nil {
			log.Fatalf("error checking candidate against %d passwords list: %s",
				elemsLen, err.Error())
		}
		if exists && (minRarity >= elemsLen) {
			log.Fatalf("password too common, found in list of %d passwords, "+
				"minimum rarity set to %d", elemsLen, minRarity)
		} else if exists {
			log.Printf("password common, found in list with %d elements, "+
				"but not more common than %d, minimum set rarity.", elemsLen, minRarity)
			break
		}
	}

	// give the user back information about the password
	fmt.Printf("%+v\n", entropyCandidate)

}

func getStdIn(stdin io.Reader) string {
	output := []rune{}
	reader := bufio.NewReader(stdin)
	for {
		input, _, err := reader.ReadRune()
		if (err != nil && err == io.EOF) || string(input) == "\n" {
			break
		}
		output = append(output, input)
	}
	return string(output)
}

func checkEntropy(pwCandPtr *string) (candidate critic.PassCandidate, err error) {
	// load the password and check the Entropy
	candidate = critic.PassCandidate{}
	candidate.StringVal = *pwCandPtr
	h, err := candidate.Entropy()
	// fmt.Println("Entropy of the password candidate: ", h)
	fmt.Printf("Entropy of the password candidate: %f\n", h)
	if err != nil {
		if _, ok := err.(*types.HomogeneityError); !ok {
			err = fmt.Errorf("non 'homogeneity' type error encounter checking entropy"+
				" of candidate: %s", err.Error())
			return
		}
		hmgError := err.(*types.HomogeneityError)
		if hmgError.LowestProbability < charProbFloor ||
			hmgError.Cardinality < minCandidateCardinality {
			// give an error msg about the least frequent character being too common
			err = fmt.Errorf("high repetition of characters: minimum %f (percentage 0 to 1)",
				hmgError.LowestProbability)
			return
		}
		// give a default case error msg
		err = fmt.Errorf("low entropy for password: mix of low variety and length")
		return
	}
	return
}
