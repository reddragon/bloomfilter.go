/*
Copyright (c) 2013, Gaurav Menghani
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package bloomfilter

import (
	"hash"
	"hash/fnv"
	"math"
)

// The standard bloom filter, which allows adding of 
// elements, and checking for their existence
type BloomFilter struct {
	bitmap []bool      // The bloom-filter bitmap
	k      int         // Number of hash functions
	n      int         // Number of elements in the filter
	m      int         // Size of the bloom filter
	hashfn hash.Hash64 // The hash function
}

// Returns a new BloomFilter object, if you pass the 
// number of Hash Functions to use and the maximum
// size of the Bloom Filter
func NewBloomFilter(numHashFuncs, bfSize int) *BloomFilter {
	bf := new(BloomFilter)
	bf.bitmap = make([]bool, bfSize)
	bf.k, bf.m = numHashFuncs, bfSize
	bf.n = 0
	bf.hashfn = fnv.New64()
	return bf
}

func (bf *BloomFilter) getHash(b []byte) (uint32, uint32) {
	bf.hashfn.Reset()
	bf.hashfn.Write(b)
	hash64 := bf.hashfn.Sum64()
	h1 := uint32(hash64 & ((1 << 32) - 1))
	h2 := uint32(hash64 >> 32)
	return h1, h2
}

// Adds an element (in byte-array form) to the Bloom Filter
func (bf *BloomFilter) Add(e []byte) {
	h1, h2 := bf.getHash(e)
	for i := 0; i < bf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.m)
		bf.bitmap[ind] = true
	}
	bf.n++
}

// Checks if an element (in byte-array form) exists in the 
// Bloom Filter
func (bf *BloomFilter) Check(x []byte) bool {
	h1, h2 := bf.getHash(x)
	result := true
	for i := 0; i < bf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(bf.m)
		result = result && bf.bitmap[ind]
	}
	return result
}

// Returns the current False Positive Rate of the Bloom Filter
func (bf *BloomFilter) FalsePositiveRate() float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/
		float64(bf.m))), float64(bf.k))
}

// A Bloom Filter which allows deletion of elements. 
// An 8-bit counter is maintained for each slot. This should
// be accounted for while deciding the size of the new filter.
type CountingBloomFilter struct {
	counts []uint8     // The bloom-filter bitmap
	k      int         // Number of hash functions
	n      int         // Number of elements in the filter
	m      int         // Size of the bloom filter
	hashfn hash.Hash64 // The hash function
}

// Creates a new Counting Bloom Filter
func NewCountingBloomFilter(numHashFuncs, cbfSize int) *CountingBloomFilter {
	cbf := new(CountingBloomFilter)
	cbf.counts = make([]uint8, cbfSize)
	cbf.k, cbf.m = numHashFuncs, cbfSize
	cbf.n = 0
	cbf.hashfn = fnv.New64()
	return cbf
}

func (cbf *CountingBloomFilter) getHash(b []byte) (uint32, uint32) {
	cbf.hashfn.Reset()
	cbf.hashfn.Write(b)
	hash64 := cbf.hashfn.Sum64()
	h1 := uint32(hash64 & ((1 << 32) - 1))
	h2 := uint32(hash64 >> 32)
	return h1, h2
}

// Adds an element (in byte-array form) to the Counting Bloom Filter
func (cbf *CountingBloomFilter) Add(e []byte) {
	h1, h2 := cbf.getHash(e)
	for i := 0; i < cbf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(cbf.m)
		// Guarding against an overflow
		if cbf.counts[ind] < 0xFF {
			cbf.counts[ind] += 1
		}
	}
	cbf.n++
}

// Removes an element (in byte-array form) from the Counting Bloom Filter
func (cbf *CountingBloomFilter) Remove(e []byte) {
	h1, h2 := cbf.getHash(e)
	for i := 0; i < cbf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(cbf.m)

		if cbf.counts[ind] > 0 {
			// Guarding against an underflow
			cbf.counts[ind] -= 1
		}
	}
	cbf.n--
}

// Checks if an element (in byte-array form) exists in the 
// Counting Bloom Filter
func (cbf *CountingBloomFilter) Check(x []byte) bool {
	h1, h2 := cbf.getHash(x)
	result := true
	for i := 0; i < cbf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(cbf.m)
		result = result && (cbf.counts[ind] > 0)
	}
	return result
}

// A scalable bloom filter, which allows adding of 
// elements, and checking for their existence
type ScalableBloomFilter struct {
	bfArr []BloomFilter // The list of Bloom Filters
	k     int           // Number of hash functions
	n     int           // Number of elements in the filter
	m     int           // Size of the smallest bloom filter
	p     int           // Maximum number of bloom filters to support. 	
	q     int           // Number of bloom filters present in the list.
	r     int           // Multiplication factor for new bloom filter sizes
	s     int           // Size of the current bloom filter
	f     float64       // Target False Positive rate / bf
}

// Returns a new Scalable BloomFilter object, if you pass in
// valid values for all the required fields.
// firstBFSize is the size of the first Bloom Filter which
// will be created.
// maxBloomFilters is the upper limit on the number of 
// Bloom Filters to create
// growthFactor is the rate at which the Bloom Filter size grows.
// targetFPR is the maximum false positive rate allowed for each
// of the constituent bloom filters, after which a new Bloom
// Filter would be created and used
func NewScalableBloomFilter(numHashFuncs, firstBFSize, maxBloomFilters, growthFactor int, targetFPR float64) *ScalableBloomFilter {
	sbf := new(ScalableBloomFilter)
	sbf.k, sbf.n, sbf.m, sbf.p, sbf.q, sbf.r, sbf.f = numHashFuncs, 0, firstBFSize, maxBloomFilters, 1, growthFactor, targetFPR
	sbf.s = sbf.m
	sbf.bfArr = make([]BloomFilter, 0, maxBloomFilters)
	bf := NewBloomFilter(sbf.k, sbf.m)
	sbf.bfArr = append(sbf.bfArr, *bf)
	return sbf
}

// Adds an element of type byte-array to the Bloom Filter
func (sbf *ScalableBloomFilter) Add(e []byte) {
	inuseFilter := sbf.q - 1
	fpr := sbf.bfArr[inuseFilter].FalsePositiveRate()
	if fpr <= sbf.f {
		sbf.bfArr[inuseFilter].Add(e)
		sbf.n++
	} else {
		if sbf.p == sbf.q {
			return
		}
		sbf.s = sbf.s * sbf.r
		bf := NewBloomFilter(sbf.k, sbf.s)
		sbf.bfArr = append(sbf.bfArr, *bf)
		sbf.q++
		inuseFilter = sbf.q - 1
		sbf.bfArr[inuseFilter].Add(e)
		sbf.n++
	}
}

// Returns the cumulative False Positive Rate of the filter
func (sbf *ScalableBloomFilter) FalsePositiveRate() float64 {
	res := 1.0
	for i := 0; i < sbf.q; i++ {
		res *= (1.0 - sbf.bfArr[i].FalsePositiveRate())
	}
	return 1.0 - res
}

// Checks if an element (in byte-array form) exists
func (sbf *ScalableBloomFilter) Check(e []byte) bool {
	for i := 0; i < sbf.q; i++ {
		if sbf.bfArr[i].Check(e) {
			return true
		}
	}
	return false
}
