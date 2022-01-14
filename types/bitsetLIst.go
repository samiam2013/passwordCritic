package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// BitsetList holds lists of the built filters for json storage/loading
type BitsetList struct {
	// map of the bitsets indexed by # of elements (pws) in the filter
	List map[int][]ZeroOneBool `json:"list"`
}

type ZeroOneBool bool

func (bit *ZeroOneBool) MarshalJSON() ([]byte, error) {
	if *bit {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

func (bit *ZeroOneBool) UnmarshalJSON(data []byte) error {
	asInt, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	} else if asInt == 1 || string(data) == "true" {
		*bit = true
	} else if asInt == 0 || string(data) == "false" {
		*bit = false
	} else {
		return fmt.Errorf("boolean unmarshall error: invalid input %s", data)
	}
	return nil
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
	bools := make([]ZeroOneBool, len(b.Bitset))
	for i := 0; i < len(b.Bitset); i++ {
		bools[i] = ZeroOneBool(b.Bitset[i]) //b.Bitset[i]
	}
	bl.List[elems] = bools
	return nil
}

func (bl *BitsetList) getFilters() (list map[int]BloomFilter) {
	list = make(map[int]BloomFilter)
	for elems, Bits := range bl.List {
		newFilter := *NewBloom(len(Bits))
		for i, bit := range Bits {
			newFilter.Bitset[i] = bool(bit)
		}
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
