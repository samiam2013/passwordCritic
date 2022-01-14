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

// ZeroOneBool masks the JSON marshal behavior of bool to be int 0|1 instead of string true|false
type ZeroOneBool bool

// MarshalJSON marshals the boolean types to string 0 or 1
func (bit *ZeroOneBool) MarshalJSON() ([]byte, error) {
	if *bit {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// UnmarshalJSON unpacks the string 1 or 0 to a boolean type
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

// LoadFromRebuild takes in the map from LoadFromFile or DownloadLists and puts the filters into the BitsetList
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

// WriteToFile puts the BitsetList in a persistent filestor in JSON format
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

// LoadFromFile reads in the BitsetList stored in JSON by BitsetList.WriteToFile()
func (bl *BitsetList) LoadFromFile(pathToFile string) (list map[int]BloomFilter, err error) {
	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileBytes, bl)

	list = bl.getFilters()

	return
}
