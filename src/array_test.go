package Array

import (
	"testing"

	Error "github.com/golang-oop/error/src"
	Result "github.com/golang-oop/result/src"
)

// errResult builds a Result that reports HasError() == true so that the
// Each iterators short-circuit on it. Because Result.HasError() is defined as
// error.IsNull(), a fresh Result.New() already reports HasError() == true.
func errResult() Result.Interface {
	return Result.New()
}

// okResult builds a Result that reports HasError() == false, allowing Each to
// continue iterating. It carries a non-null error sentinel.
func okResult() Result.Interface {
	return Result.New(Result.WithError(Error.New("sentinel")))
}

func TestNew(t *testing.T) {
	if New() == nil {
		t.Fatal("New() returned nil")
	}
}

func TestPushPopFirstFetch(t *testing.T) {
	a := New()

	a.Push("a")
	r := a.Push("b")
	if r.HasError() == false {
		// Push returns a payload-bearing result; just ensure it is non-nil.
		_ = r.Payload()
	}

	if got := a.First().Payload(); got != "a" {
		t.Fatalf("First() = %v, want a", got)
	}
	if got := a.Fetch(1).Payload(); got != "b" {
		t.Fatalf("Fetch(1) = %v, want b", got)
	}

	if got := a.Pop().Payload(); got != "b" {
		t.Fatalf("Pop() = %v, want b", got)
	}
	if got := a.First().Payload(); got != "a" {
		t.Fatalf("First() after Pop = %v, want a", got)
	}
}

func TestClear(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push(2)
	a.Clear()
	// After Clear the slice is empty; pushing again should start from index 0.
	a.Push("x")
	if got := a.First().Payload(); got != "x" {
		t.Fatalf("First() after Clear = %v, want x", got)
	}
}

func TestCopy(t *testing.T) {
	a := New()
	if a.Copy() == nil {
		t.Fatal("Copy() returned nil")
	}
}

func TestMethodEach(t *testing.T) {
	a := New()
	a.Push(10)
	a.Push(20)
	a.Push(30)

	count := 0
	res := a.Each(func(i int, item interface{}) Result.Interface {
		count++
		return okResult()
	})
	if count != 3 {
		t.Fatalf("Each visited %d items, want 3", count)
	}
	// On a full iteration Each returns Result.New() (its completion result).
	if res == nil {
		t.Fatal("Each returned nil result")
	}

	// Short-circuit: stop on the first item via an error-reporting result.
	count = 0
	a.Each(func(i int, item interface{}) Result.Interface {
		count++
		return errResult()
	})
	if count != 1 {
		t.Fatalf("Each (short-circuit) visited %d items, want 1", count)
	}
}

func TestPackageEach(t *testing.T) {
	items := []interface{}{1, 2, 3}

	count := 0
	res := Each(items, func(i int, item interface{}) Result.Interface {
		count++
		return okResult()
	})
	if count != 3 {
		t.Fatalf("Each visited %d items, want 3", count)
	}
	if res == nil {
		t.Fatal("Each returned nil result")
	}

	count = 0
	Each(items, func(i int, item interface{}) Result.Interface {
		count++
		return errResult()
	})
	if count != 1 {
		t.Fatalf("Each (short-circuit) visited %d items, want 1", count)
	}
}
