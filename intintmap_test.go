package intintmap

import (
	"math/rand"
	"runtime"
	"testing"
)

func TestMapSimple(t *testing.T) {
	m := New(10, 0.99)
	var i uint64
	var v uint64
	var ok bool

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 1; i < 20000; i += 2 {
		m.Put(i, i)
	}
	for i = 1; i < 20000; i += 2 {
		if v, ok = m.Get(i); !ok || v != i {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	if m.Size() != int(20000/2) {
		t.Errorf("size (%d) is not right, should be %d", m.Size(), int(20000/2))
	}

	// --------------------------------------------------------------------
	// Keys()

	m0 := make(map[uint64]uint64, 1000)
	for i = 1; i < 20000; i += 2 {
		m0[i] = i
	}
	n := len(m0)

	m.Iter(func(k uint64, _ uint64) bool {
		m0[k] = -k
		return true
	})

	if n != len(m0) {
		t.Errorf("get unexpected more keys")
	}

	for k, v := range m0 {
		if k != -v {
			t.Errorf("didn't get expected changed value")
		}
	}

	// --------------------------------------------------------------------
	// Items()

	m0 = make(map[uint64]uint64, 1000)
	for i = 1; i < 20000; i += 2 {
		m0[i] = i
	}
	n = len(m0)

	m.Iter(func(k uint64, v uint64) bool {
		m0[k] = -v
		if v != v {
			t.Errorf("didn't get expected key-value pair")
		}
		return true
	})

	if n != len(m0) {
		t.Errorf("get unexpected more keys")
	}

	for k, v := range m0 {
		if k != -v {
			t.Errorf("didn't get expected changed value")
		}
	}

	// --------------------------------------------------------------------
	// Del()

	for i = 1; i < 20000; i += 2 {
		m.Del(i)
	}
	for i = 1; i < 20000; i += 2 {
		if _, ok = m.Get(i); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 1; i < 20000; i += 2 {
		m.Put(i, i*2)
	}
	for i = 1; i < 20000; i += 2 {
		if v, ok = m.Get(i); !ok || v != i*2 {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

}

func TestMap(t *testing.T) {
	m := New(10, 0.6)
	var ok bool
	var v uint64

	step := uint64(61)

	var i uint64
	m.Put(1, 12345)
	for i = 2; i < 100000000; i += step {
		m.Put(i, i+7)
		m.Put(-i, i-7)

		if v, ok = m.Get(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Get(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}
	}
	for i = 2; i < 100000000; i += step {
		if v, ok = m.Get(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Get(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}

		for j := i + 1; j < i+step; j++ {
			if v, ok = m.Get(j); ok {
				t.Errorf("expected 'not found' flag for %d, found %d", j, v)
			}
		}
	}

	if v, ok = m.Get(1); !ok || v != 12345 {
		t.Errorf("expected 12345 for key 0")
	}
}

func TestReset(t *testing.T) {
	m := New(10, 0.6)
	test := func() {
		var ok bool
		var v uint64

		step := uint64(61)

		var i uint64
		m.Put(1, 12345)
		for i = 2; i < 100000000; i += step {
			m.Put(i, i+7)
			m.Put(-i, i-7)

			if v, ok = m.Get(i); !ok || v != i+7 {
				t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
			}
			if v, ok = m.Get(-i); !ok || v != i-7 {
				t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
			}
		}
		for i = 2; i < 100000000; i += step {
			if v, ok = m.Get(i); !ok || v != i+7 {
				t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
			}
			if v, ok = m.Get(-i); !ok || v != i-7 {
				t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
			}

			for j := i + 1; j < i+step; j++ {
				if v, ok = m.Get(j); ok {
					t.Errorf("expected 'not found' flag for %d, found %d", j, v)
				}
			}
		}

		if v, ok = m.Get(1); !ok || v != 12345 {
			t.Errorf("expected 12345 for key 0")
		}
	}

	test()
	m.Clear()
	test()
}

func BenchmarkFillSequential(b *testing.B) {
	b.ReportAllocs()
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		m := make(map[uint64]uint64, 2048)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[j] = 0
			}
		}
	})

	b.Run("IntInt", func(b *testing.B) {
		b.ReportAllocs()
		m := New(2048, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(j, 0)
			}
		}
	})
}

func BenchmarkFillSequentialPreAllocated(b *testing.B) {
	b.ReportAllocs()
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		m := make(map[uint64]uint64, 1_000_000)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[j] = 0
			}
		}
	})

	b.Run("IntInt", func(b *testing.B) {
		b.ReportAllocs()
		m := New(1_000_000, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(j, 0)
			}
		}
	})
}

func BenchmarkFillRandom(b *testing.B) {
	rand.Seed(0)
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		m := make(map[uint64]uint64, 2048)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[rand.Uint64()+1] = 0
			}
		}
	})

	b.Run("IntInt", func(b *testing.B) {
		b.ReportAllocs()
		m := New(2048, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(rand.Uint64()+1, 0)
			}
		}
	})
}

func BenchmarkFillRandomPreAllocated(b *testing.B) {
	rand.Seed(0)
	b.Run("Std", func(b *testing.B) {
		b.ReportAllocs()
		m := make(map[uint64]uint64, 1_000_000)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[rand.Uint64()+1] = 0
			}
		}
	})

	b.Run("IntInt", func(b *testing.B) {
		b.ReportAllocs()
		m := New(1_000_000, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(rand.Uint64()+1, 0)
			}
		}
	})
}

func BenchmarkLookupSequential(b *testing.B) {
	fillIntInt := func() *Map {
		m := New(2048, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(j, 0)
			}
		}
		return m
	}
	fillStd := func() map[uint64]uint64 {
		m := make(map[uint64]uint64, 2048)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[j] = 0
			}
		}
		return m
	}

	b.ReportAllocs()
	b.Run("Std", func(b *testing.B) {
		m := fillStd()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m[j]
			}
		}
		runtime.KeepAlive(dummy)
	})

	b.Run("IntInt", func(b *testing.B) {
		m := fillIntInt()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m.Get(j)
			}
		}
		runtime.KeepAlive(dummy)
	})
}

func BenchmarkLookupRandom(b *testing.B) {
	fillIntInt := func() *Map {
		m := New(2048, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(j, 0)
			}
		}
		return m
	}
	fillStd := func() map[uint64]uint64 {
		m := make(map[uint64]uint64, 2048)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[j] = 0
			}
		}
		return m
	}

	rand.Seed(0)
	b.ReportAllocs()
	b.Run("Std", func(b *testing.B) {
		m := fillStd()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m[uint64(rand.Int63n(1_000_000))+1]
			}
		}
		runtime.KeepAlive(dummy)
	})

	b.Run("IntInt", func(b *testing.B) {
		m := fillIntInt()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m.Get(uint64(rand.Int63n(1_000_000)) + 1)
			}
		}
		runtime.KeepAlive(dummy)
	})
}

func BenchmarkLookupNoHit(b *testing.B) {
	fillIntInt := func() *Map {
		m := New(2048, 0.6)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m.Put(j, 0)
			}
		}
		return m
	}
	fillStd := func() map[uint64]uint64 {
		m := make(map[uint64]uint64, 2048)
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				m[j] = 0
			}
		}
		return m
	}

	rand.Seed(0)
	b.ReportAllocs()
	b.Run("Std", func(b *testing.B) {
		m := fillStd()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m[uint64(rand.Int63n(1_000_000)+1_000_000)]
			}
		}
		runtime.KeepAlive(dummy)
	})

	b.Run("IntInt", func(b *testing.B) {
		m := fillIntInt()
		b.ReportAllocs()

		var dummy uint64
		for i := 0; i < b.N; i++ {
			for j := uint64(1); j < 1_000_000; j++ {
				dummy, _ = m.Get(uint64(rand.Int63n(1_000_000) + 1_000_000))
			}
		}
		runtime.KeepAlive(dummy)
	})
}
