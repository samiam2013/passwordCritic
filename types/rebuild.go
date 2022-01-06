package types

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// CacheFolderPath keeps a single reference to password list directory
const CacheFolderPath = "../cache"

func RebuildFilters() (map[int]BloomFilter, error) {
	files := map[int]string{
		1_000:     CacheFolderPath + "/1000.txt",
		10_000:    CacheFolderPath + "/10000.txt",
		100_000:   CacheFolderPath + "/100000.txt",
		1_000_000: CacheFolderPath + "/1000000.txt",
	}
	filters := make(map[int]BloomFilter)
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
		filters[count] = *newFilter
	}

	return filters, nil
}

func buildFilter(bits int, fh io.Reader) (bFilter *BloomFilter, err error) {
	err = nil
	bFilter = NewBloom(bits)

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		bFilter.Add(scanner.Bytes())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
