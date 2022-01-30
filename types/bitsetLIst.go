package types

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// BitSetMap holds lists of the built filters for json storage/loading
type BitSetMap struct {
	// map of the bitsets indexed by # of elements (pws) in the filter
	List map[int]BitSet `json:"list"`
}

// MarshalJSON overrides the interface{} marshalling behavior or BitsetMap
func (bm *BitSetMap) MarshalJSON() ([]byte, error) {
	return []byte(`{"list":{ "1000": {"bitset": [true, false, true]} }}`), nil
}

// BitSet holds a slice of ZeroOneBools for custom Marshall/Unmarshall
type BitSet struct {
	Set []bool `json:"bitset"`
}

// MarshalJSON marshals the boolean types to string 0 or 1
func (bs *BitSet) MarshalJSON() ([]byte, error) {
	// base64 holds 6 bits in a byte, so iterate over each 6 bits getting 1 byte
	bitLen := len(bs.Set)
	if bitLen%6 != 0 {
		bitRemainder := 6 - (bitLen % 6)
		bitLen += bitRemainder
		bs.Set = append(bs.Set, make([]bool, bitRemainder)...)
	}
	noBytes := bitLen / 6
	fmt.Printf("number of bytes to marshall: %d", noBytes)
	byteArr := make([]byte, noBytes)
	for byteIdx := 0; byteIdx < noBytes; byteIdx++ {
		for bit := 0; bit < 6; bit++ {
			maskVal := int(math.Pow(2, float64(bit)))
			// fmt.Printf("maskVal: %d\n", maskVal)
			i := bit + ((byteIdx) * 6)
			if bs.Set[i] {
				// mask the current bit
				byteArr[byteIdx] |= byte(maskVal)
			} else {
				const ones = byte(255)
				byteArr[byteIdx] &= (byte(maskVal) ^ ones)
			}
		}
		byteArr[byteIdx] += byte(48)
	}
	return byteArr, nil
}

// UnmarshalJSON unpacks the string 1 or 0 to a set of sythetic boolean type
func (bs *BitSet) UnmarshalJSON(data []byte) error {
	bs.Set = make([]bool, len(bs.Set)*6)
	for i, byt := range data {
		bitOffset := i * 6
		for j := 0; j < 6; j++ {
			// iterate over the bits needing unpacked and-ing them against
			//	the 2^j generated mask and checking if they are set
			//	by comparing against 0
			mask := byte(2 ^ j)
			bs.Set[bitOffset+j] = (mask & byt) != byte(0)
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
		Set: make([]bool, len(b.Bitset)),
	}
	for i := 0; i < len(b.Bitset); i++ {
		bs.Set[i] = b.Bitset[i]
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
