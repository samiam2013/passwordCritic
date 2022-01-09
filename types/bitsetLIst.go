package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// BitsetList holds lists of the built filters for json storage/loading
type BitsetList struct {
	// map of the bitsets indexed by # of elements (pws) in the filter
	List map[int][]bool `json:"list"`
}

func (bl *BitsetList) LoadFromRebuild(filters map[int]BloomFilter) error {
	for elems, bFilter := range filters {
		err := bl.addFilter(elems, bFilter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bl *BitsetList) addFilter(elems int, b BloomFilter) error {
	if _, ok := bl.List[elems]; ok {
		return fmt.Errorf("key (# passwords) '%d' already set", elems)
	}
	bl.List[elems] = b.Bitset
	return nil
}

func (bl *BitsetList) getFilters() (list map[int]BloomFilter) {
	list = make(map[int]BloomFilter)
	for elems, Bits := range bl.List {
		newFilter := *NewBloom(len(Bits))
		newFilter.Bitset = Bits
		newFilter.N = uint(elems)
		list[elems] = newFilter
	}
	return
}

func (bl *BitsetList) WriteToFile(pathToFile string) error {
	fh, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer fh.Close()

	jsonBytes, err := json.Marshal(bl)
	if err != nil {
		return err
	}

	_, err = fh.Write(jsonBytes)
	return err
}

func (bl *BitsetList) LoadFromFile(pathToFile string) (list map[int]BloomFilter, err error) {
	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileBytes, bl)

	list = bl.getFilters()

	return
}
