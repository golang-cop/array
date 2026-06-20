package Array

import (
	"testing"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

// errResult builds a Result that reports HasError() == true so the Each
// iterators short-circuit on it: Result.HasError() is defined as
// !error.IsNull(), so a real (non-null) error makes it true.
func errResult() Result.Interface {
	return Result.New(Result.WithError(Error.New("sentinel")))
}

// okResult builds a Result that reports HasError() == false, allowing Each to
// continue iterating: a fresh Result.New() defaults to a NullError, whose
// IsNull() is true, so HasError() is false.
func okResult() Result.Interface {
	return Result.New()
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

// truthy is a tiny structural value exercising isTruthy's IsTrue() case.
type truthy struct{ v bool }

func (t truthy) IsTrue() bool { return t.v }

func TestLast(t *testing.T) {
	a := New()
	a.Push("x")
	a.Push("y")
	if got := a.Last().Payload(); got != "y" {
		t.Fatalf("Last() = %v, want y", got)
	}
}

func TestMap(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push(2)
	a.Push(3)
	res := a.Map(func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(item.(int) * 2))
	})
	if res.HasError() {
		t.Fatal("Map: unexpected error")
	}
	out := res.Payload().(Interface)
	if out.First().Payload() != 2 || out.Last().Payload() != 6 {
		t.Fatalf("Map = %v..%v, want 2..6", out.First().Payload(), out.Last().Payload())
	}
	if !a.Map(func(int, interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("Map: expected error propagation")
	}
}

func TestFilterAndIsTruthy(t *testing.T) {
	a := New()
	for _, v := range []string{"bt", "bf", "it", "def", "nilp"} {
		a.Push(v)
	}
	res := a.Filter(func(_ int, item interface{}) Result.Interface {
		switch item {
		case "bt":
			return Result.New(Result.WithPayload(true)) // bool true
		case "bf":
			return Result.New(Result.WithPayload(false)) // bool false
		case "it":
			return Result.New(Result.WithPayload(truthy{true})) // IsTrue()
		case "def":
			return Result.New(Result.WithPayload(42)) // default -> true
		default:
			return Result.New(Result.WithPayload(nil)) // nil -> false
		}
	})
	if res.HasError() {
		t.Fatal("Filter: unexpected error")
	}
	out := res.Payload().(Interface)
	// kept: bt, it, def  (3); dropped: bf, nilp
	if out.First().Payload() != "bt" || out.Last().Payload() != "def" {
		t.Fatalf("Filter kept %v..%v, want bt..def", out.First().Payload(), out.Last().Payload())
	}
	if !a.Filter(func(int, interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("Filter: expected error propagation")
	}
}

func TestReduce(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push(2)
	a.Push(3)
	res := a.Reduce(0, func(acc, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(acc.(int) + item.(int)))
	})
	if res.HasError() || res.Payload() != 6 {
		t.Fatalf("Reduce = %v, want 6", res.Payload())
	}
	if !a.Reduce(0, func(_, _ interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("Reduce: expected error propagation")
	}
}

func TestFind(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push(2)
	a.Push(3)
	eq2 := func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(item.(int) == 2))
	}
	if got := a.Find(eq2).Payload(); got != 2 {
		t.Fatalf("Find = %v, want 2", got)
	}
	// no match -> error
	if !a.Find(func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(false))
	}).HasError() {
		t.Fatal("Find: expected not-found error")
	}
	// pred error propagates
	if !a.Find(func(int, interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("Find: expected error propagation")
	}
}

func TestAnyAll(t *testing.T) {
	a := New()
	a.Push(2)
	a.Push(4)
	a.Push(6)
	even := func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(item.(int)%2 == 0))
	}
	isFour := func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(item.(int) == 4))
	}
	if a.Any(isFour).Payload() != true {
		t.Fatal("Any: expected true")
	}
	if a.Any(func(_ int, item interface{}) Result.Interface {
		return Result.New(Result.WithPayload(false))
	}).Payload() != false {
		t.Fatal("Any: expected false")
	}
	if a.All(even).Payload() != true {
		t.Fatal("All: expected true")
	}
	if a.All(isFour).Payload() != false {
		t.Fatal("All: expected false (not all 4)")
	}
	// pred error propagates in both
	if !a.Any(func(int, interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("Any: expected error propagation")
	}
	if !a.All(func(int, interface{}) Result.Interface { return errResult() }).HasError() {
		t.Fatal("All: expected error propagation")
	}
}
