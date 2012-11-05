package main
import "fmt"

type BloomFilter struct {
	// TODO use a more conservative data-type like 'bit' for bitmap
	bitmap []int // The bloom-filter bitmap
	k int // The number of hash functions
	n int // Number of elements in the filter
	m int // Size of the bloom filter
}

func newBloomFilter(k, m int) *BloomFilter {
	bf := new (BloomFilter)
	bf.bitmap = make([]int, m)
	bf.k, bf.m = k, m
	bf.n = 0
	return bf
}

// TODO Fill this
func (bf *BloomFilter) add() {
}

// TODO Fill this
func (bf *BloomFilter) check(x int64) bool {
	return false
}

func main() {
	fmt.Println("Testing out the Bloom Filter")
	bf := newBloomFilter(10, 20)
	fmt.Println(bf)
}
