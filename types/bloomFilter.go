package types

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/spaolacci/murmur3"
)

// this code adapted from https://codeburst.io/lets-implement-a-bloom-filter-in-go-b2da8a4b849f

type BFilter interface {
	Add(item []byte) error          // Adds the item into the Set
	Test(item []byte) (bool, error) // Check if items is maybe in the Set
}

// BloomFilter probabilistic bloom filter struct
type BloomFilter struct {
	bitset    []bool // The bloom-filter bitset
	k         uint   // Number of hash values
	n         uint   // Number of elements in the filter
	hashFuncs []hash.Hash64
}

// force Bloomfilter Struct to fit the interface defined by BFilter
var _ BFilter = &BloomFilter{}

// Returns a new BloomFilter struct
func NewBloom(size int) *BloomFilter {
	return &BloomFilter{
		bitset:    make([]bool, size),
		k:         3, // we have 3 hash functions (for now)
		n:         uint(0),
		hashFuncs: []hash.Hash64{murmur3.New64(), fnv.New64(), fnv.New64a()},
	}
}

// Adds the item into the bloom filter set by hashing in over the hash functions
func (bf *BloomFilter) Add(item []byte) error {
	hashes, err := bf.hashValues(item)
	if err != nil {
		return fmt.Errorf("couldn't get hashes for adding item to bloom filter: %s", err.Error())
	}

	m := uint(len(bf.bitset))
	for i := uint(0); i < bf.k; i++ {
		position := uint(hashes[i]) % m
		bf.bitset[uint(position)] = true
	}
	bf.n++
	return nil
}

// Calculates all the hash values by applying in the item over the hash functions
func (bf *BloomFilter) hashValues(item []byte) ([]uint64, error) {
	var result []uint64
	for _, hashFunc := range bf.hashFuncs {
		_, err := hashFunc.Write(item)
		if err != nil {
			return result, fmt.Errorf(
				"trouble getting hash sum from '%+v': %s", hashFunc, err.Error())
		}
		result = append(result, hashFunc.Sum64())
		hashFunc.Reset()
	}
	return result, nil
}

// Test if the item into the bloom filter is set by hashing in over // the hash functions
func (bf *BloomFilter) Test(item []byte) (exists bool, failure error) {
	hashes, err := bf.hashValues(item)
	if err != nil {
		failure = fmt.Errorf("failed hashing to test item in filter: %s", err.Error())
		return
	}

	exists = true
	m := uint(len(bf.bitset))
	for i := uint(0); i < bf.k; i++ {
		position := uint(hashes[i]) % m
		if !bf.bitset[uint(position)] {
			exists = false
			break
		}
	}
	return
}
