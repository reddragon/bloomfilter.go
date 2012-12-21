package bloomfilter

import (
	"hash"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitmap []bool      // The bloom-filter bitmap
	k      int         // Number of hash functions
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
	bf.hashfn.Write(b)
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
	bf.n++
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

func (bf *BloomFilter) falsePositiveRate() float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/
		float64(bf.m))), float64(bf.k))
}

type CountingBloomFilter struct {
	counts []uint8     // The bloom-filter bitmap
	k      int         // Number of hash functions
	n      int         // Number of elements in the filter
	m      int         // Size of the bloom filter
	hashfn hash.Hash64 // The hash function
}

func newCountingBloomFilter(k, m int) *CountingBloomFilter {
	cbf := new(CountingBloomFilter)
	cbf.counts = make([]uint8, m)
	cbf.k, cbf.m = k, m
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

func (cbf *CountingBloomFilter) add(e []byte) {
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

func (cbf *CountingBloomFilter) remove(e []byte) {
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

func (cbf *CountingBloomFilter) check(x []byte) bool {
	h1, h2 := cbf.getHash(x)
	result := true
	for i := 0; i < cbf.k; i++ {
		ind := (h1 + uint32(i)*h2) % uint32(cbf.m)
		result = result && (cbf.counts[ind] > 0)
	}
	return result
}

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

func newScalableBloomFilter(k, m, p, r int, f float64) *ScalableBloomFilter {
	sbf := new(ScalableBloomFilter)
	sbf.k, sbf.n, sbf.m, sbf.p, sbf.q, sbf.r, sbf.f = k, 0, m, p, 1, r, f
	sbf.s = sbf.m
	sbf.bfArr = make([]BloomFilter, 0, p)
	bf := newBloomFilter(sbf.k, sbf.m)
	sbf.bfArr = append(sbf.bfArr, *bf)
	return sbf
}

func (sbf *ScalableBloomFilter) add(e []byte) {
	inuseFilter := sbf.q - 1
	fpr := sbf.bfArr[inuseFilter].falsePositiveRate()
	if fpr <= sbf.f {
		sbf.bfArr[inuseFilter].add(e)
		sbf.n++
	} else {
		if sbf.p == sbf.q {
			return
		}
		sbf.s = sbf.s * sbf.r
		bf := newBloomFilter(sbf.k, sbf.s)
		sbf.bfArr = append(sbf.bfArr, *bf)
		sbf.q++
		inuseFilter = sbf.q - 1
		sbf.bfArr[inuseFilter].add(e)
		sbf.n++
	}
}

func (sbf *ScalableBloomFilter) falsePositiveRate() float64 {
	res := 1.0
	for i := 0; i < sbf.q; i++ {
		res *= (1.0 - sbf.bfArr[i].falsePositiveRate())
	}
	return 1.0 - res
}

func (sbf *ScalableBloomFilter) check(e []byte) bool {
	for i := 0; i < sbf.q; i++ {
		if sbf.bfArr[i].check(e) {
			return true
		}
	}
	return false
}
