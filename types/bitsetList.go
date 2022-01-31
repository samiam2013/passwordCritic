package types

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
)

// BitSetMap holds lists of the built filters for json storage/loading
type BitSetMap struct {
	// map of the bitsets indexed by # of elements (pws) in the filter
	List map[int]BitSet `json:"list"`
}

// The struct above creates dependency for example List needs to be index 0

// MarshalJSON overrides the interface{} marshalling behavior or BitsetMap
func (bl *BitSetMap) MarshalJSON() ([]byte, error) {
	list := make(map[int]string, len(bl.List))
	for nElem, bitSet := range bl.List {
		bytes, err := bitSet.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var unmrsld map[string]string
		if err = json.Unmarshal(bytes, &unmrsld); err != nil {
			return nil, fmt.Errorf("failed unmarshalling bitset into interface: %s", err.Error())
		}
		list[nElem] = string(unmrsld["bitset"])
	}

	// get the tag name for marshall
	val := reflect.ValueOf(*bl)
	tag, ok := val.Type().Field(0).Tag.Lookup("json")
	if !ok {
		return []byte{}, fmt.Errorf("error getting json tag for first field in BitSetMap")
	}

	toMarshal := map[string]interface{}{
		tag: map[string]map[int]string{
			"bitset": list,
		},
	}
	bytes, err := json.Marshal(toMarshal)
	if err != nil {
		return []byte{}, fmt.Errorf("failed marshalling final for BitSetList.MarshallJSON():%s", err.Error())
	}
	return bytes, nil
}

// UnmarshalJSON defines custom interaction with BitSet.UnmarshalJSON
func (bl *BitSetMap) UnmarshalJSON(data []byte) error {
	return nil
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

		// 35 - 90 , 97-122
		byteArr[byteIdx] += byte(42)
		if byteArr[byteIdx] >= byte(60) {
			byteArr[byteIdx] += byte(3)
		}
		if byteArr[byteIdx] >= byte(90) {
			byteArr[byteIdx] += byte(7)
		}
	}
	toMarshal := map[string]string{
		"bitset": string(byteArr),
	}
	bytes, err := json.Marshal(toMarshal)
	if err != nil {
		return nil, fmt.Errorf("faild marshalling value for bitset: %s", err.Error())
	}
	return bytes, nil

}

// UnmarshalJSON unpacks the string 1 or 0 to a set of sythetic boolean type
func (bs *BitSet) UnmarshalJSON(data []byte) error {
	bs.Set = []bool{}
	for _, byt := range data {
		bools := make([]bool, 6)
		// work out the value from the character
		// [42-60] +3 [63-90] +7 [97-?)
		if byt >= byte(63) {
			byt -= byte(3)
		}
		if byt >= byte(90) {
			byt -= byte(7)
		}
		byt -= byte(42)
		for j := 0; j < 6; j++ {
			powMask := int(math.Pow(2.0, float64(j)))
			masked := byt & byte(powMask)
			bools[j] = masked != byte(0)
		}
		bs.Set = append(bs.Set, bools...)
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
