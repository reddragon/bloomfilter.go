package main
import "fmt"

type BloomFilter struct {
	bitmap []int // The bloom-filter bitmap
	k int // The number of hash functions
	n int // Number of elements in the filter
	m int // Size of the bloom filter
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
}
