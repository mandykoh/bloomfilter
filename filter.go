package bloomfilter

import (
	"hash/fnv"
	"math"
)

type Filter interface {
	Add(value string) uint
	Clear()
	Count(value string) uint
	Remove(value string) uint
	Reset(value string)
}

func filterSizeFor(valuesToAdd uint, falsePositiveRate float64) (hashCount, size uint) {
	logFP := math.Log(falsePositiveRate)
	log2 := math.Log(2)

	hashCount = uint(-(logFP / log2))
	size = uint(-(float64(valuesToAdd) * logFP) / math.Pow(log2, 2))
	return
}

func hashesFor(value string) (low, high uint) {
	fnv := fnv.New64a()
	fnv.Write([]byte(value))
	hash := fnv.Sum64()
	return uint(hash >> 32), uint(hash & 0xffffffff)
}
