package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/samiam2013/passwordcritic/types"
)

func RebuildFilters() ([]types.BloomFilter, error) {
	files := map[int]string{
		1_000:     "../cache/10-million-password-list-top-1000.txt",
		10_000:    "../cache/10-million-password-list-top-10000.txt",
		100_000:   "../cache/10-million-password-list-top-100000.txt",
		1_000_000: "../cache/10-million-password-list-top-1000000.txt",
	}
	filters := make([]types.BloomFilter, 0)
	for count, filepath := range files {
		bitsNeeded := int(float32(count) * 12.364167) // only works for 3 hash functions

		fh, err := os.Open(filepath)
		if err != nil {
			return filters, err
		}
		defer fh.Close()

		newFilter, err := buildFilter(bitsNeeded, fh)
		if err != nil {
			return nil, fmt.Errorf("error building filter for file '%s': %s",
				fh.Name(), err.Error())
		}
		filters = append(filters, *newFilter)
	}

	return filters, nil
}

func buildFilter(bits int, fh io.Reader) (bFilter *types.BloomFilter, err error) {
	err = nil
	bFilter = types.NewBloom(bits)

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		bFilter.Add(scanner.Bytes())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
