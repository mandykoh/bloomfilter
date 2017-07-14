# bloomfilter

[![GoDoc](https://godoc.org/github.com/mandykoh/bloomfilter?status.svg)](https://godoc.org/github.com/mandykoh/bloomfilter)
[![Go Report Card](https://goreportcard.com/badge/github.com/mandykoh/bloomfilter)](https://goreportcard.com/report/github.com/mandykoh/bloomfilter)
[![Build Status](https://travis-ci.org/mandykoh/bloomfilter.svg?branch=master)](https://travis-ci.org/mandykoh/bloomfilter)

A counting Bloom filter implementation for Go.

## Introduction

This package provides thread-safe implementations of counting Bloom filters with 8- and 16-bit counters.


## Example

Use `bloomfilter` like this:

```go
f8 := bloomfilter.New8Bit(valuesToInsert, falsePositiveRate)

f16 := bloomfilter.New16Bit(valuesToInsert, falsePositiveRate)
```

where `valuesToInsert` is the expected number of items to add to the filter, and `falsePositiveRate` is the desired rate of false positives (eg 0.1 for a 10% false positive rate).

The counts that a filter will keep track of are limited by the resolution of the filter; eg an 8-bit filter will allow counts up to 255, at which point they stop increasing (counters do not overflow). A 16-bit filter allows for counts up to 65,535. Higher resolution filters also take up more memory.

Add values to the filter like this:

```go
count = f.Add("some value")
```

where the returned `count` is the upper bound for the number of times that value has been added before. (The value may have been added up to this many times, but definitely not more.)

Explicitly check the count for a value like this:

```go
count = f.Count("some value")
```

Remove a value from the filter like this:

```go
count = f.Remove("some value")
```

where the returned `count` is the count for that value after removal.

Reset the count for a single value like this:

```go
f.Reset("some value")
```

Immediately after this, the count for that value is guaranteed to be zero.

Clear all values from the filter like this:

```go
f.Clear()
```

Immediately after this, the counts for all values will be zero.
