package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// BitSetMap holds lists of the built filters for json storage/loading
type BitSetMap struct {
	// map of the bitsets indexed by # of elements (pws) in the filter
	List map[int]BitSet `json:"list"`
}

// ZeroOneBool masks the JSON marshal behavior of bool to be int 0|1 instead of string true|false
type ZeroOneBool bool

// BitSet holds a slice of ZeroOneBools for custom Marshall/Unmarshall
type BitSet struct {
	Set []ZeroOneBool
}

// MarshalJSON marshals the boolean types to string 0 or 1
func (bs BitSet) MarshalJSON() ([]byte, error) {
	// base64 holds 6 bits in a byte, so iterate over each 6 bits getting 1 byte
	bitLen := len(bs.Set)
	bitLen += 6 - (bitLen % 6)
	noBytes := bitLen / 6
	byteArr := make([]byte, noBytes)
	for byteIdx := 0; byteIdx < noBytes-1; byteIdx++ {
		for i := byteIdx * 6; i < (byteIdx*6)+6; i++ {
			if bs.Set[i] {
				// mask the current bit
				byteArr[byteIdx] |= byte(2 ^ i)
			} else {
				byteArr[byteIdx] &= (byte(2^i) ^ byte(255))
			}
		}
	}
	return []byte(string(byteArr)), nil
}

// UnmarshalJSON unpacks the string 1 or 0 to a set of sythetic boolean type
func (bs *BitSet) UnmarshalJSON(data []byte) error {
	bs.Set = make([]ZeroOneBool, len(bs.Set)*6)
	for i, byt := range data {
		bitOffset := i * 6
		for j := 0; j < 6; j++ {
			// iterate over the bits needing unpacked and-ing them against
			//	the 2^j generated mask and checking if they are set
			//	by comparing against 0
			mask := byte(2 ^ j)
			bs.Set[bitOffset+j] = ZeroOneBool((mask & byt) != byte(0))
		}
	}
	return nil
}

// LoadFromRebuild takes in the map from LoadFromFile or DownloadLists and puts the filters into the BitSetMap
func (bl *BitSetMap) LoadFromRebuild(filters map[int]BloomFilter) error {
	for elems, bFilter := range filters {
		err := bl.addFilter(elems, bFilter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bl *BitSetMap) addFilter(elems int, b BloomFilter) error {
	if _, ok := bl.List[elems]; ok {
		return fmt.Errorf("key (# passwords) '%d' already set", elems)
	}
	bs := BitSet{
		Set: make([]ZeroOneBool, len(b.Bitset)),
	}
	for i := 0; i < len(b.Bitset); i++ {
		bs.Set[i] = ZeroOneBool(b.Bitset[i]) //b.Bitset[i]
	}
	bl.List[elems] = bs
	return nil
}

func (bl *BitSetMap) getFilters() (list map[int]BloomFilter) {
	list = make(map[int]BloomFilter)
	for elems, Bits := range bl.List {
		newFilter := *NewBloom(len(Bits.Set))
		for i, bit := range Bits.Set {
			newFilter.Bitset[i] = bool(bit)
		}
		newFilter.N = uint(elems)
		list[elems] = newFilter
	}
	return
}

// WriteToFile puts the BitSetMap in a persistent filestor in JSON format
func (bl *BitSetMap) WriteToFile(pathToFile string) error {
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

// LoadFromFile reads in the BitSetMap stored in JSON by BitSetMap.WriteToFile()
func (bl *BitSetMap) LoadFromFile(pathToFile string) (list map[int]BloomFilter, err error) {
	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileBytes, bl)

	list = bl.getFilters()

	return
}
