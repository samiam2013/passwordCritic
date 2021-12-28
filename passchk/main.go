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

const minLeastCommonCharProb = 0.1
const minCandidateCardinality = 8

func main() {
	// check for (-r|--rebuild-filter) and rebuild the filter if needed
	// else try reading in filter from serialized data (json?)
	// TODO finish!
	//rebuildPtr := flag.Bool("rebuild", false, "set to rebuild password bloom filter")

	// check the command line for (-p|--password-candidate)
	pwCandPtr := flag.String("pw", "", "password to check")
	flag.Parse()

	if len(*pwCandPtr) == 0 {
		*pwCandPtr = getStdIn()
	}

	entropyCandidate, err := checkEntropy(pwCandPtr)
	if err != nil {
		log.Println(err)
	}
	// check if candidate occurs in the 10,000 most common filter
	// if it is error out saying it's far too common

	// ?check if it's in the 10 million most common filter
	// if it is error out saying it's a common password, but not the most common

	// give the user back information about the password
	fmt.Printf("%+v\n", entropyCandidate)

}

func getStdIn() string {
	output := []rune{}
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
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
	fmt.Println("Entropy of the password candidate: ", h)
	if err != nil {
		if _, ok := err.(*types.HomogeneityError); !ok {
			err = fmt.Errorf("non-homogeneity error encounter checking entropy"+
				" of candidate: %s", err.Error())
			return
		}
		hmgError := err.(*types.HomogeneityError)
		if hmgError.LowestProbability > minLeastCommonCharProb {
			// give an error msg about the least frequent character being too common
			err = fmt.Errorf("high repetition of characters: minimum %f (percentage 0 to 1)\n",
				hmgError.LowestProbability)
			return
		} else if hmgError.Cardinality < minCandidateCardinality {
			// give an error msg about the repetition of characters
			err = fmt.Errorf("low variety of characters: %d (expect > %d)\n",
				hmgError.Cardinality, minCandidateCardinality)
			return
		} else {
			// give a default case error msg
			err = fmt.Errorf("low entropy for password: mix of low variety and length\n")
			return
		}
	}
	return
}
