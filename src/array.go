package Array

import (
	"reflect"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

type Interface interface {
	Each(fn func(int, interface{}) Result.Interface) Result.Interface
	Map(fn func(int, interface{}) Result.Interface) Result.Interface
	Filter(pred func(int, interface{}) Result.Interface) Result.Interface
	Reduce(seed interface{}, fn func(acc, item interface{}) Result.Interface) Result.Interface
	Find(pred func(int, interface{}) Result.Interface) Result.Interface
	Any(pred func(int, interface{}) Result.Interface) Result.Interface
	All(pred func(int, interface{}) Result.Interface) Result.Interface
	Push(interface{}) Result.Interface
	Pop() Result.Interface
	First() Result.Interface
	Fetch(int) Result.Interface
	Last() Result.Interface
	Clear() Result.Interface
	Copy() Result.Interface
	Len() int
	IsEmpty() bool
	Contains(item interface{}) Result.Interface
	Insert(index int, item interface{}) Result.Interface
	Delete(index int) Result.Interface
	Reverse() Result.Interface
	IsNull() bool
}

type data struct {
	value []interface{}
}

func New() Interface {
	return &data{
		value: make([]interface{}, 0),
	}
}

/*
Iterate over an Array of interface{}
*/
func Each(
	items []interface{},
	fn func(int, interface{}) Result.Interface,
) Result.Interface {
	for index, item := range items {
		if result := fn(index, item); result.HasError() {
			return result
		}
	}
	return Result.New()
}

/*
Iterate over an Array of interface{}
and execute fn on each item
*/
func (d *data) Each(
	fn func(int, interface{}) Result.Interface,
) Result.Interface {
	for index, item := range d.value {
		if result := fn(index, item); result.HasError() {
			return result
		}
	}
	return Result.New()
}

/*
Add a new item to the end of the Array.
*/
func (d *data) Push(item interface{}) Result.Interface {
	temp := append(d.value, item)
	d.value = temp
	return Result.New(
		Result.WithPayload(d),
	)
}

/*
Remove the last item of the Array.
*/
func (d *data) Pop() Result.Interface {
	size := len(d.value)
	item := (d.value)[size-1]
	d.value = (d.value)[:size-1]
	return Result.New(
		Result.WithPayload(item),
	)
}

func (d data) First() Result.Interface {
	return Result.New(
		Result.WithPayload(d.value[0]),
	)
}
func (d data) Fetch(index int) Result.Interface {
	return Result.New(
		Result.WithPayload(d.value[index]),
	)
}
func (d data) Last() Result.Interface {
	return Result.New(
		Result.WithPayload(d.value[len(d.value)-1]),
	)
}
func (d *data) Clear() Result.Interface {
	d.value = make([]interface{}, 0)
	return Result.New()
}
func (d data) Copy() Result.Interface {
	out := &data{
		value: make([]interface{}, len(d.value)),
	}
	copy(out.value, d.value)
	return Result.New(
		Result.WithPayload(Interface(out)),
	)
}

// Len returns the number of elements in the Array.
func (d data) Len() int {
	return len(d.value)
}

// IsEmpty reports whether the Array has no elements.
func (d data) IsEmpty() bool {
	return len(d.value) == 0
}

// Contains reports, as a Result whose payload is a Go bool, whether item is
// present in the Array (compared with reflect.DeepEqual).
func (d data) Contains(item interface{}) Result.Interface {
	for _, candidate := range d.value {
		if reflect.DeepEqual(candidate, item) {
			return Result.New(Result.WithPayload(true))
		}
	}
	return Result.New(Result.WithPayload(false))
}

// Insert places item at index, shifting the tail right. A valid index ranges
// from 0 to Len() inclusive (Len() appends). Out-of-range indices yield an
// error Result; on success the payload is the Array itself.
func (d *data) Insert(index int, item interface{}) Result.Interface {
	if index < 0 || index > len(d.value) {
		return Result.New(
			Result.WithError(Error.New("Array.Insert: index out of range")),
		)
	}
	d.value = append(d.value, nil)
	copy(d.value[index+1:], d.value[index:])
	d.value[index] = item
	return Result.New(
		Result.WithPayload(Interface(d)),
	)
}

// Delete removes the element at index, shifting the tail left. Out-of-range
// indices yield an error Result; on success the payload is the removed item.
func (d *data) Delete(index int) Result.Interface {
	if index < 0 || index >= len(d.value) {
		return Result.New(
			Result.WithError(Error.New("Array.Delete: index out of range")),
		)
	}
	item := d.value[index]
	d.value = append(d.value[:index], d.value[index+1:]...)
	return Result.New(
		Result.WithPayload(item),
	)
}

// Reverse reverses the Array in place; the payload is the Array itself.
func (d *data) Reverse() Result.Interface {
	for i, j := 0, len(d.value)-1; i < j; i, j = i+1, j-1 {
		d.value[i], d.value[j] = d.value[j], d.value[i]
	}
	return Result.New(
		Result.WithPayload(Interface(d)),
	)
}

// IsNull reports that this is a real (non-null) Array.
func (d data) IsNull() bool {
	return false
}

/*
isTruthy decides whether a predicate's payload counts as a match.

It understands Go native booleans, any value with an IsTrue() bool method
(e.g. a Boolean.Interface — matched structurally to avoid importing boolean,
which would create an array→boolean→inspect→string→array cycle) and the
nil/absent case. Any other non-nil value is treated as truthy, mirroring the
"non-error / present" intuition.
*/
func isTruthy(payload interface{}) bool {
	switch p := payload.(type) {
	case nil:
		return false
	case bool:
		return p
	case interface{ IsTrue() bool }:
		return p.IsTrue()
	default:
		return true
	}
}

/*
Map applies fn to each item of an Array of interface{}, collecting the
non-error payloads into a brand new Array. It short-circuits, returning the
first error Result it encounters; on a full pass it returns a Result whose
payload is the new Array.Interface.
*/
func Map(
	items []interface{},
	fn func(int, interface{}) Result.Interface,
) Result.Interface {
	out := New()
	for index, item := range items {
		result := fn(index, item)
		if result.HasError() {
			return result
		}
		out.Push(result.Payload())
	}
	return Result.New(
		Result.WithPayload(out),
	)
}

/*
Map applies fn to each item and collects the resulting payloads into a new
Array, short-circuiting on the first error Result.
*/
func (d *data) Map(
	fn func(int, interface{}) Result.Interface,
) Result.Interface {
	return Map(d.value, fn)
}

/*
Filter keeps the items of an Array of interface{} whose pred returns a
non-error, truthy Result, collecting them into a new Array. It short-circuits
on the first error Result; on a full pass it returns a Result whose payload is
the new Array.Interface.
*/
func Filter(
	items []interface{},
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	out := New()
	for index, item := range items {
		result := pred(index, item)
		if result.HasError() {
			return result
		}
		if isTruthy(result.Payload()) {
			out.Push(item)
		}
	}
	return Result.New(
		Result.WithPayload(out),
	)
}

/*
Filter keeps the items whose pred returns a non-error, truthy Result,
short-circuiting on the first error Result.
*/
func (d *data) Filter(
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	return Filter(d.value, pred)
}

/*
Reduce performs a left fold over an Array of interface{}, threading an
accumulator through fn while propagating Result errors. It short-circuits on
the first error Result; on a full pass it returns a Result whose payload is the
final accumulator value.
*/
func Reduce(
	items []interface{},
	seed interface{},
	fn func(acc, item interface{}) Result.Interface,
) Result.Interface {
	acc := seed
	for _, item := range items {
		result := fn(acc, item)
		if result.HasError() {
			return result
		}
		acc = result.Payload()
	}
	return Result.New(
		Result.WithPayload(acc),
	)
}

/*
Reduce performs a left fold over the Array, threading an accumulator through
fn and short-circuiting on the first error Result.
*/
func (d *data) Reduce(
	seed interface{},
	fn func(acc, item interface{}) Result.Interface,
) Result.Interface {
	return Reduce(d.value, seed, fn)
}

/*
Find returns the first item of an Array of interface{} for which pred returns a
non-error, truthy Result, as the payload of a successful Result. It
short-circuits on the first error Result; when no item matches it returns a
Result carrying a not-found Error.
*/
func Find(
	items []interface{},
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	for index, item := range items {
		result := pred(index, item)
		if result.HasError() {
			return result
		}
		if isTruthy(result.Payload()) {
			return Result.New(
				Result.WithPayload(item),
			)
		}
	}
	return Result.New(
		Result.WithError(Error.New("Array.Find: no matching item")),
	)
}

/*
Find returns the first item for which pred returns a non-error, truthy Result,
short-circuiting on the first error Result and reporting a not-found Error when
nothing matches.
*/
func (d *data) Find(
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	return Find(d.value, pred)
}

/*
Any reports whether at least one item of an Array of interface{} satisfies
pred. It returns a Result whose payload is a Go bool. It short-circuits: a
pred error is propagated (the error Result is returned), the first truthy item
yields a true payload, and an empty Array (or no match) yields a false payload.
*/
func Any(
	items []interface{},
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	for index, item := range items {
		result := pred(index, item)
		if result.HasError() {
			return result
		}
		if isTruthy(result.Payload()) {
			return Result.New(Result.WithPayload(true))
		}
	}
	return Result.New(Result.WithPayload(false))
}

/*
Any reports whether at least one item satisfies pred, short-circuiting on the
first truthy match (or a pred error).
*/
func (d *data) Any(
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	return Any(d.value, pred)
}

/*
All reports whether every item of an Array of interface{} satisfies pred. It
returns a Result whose payload is a Go bool. It short-circuits: a pred error is
propagated, the first falsy item yields a false payload, and an empty Array
yields a true payload (vacuous truth).
*/
func All(
	items []interface{},
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	for index, item := range items {
		result := pred(index, item)
		if result.HasError() {
			return result
		}
		if !isTruthy(result.Payload()) {
			return Result.New(Result.WithPayload(false))
		}
	}
	return Result.New(Result.WithPayload(true))
}

/*
All reports whether every item satisfies pred, short-circuiting on the first
failing item (or a pred error).
*/
func (d *data) All(
	pred func(int, interface{}) Result.Interface,
) Result.Interface {
	return All(d.value, pred)
}

// null is the Null-Object variant of an Array: an empty, immutable placeholder
// that honours the full Interface without ever being nil. Mutating methods are
// no-ops returning a successful Result; lookups and queries return empty/false/
// zero values.
type null struct{}

// Null returns the Null-Object Array.
func Null() Interface {
	return &null{}
}

func (n *null) Each(fn func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New()
}

func (n *null) Map(fn func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New(Result.WithPayload(Null()))
}

func (n *null) Filter(pred func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New(Result.WithPayload(Null()))
}

func (n *null) Reduce(
	seed interface{},
	fn func(acc, item interface{}) Result.Interface,
) Result.Interface {
	return Result.New(Result.WithPayload(seed))
}

func (n *null) Find(pred func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.Find: no matching item")),
	)
}

func (n *null) Any(pred func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New(Result.WithPayload(false))
}

func (n *null) All(pred func(int, interface{}) Result.Interface) Result.Interface {
	return Result.New(Result.WithPayload(true))
}

func (n *null) Push(item interface{}) Result.Interface {
	return Result.New(Result.WithPayload(Interface(n)))
}

func (n *null) Pop() Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.Pop: empty array")),
	)
}

func (n *null) First() Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.First: empty array")),
	)
}

func (n *null) Fetch(index int) Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.Fetch: index out of range")),
	)
}

func (n *null) Last() Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.Last: empty array")),
	)
}

func (n *null) Clear() Result.Interface {
	return Result.New()
}

func (n *null) Copy() Result.Interface {
	return Result.New(Result.WithPayload(Null()))
}

func (n *null) Len() int {
	return 0
}

func (n *null) IsEmpty() bool {
	return true
}

func (n *null) Contains(item interface{}) Result.Interface {
	return Result.New(Result.WithPayload(false))
}

func (n *null) Insert(index int, item interface{}) Result.Interface {
	return Result.New(Result.WithPayload(Interface(n)))
}

func (n *null) Delete(index int) Result.Interface {
	return Result.New(
		Result.WithError(Error.New("Array.Delete: index out of range")),
	)
}

func (n *null) Reverse() Result.Interface {
	return Result.New(Result.WithPayload(Interface(n)))
}

// IsNull reports that this is the null Array.
func (n *null) IsNull() bool {
	return true
}
