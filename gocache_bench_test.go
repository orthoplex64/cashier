package gocache

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, NoExpiration)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(NoItemLimit, exp, 0)
	tc.Set("foo", "bar", DefaultExpiration)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo")
	}
}

func BenchmarkRWMutexMapGet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m["foo"]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetStruct(b *testing.B) {
	b.StopTimer()
	s := struct{ name string }{name: "foo"}
	m := map[interface{}]string{
		s: "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m[s]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetString(b *testing.B) {
	b.StopTimer()
	m := map[interface{}]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m["foo"]
		mu.RUnlock()
	}
}

func BenchmarkCacheGetConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, NoExpiration)
}

func benchmarkCacheGetConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(NoItemLimit, exp, 0)
	tc.Set("foo", "bar", DefaultExpiration)
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				tc.Get("foo")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRWMutexMapGetConcurrent(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				mu.RLock()
				_ = m["foo"]
				mu.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, NoExpiration)
}

func benchmarkCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	// This is the same as BenchmarkCacheGetConcurrent, but its result
	// can be compared against BenchmarkShardedCacheGetManyConcurrent
	// in sharded_test.go.
	b.StopTimer()
	n := 10000
	tc := New(NoItemLimit, exp, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		tc.Set(k, "bar", DefaultExpiration)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				tc.Get(k)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}

func BenchmarkCacheSetExpiring(b *testing.B) {
	benchmarkCacheSet(b, 5*time.Minute)
}

func BenchmarkCacheSetNotExpiring(b *testing.B) {
	benchmarkCacheSet(b, NoExpiration)
}

func benchmarkCacheSet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := New(NoItemLimit, exp, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", DefaultExpiration)
	}
}

func BenchmarkRWMutexMapSet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		mu.Unlock()
	}
}

func BenchmarkCacheSetDelete(b *testing.B) {
	b.StopTimer()
	tc := New(NoItemLimit, DefaultExpiration, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", DefaultExpiration)
		tc.Delete("foo")
	}
}

func BenchmarkRWMutexMapSetDelete(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		mu.Unlock()
		mu.Lock()
		delete(m, "foo")
		mu.Unlock()
	}
}

func BenchmarkCacheSetDeleteSingleLock(b *testing.B) {
	b.StopTimer()
	tc := New(NoItemLimit, DefaultExpiration, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.mu.Lock()
		tc.set("foo", "bar", DefaultExpiration)
		tc.delete("foo")
		tc.mu.Unlock()
	}
}

func BenchmarkRWMutexMapSetDeleteSingleLock(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		delete(m, "foo")
		mu.Unlock()
	}
}

func BenchmarkDeleteExpiredLoop(b *testing.B) {
	b.StopTimer()
	tc := New(NoItemLimit, 5*time.Minute, 0)
	tc.mu.Lock()
	for i := 0; i < 100000; i++ {
		tc.set(strconv.Itoa(i), "bar", DefaultExpiration)
	}
	tc.mu.Unlock()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.DeleteExpired()
	}
}

func TestGetWithExpiration(t *testing.T) {
	tc := New(NoItemLimit, DefaultExpiration, 0)

	a, expiration, found := tc.GetWithExpiration("a")
	if found || a != nil || !expiration.IsZero() {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, expiration, found := tc.GetWithExpiration("b")
	if found || b != nil || !expiration.IsZero() {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, expiration, found := tc.GetWithExpiration("c")
	if found || c != nil || !expiration.IsZero() {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.Set("a", 1, DefaultExpiration)
	tc.Set("b", "b", DefaultExpiration)
	tc.Set("c", 3.5, DefaultExpiration)
	tc.Set("d", 1, NoExpiration)
	tc.Set("e", 1, 50*time.Millisecond)

	x, expiration, found := tc.GetWithExpiration("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}
	if !expiration.IsZero() {
		t.Error("Expiration for a is not a zeroed time")
	}

	x, expiration, found = tc.GetWithExpiration("b")
	if !found {
		t.Error("b was not found while getting b2")
	}
	if x == nil {
		t.Error("x for b is nil")
	} else if b2 := x.(string); b2+"B" != "bB" {
		t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
	}
	if !expiration.IsZero() {
		t.Error("Expiration for b is not a zeroed time")
	}

	x, expiration, found = tc.GetWithExpiration("c")
	if !found {
		t.Error("c was not found while getting c2")
	}
	if x == nil {
		t.Error("x for c is nil")
	} else if c2 := x.(float64); c2+1.2 != 4.7 {
		t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
	}
	if !expiration.IsZero() {
		t.Error("Expiration for c is not a zeroed time")
	}

	x, expiration, found = tc.GetWithExpiration("d")
	if !found {
		t.Error("d was not found while getting d2")
	}
	if x == nil {
		t.Error("x for d is nil")
	} else if d2 := x.(int); d2+2 != 3 {
		t.Error("d (which should be 1) plus 2 does not equal 3; value:", d2)
	}
	if !expiration.IsZero() {
		t.Error("Expiration for d is not a zeroed time")
	}

	x, expiration, found = tc.GetWithExpiration("e")
	if !found {
		t.Error("e was not found while getting e2")
	}
	if x == nil {
		t.Error("x for e is nil")
	} else if e2 := x.(int); e2+2 != 3 {
		t.Error("e (which should be 1) plus 2 does not equal 3; value:", e2)
	}
	if expiration.UnixNano() != tc.items["e"].Expiration {
		t.Error("Expiration for e is not the correct time")
	}
	if expiration.UnixNano() < time.Now().UnixNano() {
		t.Error("Expiration for e is in the past")
	}
}
