// Package intintmap is a fast uint64 key -> uint64 value map.
//
// It is copied nearly verbatim from http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/
package intintmap

import (
	"math"
)

// IntPhi is for scrambling the keys
const IntPhi = uint64(0x9E3779B9)

func phiMix(x uint64) uint64 {
	h := x * IntPhi
	return h ^ (h >> 16)
}

// Map is a map-like data-structure for int64s
type Map struct {
	data       []uint64 // interleaved keys and values
	fillFactor float64
	threshold  int // we will resize a map once it reaches this size
	size       int

	mask  uint64 // mask to calculate the original position
	mask2 uint64
}

func nextPowerOf2(x uint32) uint32 {
	if x == math.MaxUint32 {
		return x
	}

	if x == 0 {
		return 1
	}

	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16

	return x + 1
}

func arraySize(exp int, fill float64) int {
	s := nextPowerOf2(uint32(math.Ceil(float64(exp) / fill)))
	if s < 2 {
		s = 2
	}
	return int(s)
}

// New returns a map initialized with n spaces and uses the stated fillFactor.
// The map will grow as needed.
func New(size int, fillFactor float64) *Map {
	if fillFactor <= 0 || fillFactor >= 1 {
		panic("FillFactor must be in (0, 1)")
	}
	if size <= 0 {
		panic("Size must be positive")
	}

	capacity := arraySize(size, fillFactor)
	return &Map{
		data:       make([]uint64, 2*capacity),
		fillFactor: fillFactor,
		threshold:  int(math.Floor(float64(capacity) * fillFactor)),
		mask:       uint64(capacity - 1),
		mask2:      uint64(2*capacity - 1),
	}
}

// Clear clears all key->value associations in the map but preserves memory
func (m *Map) Clear() {
	m.size = 0
	for i := range m.data {
		m.data[i] = 0
	}
}

// Get returns the value if the key is found.
func (m *Map) Get(key uint64) (uint64, bool) {
	if key == 0 {
		return 0, false
	}

	ptr := (phiMix(key) & m.mask) << 1
	if ptr < 0 || ptr >= uint64(len(m.data)) { // Check to help to compiler to eliminate a bounds check below.
		return 0, false
	}
	k := m.data[ptr]

	if k == 0 { // end of chain already
		return 0, false
	}
	if k == key { // we check FREE prior to this call
		return m.data[ptr+1], true
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]
		if k == 0 {
			return 0, false
		}
		if k == key {
			return m.data[ptr+1], true
		}
	}
}

// Put adds or updates key with value val.
func (m *Map) Put(key uint64, val uint64) {
	if key == 0 {
		panic("zero key are illegal")
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == 0 { // end of chain already
		m.data[ptr] = key
		m.data[ptr+1] = val
		if m.size >= m.threshold {
			m.rehash()
		} else {
			m.size++
		}
		return
	} else if k == key { // overwrite existed value
		m.data[ptr+1] = val
		return
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]

		if k == 0 {
			m.data[ptr] = key
			m.data[ptr+1] = val
			if m.size >= m.threshold {
				m.rehash()
			} else {
				m.size++
			}
			return
		} else if k == key {
			m.data[ptr+1] = val
			return
		}
	}

}

// Del deletes a key and its value.
func (m *Map) Del(key uint64) {
	if key == 0 {
		return
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == key {
		m.shiftKeys(ptr)
		m.size--
		return
	} else if k == 0 { // end of chain already
		return
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]

		if k == key {
			m.shiftKeys(ptr)
			m.size--
			return
		} else if k == 0 {
			return
		}

	}
}

func (m *Map) shiftKeys(pos uint64) uint64 {
	// Shift entries with the same hash.
	var last, slot uint64
	var k uint64
	var data = m.data
	for {
		last = pos
		pos = (last + 2) & m.mask2
		for {
			k = data[pos]
			if k == 0 {
				data[last] = 0
				return last
			}

			slot = (phiMix(k) & m.mask) << 1
			if last <= pos {
				if last >= slot || slot > pos {
					break
				}
			} else {
				if last >= slot && slot > pos {
					break
				}
			}
			pos = (pos + 2) & m.mask2
		}
		data[last] = k
		data[last+1] = data[pos+1]
	}
}

func (m *Map) rehash() {
	newCapacity := len(m.data) * 2
	m.threshold = int(math.Floor(float64(newCapacity/2) * m.fillFactor))
	m.mask = uint64(newCapacity/2 - 1)
	m.mask2 = uint64(newCapacity - 1)

	data := make([]uint64, len(m.data)) // copy of original data
	copy(data, m.data)

	m.data = make([]uint64, newCapacity)
	m.size = 0

	var o uint64
	for i := 0; i < len(data); i += 2 {
		o = data[i]
		if o != 0 {
			m.Put(o, data[i+1])
		}
	}
}

// Size returns size of the map.
func (m *Map) Size() int {
	return m.size
}

// Cap returns the capacity of the map
func (m *Map) Cap() int {
	return m.threshold
}

// Iter call the provided function for each key, value pair.
// The provided function should return true if the iteration should continue
func (m *Map) Iter(fn func(uint64, uint64) bool) {
	data := m.data
	var k uint64

	for i := 0; i < len(data); i += 2 {
		k = data[i]
		if k == 0 {
			continue
		}
		if !fn(k, data[i+1]) {
			return
		}
	}
}
