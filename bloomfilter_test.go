package bloomfilter

import (
	"testing"

//"fmt"
)

func TestBasic(t *testing.T) {
	bf := newBloomFilter(3, 100)
	d1, d2 := []byte("Hello"), []byte("Jello")
	bf.add(d1)

	if !bf.check(d1) {
		t.Errorf("d1 should be present in the BloomFilter")
	}

	if bf.check(d2) {
		t.Errorf("d2 should be absent from the BloomFilter")
	}
}
