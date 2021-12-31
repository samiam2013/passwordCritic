package main

import (
	"fmt"

	"github.com/samiam2013/passwordcritic/types"
)

func RebuildFilters() ([]types.BloomFilter, error) {
	files := map[int]string{
		1_000:     "./cache/10-million-password-list-top-1000.txt",
		10_000:    "./cache/10-million-password-list-top-10000",
		100_000:   "./cache/10-million-password-list-top-100000.txt",
		1_000_000: "./cache/10-million-password-list-top-1000000.txt",
	}
	filters := make([]types.BloomFilter, 0)
	for count, filePath := range files {
		bitsNeeded := int(float32(count) * 12.364167) // only works for 3 hash functions
		newFilter, err := buildFilter(count, bitsNeeded, filePath)
		if err != nil {
			return nil, fmt.Errorf("error building filter for file '%s': %s",
				filePath, err.Error())
		}
		filters = append(filters, *newFilter)
	}

	return filters, nil
}

func buildFilter(limit, bits int, filepath string) (bFilter *types.BloomFilter, err error) {
	err = nil
	bFilter = types.NewBloom(bits)

	return
}

// what's it like typing on this keyboard instead? not grea
// 	i NEED to make money, NOW-ish

// it would be hilariously difficult to _not_ get a job in 30 minutes
//  if I got a kubernetes certificate7
