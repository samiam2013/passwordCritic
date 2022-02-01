package types

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// CacheFolder keeps a single reference to password list directory
const CacheFolder = "cache" + string(os.PathSeparator)

// DefaultBitsetFilename determines what file the pre-compiled filter exists or gets rebuilt
const DefaultBitsetFilename = "defaultFilter.json"

func getList() map[int]string {
	CacheFolderPath := getCacheFolder()
	return map[int]string{
		1_000:   CacheFolderPath + "1000.txt",
		10_000:  CacheFolderPath + "10000.txt",
		100_000: CacheFolderPath + "100000.txt",
		// 1_000_000: CacheFolderPath + "1000000.txt",
	}
}

func getCacheFolder() string {
	folder := "../" + CacheFolder
	envVar := os.Getenv("CACHE_FOLDER")
	if envVar != "" {
		folder = envVar
	}
	return folder
}

// RebuildFilters looks at the default filter paths and rebuilds the Bloomfilters
func RebuildFilters() (map[int]BloomFilter, error) {
	filters := make(map[int]BloomFilter)
	for count, filepath := range getList() {
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

	bitset := BitSetMap{
		List: make(map[int]BitSet),
	}
	for elems, filter := range filters {
		bitset.addFilter(elems, filter)
	}
	err := bitset.WriteToFile(getCacheFolder() + DefaultBitsetFilename)
	if err != nil {
		return nil, fmt.Errorf("failed writing filters to file in RebuildFilters(): %s", err.Error())
	}

	return filters, nil
}

// LoadFilters loads the BloomFilters from the Default BitsetFile location
func LoadFilters() (filters map[int]BloomFilter, err error) {
	bitset := BitSetMap{
		List: make(map[int]BitSet),
	}
	filters, err = bitset.LoadFromFile(getCacheFolder() + DefaultBitsetFilename)
	if err != nil {
		return
	}

	return
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
