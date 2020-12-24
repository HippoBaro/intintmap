> This repo is a fork of `brentp/intintmap`. The main change is that it uses `uint64` instead of `int64`.

Fast uint64 -> uint64 hash in golang.

[![GoDoc](https://godoc.org/github.com/brentp/intintmap?status.svg)](https://godoc.org/github.com/brentp/intintmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/brentp/intintmap)](https://goreportcard.com/report/github.com/brentp/intintmap)

# intintmap

    import "github.com/HippoBaro/intintmap"

Package intintmap is a fast uint64 key -> uint64 value map.

It is copied nearly verbatim from
http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/ .

It interleaves keys and values in the same underlying array to improve locality.

It is 2-5X faster than the builtin map and uses less memory:

```
BenchmarkFillSequential
BenchmarkFillSequential/Std-8           	      18	  65232116 ns/op	 4880488 B/op	    2136 allocs/op
BenchmarkFillSequential/IntInt-8        	      42	  26693671 ns/op	 2393628 B/op	       0 allocs/op
BenchmarkFillSequentialPreAllocated
BenchmarkFillSequentialPreAllocated/Std-8         18	  63034107 ns/op	 2236945 B/op	       1 allocs/op
BenchmarkFillSequentialPreAllocated/IntInt-8      44	  27033307 ns/op	  762602 B/op	       0 allocs/op
BenchmarkFillRandom
BenchmarkFillRandom/Std-8                          6	 202616442 ns/op	61217498 B/op	   36658 allocs/op
BenchmarkFillRandom/IntInt-8                      12	 204325495 ns/op	134206828 B/op	       2 allocs/op
BenchmarkFillRandomPreAllocated
BenchmarkFillRandomPreAllocated/Std-8              6	 219628375 ns/op	53283708 B/op	   30249 allocs/op
BenchmarkFillRandomPreAllocated/IntInt-8          16	 119684014 ns/op	96469003 B/op	       0 allocs/op
BenchmarkLookupSequential
BenchmarkLookupSequential/Std-8                   15	  68508601 ns/op	 5854772 B/op	    2543 allocs/op
BenchmarkLookupSequential/IntInt-8                39	  27021290 ns/op	 2577751 B/op	       0 allocs/op
BenchmarkLookupRandom
BenchmarkLookupRandom/Std-8                        7	 145837643 ns/op	12549013 B/op	    5490 allocs/op
BenchmarkLookupRandom/IntInt-8                    10	 102296550 ns/op	10053249 B/op	       2 allocs/op
BenchmarkLookupNoHit
BenchmarkLookupNoHit/Std-8                         9	 123738290 ns/op	 9761665 B/op	    4265 allocs/op
BenchmarkLookupNoHit/IntInt-8                      9	 118541631 ns/op	11170256 B/op	       2 allocs/op
```

## Usage

```go
m := intintmap.New(8096, 0.6)
m.Put(1234, 222)
m.Put(uint64(123), uint64(33))

v, ok := m.Get(uint64(222))
v, ok := m.Get(uint64(333))

m.Del(uint64(222))
m.Del(uint64(333))

fmt.Println(m.Size())

m.Iter(func (k, v uint64) bool {
fmt.Printf("key: %d, value: %d\n", k, v)
return true
})
```

#### type Map

```go
type Map struct {
}
```

Map is a map-like data-structure for int64s

#### func  New

```go
func New(size int, fillFactor float64) *Map
```

New returns a map initialized with n spaces and uses the stated fillFactor. The map will grow as needed.

#### func (*Map) Get

```go
func (m *Map) Get(key uint64) (uint64, bool)
```

Get returns the value if the key is found.

#### func (*Map) Put

```go
func (m *Map) Put(key uint64, val uint64)
```

Put adds or updates key with value val.

#### func (*Map) Del

```go
func (m *Map) Del(key uint64)
```

Del deletes a key and its value.

#### func (*Map) Iter

```go
func (m *Map) Iter(fn func (uint64, uint64) bool)
```

Iter calls the provided function for each key-value pairs

#### func (*Map) Size

```go
func (m *Map) Size() int
```

Size returns size of the map.
