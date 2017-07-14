package bloomfilter

import "sync"

type filter16 struct {
	hashCount uint
	counters  []uint16
	mutex     sync.RWMutex
}

func (f *filter16) Add(value string) uint {
	return f.applyToValueCounters(value, func(counter *uint16) {
		if *counter < 65535 {
			*counter++
		}
	})
}

func (f *filter16) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.counters = make([]uint16, len(f.counters))
}

func (f *filter16) Count(value string) uint {
	return f.applyToValueCounters(value, nil)
}

func (f *filter16) Remove(value string) uint {
	return f.applyToValueCounters(value, func(counter *uint16) {
		if *counter > 0 {
			*counter--
		}
	})
}

func (f *filter16) Reset(value string) {
	f.applyToValueCounters(value, func(counter *uint16) {
		*counter = 0
	})
}

func (f *filter16) applyToValueCounters(value string, op func(counter *uint16)) uint {
	low, high := hashesFor(value)
	count := uint16(65535)

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

func New16Bit(valuesToAdd uint, falsePositiveRate float64) Filter {
	hashCount, filterSize := filterSizeFor(valuesToAdd, falsePositiveRate)

	return &filter16{
		hashCount: hashCount,
		counters:  make([]uint16, filterSize),
	}
}
