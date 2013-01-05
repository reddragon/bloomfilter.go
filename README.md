bloomfilter.go
==============

A Bloom Filter Implementation in Go

This implementation includes:
- Standard Bloom Filter
- Counting Bloom Filter
  Standard Bloom Filter, except that it supports deletions, by keeping 8-bit counters in each slot. These counters
  are likely to overflow if we add too many elements.
- Scalable Bloom Filter 
  This is a Bloom Filter, which can scale. Standard Bloom Filters will get too full at some time. If you keep adding
  elements any further, the False Positive Rate will go up. Scalable Bloom Filters get rid of this problem, by creating
  a new bloom filter to insert new elements. The newer bloom filters grow exponentially in size, and hence it will take
  exponentially more elements to become full. This is similar to how vectors grow.


Usage
-----

- You can create a new standard bloom filter like this:
`bf := NewBloomFilter(numHashFuncs, bloomFilterSize int)`
  Where,
  1. `numHashFuncs` is the number of hash functions to use in the filter.
  2. `bloomFilterSize` is the size of the bloom filter.

- Adding a new element is as simple as:
`bf.Add(elementByteSlice)`
Where, `elementByteSlice` is a representation of the element in the byte slice format.

- Similarily, you can check if the element is present in the Bloom Filter by doing this:
`bf.Check(elementByteSlice)`

- The Check method is absolutely correct when it returns `false`. However, when it returns `true`, it might 
be wrong with a probability `bf.FalsePositiveRate()`

- The scalable bloom filter is exactly the same, except it is created as follows:
`sbf := NewScalableBloomFilter(numHashFuncs, firstBFSize, maxBloomFilters, growthFactor int, targetFPR float64)`

  Where,
  1. `firstBFSize` is the size of the first Bloom Filter which will be created.
  2. `maxBloomFilters` is the upper limit on the number of Bloom Filters to create
  4. `growthFactor` is the rate at which the Bloom Filter size grows.
  5. `targetFPR` is the maximum false positive rate allowed for each of the constituent bloom filters, after which a new Bloom
  Filter would be created and used

- The counting bloom filter is similar, except it is created as follows:
`cbf := NewCountingBloomFilter(numHashFuncs, cbfSize int)`

- The counting bloom filter also supports deletion:
`cbf.Remove(elementByteSlice)`

Tests
-----
I have written a couple of simple tests, which you can see in `bloomfilter_test.go`, and run the tests by doing `go test`

References
----------
- [Scalable Bloom Filters - Paulo Sérgio Almeida, Carlos Baquero, Nuno Preguiça, David Hutchison](http://www.sciencedirect.com/science/article/pii/S0020019006003127)
- https://github.com/bitly/dablooms
- https://github.com/willf/bloom
- [Less Hashing, Same Performance: Building a Better Bloom Filter - Adam Kirsch and Michael Mitzenmacher](http://www.eecs.harvard.edu/~kirsch/pubs/bbbf/esa06.pdf)
- [Analysis for the false positive rate in Scalable Bloom Filters](http://blog.gaurav.im/?p=278)

Credits
-------
Thanks to [Aditya](https://github.com/truncs) for suggesting the idea and pointers.

Author
------
Gaurav Menghani
