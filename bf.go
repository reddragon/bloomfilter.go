package main

import "fmt"
import "hash"
import "hash/fnv"

// TODO
// 1. Allow the user to specify a hash function?
// 2. Support for concurrent inserts?
// 3. Do performance testing

type BloomFilter struct {
	bitmap []bool      // The bloom-filter bitmap
	k      int         // The number of hash functions
	n      int         // Number of elements in the filter
	m      int         // Size of the bloom filter
	hashfn hash.Hash64 // The hash function
}

func newBloomFilter(k, m int) *BloomFilter {
	bf := new(BloomFilter)
	bf.bitmap = make([]bool, m)
	bf.k, bf.m = k, m
	bf.n = 0
	bf.hashfn = fnv.New64()
	return bf
}

func (bf *BloomFilter) getHash(b []byte) (uint32, uint32) {
	bf.hashfn.Reset()
	bf.hashfn.Sum(b)
	hash64 := bf.hashfn.Sum64()
	h1 := uint32(hash64 & ((1 << 32) - 1))
	h2 := uint32(hash64 >> 32)
	return h1, h2
}

func (bf *BloomFilter) add(e []byte) {
	h1, h2 := bf.getHash(e)
	for i := 0; i < bf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.m)
		bf.bitmap[ind] = true
	}
}

func (bf *BloomFilter) check(x []byte) bool {
	h1, h2 := bf.getHash(x)
	result := true
	for i := 0; i < bf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.m)
		result = result && bf.bitmap[ind]
	}
	return result
}

func main() {
	fmt.Println("Testing out the Bloom Filter")
	bf := newBloomFilter(10, 20)
	data := []byte("Hello")
	bf.add(data)
	fmt.Println(bf.check(data))
}
