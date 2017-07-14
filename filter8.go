package bloomfilter

import (
	"sync"
)

type filter8 struct {
	hashCount uint
	counters  []uint8
	mutex     sync.RWMutex
}

func (f *filter8) Add(value string) uint {
	return f.applyToValueCounters(value, func(counter *uint8) {
		if *counter < 255 {
			*counter++
		}
	})
}

func (f *filter8) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.counters = make([]uint8, len(f.counters))
}

func (f *filter8) Count(value string) uint {
	return f.applyToValueCounters(value, nil)
}

func (f *filter8) Remove(value string) uint {
	return f.applyToValueCounters(value, func(counter *uint8) {
		if *counter > 0 {
			*counter--
		}
	})
}

func (f *filter8) Reset(value string) {
	f.applyToValueCounters(value, func(counter *uint8) {
		*counter = 0
	})
}

func (f *filter8) applyToValueCounters(value string, op func(counter *uint8)) uint {
	low, high := hashesFor(value)
	count := uint8(255)

	if op == nil {
		f.mutex.RLock()
		defer f.mutex.RUnlock()
	} else {
		f.mutex.Lock()
		defer f.mutex.Unlock()
	}

	for i := uint(0); i < f.hashCount; i++ {
		index := (low + i*high) % uint(len(f.counters))

		if op != nil {
			op(&f.counters[index])
		}

		if f.counters[index] < count {
			count = f.counters[index]
		}
	}

	return uint(count)
}

func New8Bit(valuesToAdd uint, falsePositiveRate float64) Filter {
	hashCount, filterSize := filterSizeFor(valuesToAdd, falsePositiveRate)

	return &filter8{
		hashCount: hashCount,
		counters:  make([]uint8, filterSize),
	}
}
