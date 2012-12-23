package bloomfilter

import (
	"encoding/binary"
	"testing"
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

func TestCountingBFBasic(t *testing.T) {
	cbf := newCountingBloomFilter(3, 100)
	d1 := []byte("Hello")
	cbf.add(d1)

	if !cbf.check(d1) {
		t.Errorf("d1 should be present in the BloomFilter")
	}

	cbf.remove(d1)

	if cbf.check(d1) {
		t.Errorf("d1 should be absent from the BloomFilter after deletion")
	}
}

func TestScalableBFBasic(t *testing.T) {
	sbf := newScalableBloomFilter(3, 20, 4, 10, 0.01)

	for i := 1; i < 1000; i++ {
		buf := make([]byte, 8)
		binary.PutVarint(buf, int64(i))
		sbf.add(buf)
		if !sbf.check(buf) {
			t.Errorf("%d should be present in the BloomFilter", i)
			return
		}
	}

	for i := 1; i < 1000; i++ {
		buf := make([]byte, 8)
		binary.PutVarint(buf, int64(i))
		if !sbf.check(buf) {
			t.Errorf("%d should be present in the BloomFilter", i)
			return
		}
	}

	count := 0

	for i := 1000; i < 4000; i++ {
		buf := make([]byte, 8)
		binary.PutVarint(buf, int64(i))
		if sbf.check(buf) {
			count++
		}
	}

	if sbf.falsePositiveRate() > 0.04 {
		t.Errorf("False Positive Rate for this test should be < 0.04")
		return
	}

	sensitivity := 0.01 // TODO Make this configurable
	expectedFalsePositives :=
		(int)((4000 - 1000) * (sbf.falsePositiveRate() + sensitivity))
	if count > expectedFalsePositives {
		t.Errorf("Actual false positives %d is greater than max expected false positives %d",
			count,
			expectedFalsePositives)
		return
	}
}
