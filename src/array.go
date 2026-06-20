package Array

import (
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
	return Result.New()
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
