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
	a.Push("a")
	a.Push("b")

	res := a.Copy()
	if res == nil || res.HasError() {
		t.Fatal("Copy() returned nil or error")
	}
	cp := res.Payload().(Interface)

	// Same contents.
	if cp.Len() != 2 || cp.First().Payload() != "a" || cp.Last().Payload() != "b" {
		t.Fatalf("Copy contents = %v..%v len %d", cp.First().Payload(), cp.Last().Payload(), cp.Len())
	}

	// Independence: mutating the copy must not affect the original.
	cp.Push("c")
	if a.Len() != 2 {
		t.Fatalf("original len = %d after mutating copy, want 2", a.Len())
	}
	cp.Delete(0)
	if a.First().Payload() != "a" {
		t.Fatalf("original First() = %v after deleting from copy, want a", a.First().Payload())
	}
}

func TestLenIsEmpty(t *testing.T) {
	a := New()
	if a.Len() != 0 || !a.IsEmpty() {
		t.Fatal("new array should be empty with len 0")
	}
	a.Push(1)
	if a.Len() != 1 || a.IsEmpty() {
		t.Fatal("array with one item should have len 1 and not be empty")
	}
}

func TestContains(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push([]int{2, 3})
	if a.Contains(1).Payload() != true {
		t.Fatal("Contains(1) want true")
	}
	if a.Contains([]int{2, 3}).Payload() != true {
		t.Fatal("Contains deep slice want true")
	}
	if a.Contains(99).Payload() != false {
		t.Fatal("Contains(99) want false")
	}
}

func TestInsert(t *testing.T) {
	a := New()
	a.Push("a")
	a.Push("c")
	// Insert in the middle.
	if r := a.Insert(1, "b"); r.HasError() {
		t.Fatal("Insert(1) unexpected error")
	}
	if a.Fetch(1).Payload() != "b" || a.Len() != 3 {
		t.Fatalf("after Insert middle = %v len %d", a.Fetch(1).Payload(), a.Len())
	}
	// Insert at the end (index == Len()).
	if r := a.Insert(a.Len(), "d"); r.HasError() {
		t.Fatal("Insert at end unexpected error")
	}
	if a.Last().Payload() != "d" {
		t.Fatalf("after Insert end Last() = %v, want d", a.Last().Payload())
	}
	// Out of range.
	if !a.Insert(-1, "x").HasError() {
		t.Fatal("Insert(-1) want error")
	}
	if !a.Insert(a.Len()+1, "x").HasError() {
		t.Fatal("Insert(>Len) want error")
	}
}

func TestDelete(t *testing.T) {
	a := New()
	a.Push("a")
	a.Push("b")
	a.Push("c")
	r := a.Delete(1)
	if r.HasError() || r.Payload() != "b" {
		t.Fatalf("Delete(1) payload = %v", r.Payload())
	}
	if a.Len() != 2 || a.Fetch(1).Payload() != "c" {
		t.Fatalf("after Delete = len %d, [1]=%v", a.Len(), a.Fetch(1).Payload())
	}
	if !a.Delete(-1).HasError() {
		t.Fatal("Delete(-1) want error")
	}
	if !a.Delete(a.Len()).HasError() {
		t.Fatal("Delete(Len) want error")
	}
}

func TestReverse(t *testing.T) {
	a := New()
	a.Push(1)
	a.Push(2)
	a.Push(3)
	if r := a.Reverse(); r.HasError() {
		t.Fatal("Reverse unexpected error")
	}
	if a.First().Payload() != 3 || a.Last().Payload() != 1 {
		t.Fatalf("after Reverse = %v..%v, want 3..1", a.First().Payload(), a.Last().Payload())
	}
}

func TestIsNull(t *testing.T) {
	if New().IsNull() {
		t.Fatal("real array IsNull() want false")
	}
	if !Null().IsNull() {
		t.Fatal("Null() IsNull() want true")
	}
}

func TestNullVariant(t *testing.T) {
	n := Null()

	// Queries.
	if n.Len() != 0 || !n.IsEmpty() {
		t.Fatal("null Len/IsEmpty wrong")
	}
	if n.Contains(1).Payload() != false {
		t.Fatal("null Contains want false")
	}

	// Mutating no-ops succeed.
	for _, r := range []Result.Interface{
		n.Push(1), n.Clear(), n.Insert(0, 1), n.Reverse(), n.Copy(),
		n.Map(func(int, interface{}) Result.Interface { return okResult() }),
		n.Filter(func(int, interface{}) Result.Interface { return okResult() }),
		n.Each(func(int, interface{}) Result.Interface { return okResult() }),
		n.Reduce(0, func(_, _ interface{}) Result.Interface { return okResult() }),
	} {
		if r.HasError() {
			t.Fatal("null mutating/iter op unexpectedly errored")
		}
	}

	// Reduce returns seed; Any false; All true.
	if n.Reduce(7, func(_, _ interface{}) Result.Interface { return okResult() }).Payload() != 7 {
		t.Fatal("null Reduce want seed 7")
	}
	if n.Any(func(int, interface{}) Result.Interface { return okResult() }).Payload() != false {
		t.Fatal("null Any want false")
	}
	if n.All(func(int, interface{}) Result.Interface { return okResult() }).Payload() != true {
		t.Fatal("null All want true")
	}

	// Map/Filter/Copy payloads are null Arrays.
	if !n.Map(func(int, interface{}) Result.Interface { return okResult() }).Payload().(Interface).IsNull() {
		t.Fatal("null Map payload should be null array")
	}
	if !n.Filter(func(int, interface{}) Result.Interface { return okResult() }).Payload().(Interface).IsNull() {
		t.Fatal("null Filter payload should be null array")
	}
	if !n.Copy().Payload().(Interface).IsNull() {
		t.Fatal("null Copy payload should be null array")
	}
	if !n.Push(1).Payload().(Interface).IsNull() {
		t.Fatal("null Push payload should be null array")
	}
	if !n.Insert(0, 1).Payload().(Interface).IsNull() {
		t.Fatal("null Insert payload should be null array")
	}
	if !n.Reverse().Payload().(Interface).IsNull() {
		t.Fatal("null Reverse payload should be null array")
	}

	// Lookups/queries error or miss.
	for _, r := range []Result.Interface{
		n.Pop(), n.First(), n.Fetch(0), n.Last(), n.Delete(0), n.Find(func(int, interface{}) Result.Interface { return okResult() }),
	} {
		if !r.HasError() {
			t.Fatal("null lookup op should error")
		}
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
