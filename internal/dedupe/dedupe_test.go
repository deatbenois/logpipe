package dedupe

import (
	"fmt"
	"sync"
	"testing"
)

func TestIsDuplicate_FirstLineFalse(t *testing.T) {
	d := New(1)
	if d.IsDuplicate("src", "hello") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_ConsecutiveSameLine(t *testing.T) {
	d := New(1)
	d.IsDuplicate("src", "hello")
	if !d.IsDuplicate("src", "hello") {
		t.Fatal("second identical line should be a duplicate")
	}
}

func TestIsDuplicate_DifferentSources_Independent(t *testing.T) {
	d := New(1)
	d.IsDuplicate("a", "hello")
	d.IsDuplicate("b", "hello")
	// same line on different source should still be duplicate for that source
	if !d.IsDuplicate("b", "hello") {
		t.Fatal("same line on same source should be duplicate")
	}
	// but a new line on source a should not be duplicate
	if d.IsDuplicate("a", "world") {
		t.Fatal("different line should not be duplicate")
	}
}

func TestIsDuplicate_WindowGreaterThanOne(t *testing.T) {
	d := New(3)
	d.IsDuplicate("s", "a") // history: [a]
	d.IsDuplicate("s", "b") // history: [a, b]
	d.IsDuplicate("s", "c") // history: [a, b, c]
	// "a" is still within window of 3
	if !d.IsDuplicate("s", "a") {
		t.Fatal("line within window should be duplicate")
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	d := New(2)
	d.IsDuplicate("s", "a") // history: [a]
	d.IsDuplicate("s", "b") // history: [a, b]
	d.IsDuplicate("s", "c") // history: [b, c]  — "a" evicted
	if d.IsDuplicate("s", "a") {
		t.Fatal("evicted line should not be duplicate")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(1)
	d.IsDuplicate("src", "hello")
	d.Reset()
	if d.IsDuplicate("src", "hello") {
		t.Fatal("after reset, line should not be duplicate")
	}
}

func TestIsDuplicate_Concurrent(t *testing.T) {
	d := New(5)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			src := fmt.Sprintf("src-%d", n%5)
			line := fmt.Sprintf("line-%d", n%3)
			d.IsDuplicate(src, line)
		}(i)
	}
	wg.Wait() // should not race
}
