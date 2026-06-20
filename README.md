<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/array" width="720"></p>

# array

[![ci](https://github.com/go-composites/array/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/array/actions/workflows/ci.yml)

The ordered-collection composite of [go-composites](https://github.com/go-composites).
An `Array` holds `interface{}` items and exposes both stack-like access and
functional combinators, all returning a [`Result`](https://github.com/go-composites/result)
so failures (and short-circuits) are *values*, never panics or `nil`.

## Install

```sh
go get github.com/go-composites/array
```

## Access & mutation

`New()` returns an empty `Array`. Unless noted, each method returns a
`Result.Interface`.

| method | result payload | notes |
| --- | --- | --- |
| `Push(item)` | the `Array` | append to the end |
| `Pop()` | the removed item | error `Result` on an empty `Array` |
| `First()` | the first item | |
| `Last()` | the last item | |
| `Fetch(i)` | the item at `i` | |
| `Insert(i, item)` | the `Array` | shift the tail right; valid `i` is `0..Len()` (append at `Len()`); out-of-range → error `Result` |
| `Delete(i)` | the removed item | shift the tail left; out-of-range → error `Result` |
| `Reverse()` | the `Array` | reverse in place |
| `Clear()` | completion `Result` | drop all elements |
| `Copy()` | a new `Array` | **deep copy** of the backing slice (independent of the receiver) |

## Queries

| method | returns | notes |
| --- | --- | --- |
| `Len()` | Go `int` | element count |
| `IsEmpty()` | Go `bool` | `Len() == 0` |
| `Contains(item)` | `Result` of Go `bool` | membership, via `reflect.DeepEqual` |
| `IndexOf(item)` | `Result` of Go `int` | index of the first match, or `-1` |
| `Slice(start, end)` | a new independent `Array` | elements `[start:end)`; out-of-range → error `Result` |
| `Sort(less)` | the `Array` | stable sort in place (`sort.SliceStable`) |

## Combinators

All take a predicate/function returning `Result.Interface` and **short-circuit on
the first error Result** (propagating it), mirroring `Each`:

| method | result payload | notes |
| --- | --- | --- |
| `Each(fn)` | completion `Result` | iterate; stop on first error |
| `Map(fn)` | a new `Array` | collect each fn payload |
| `Filter(pred)` | a new `Array` | keep items whose pred is truthy |
| `Reduce(seed, fn)` | final accumulator | left fold threading the accumulator |
| `Find(pred)` | first matching item | not-found → a `Result` carrying an error |
| `Any(pred)` | Go `bool` | true on first truthy item (vacuously false on empty) |
| `All(pred)` | Go `bool` | false on first falsy item (vacuously true on empty) |

Truthiness (for `Filter`/`Find`/`Any`/`All`) accepts a Go `bool`, any value with
an `IsTrue() bool` method (e.g. a `Boolean.Interface`, matched *structurally* so
`array` need not import — and cycle with — `boolean`), `nil` (falsy), and treats
any other non-nil payload as truthy.

## Null-Object

`Null()` returns the never-nil Null-Object `Array`: an empty, immutable
placeholder honouring the full `Interface`. `IsNull()` reports `true` for it and
`false` for every concrete `Array`. Mutating methods are successful no-ops and
queries return empty/`false`/zero values.

```go
a := Array.New()
a.Push(1); a.Push(2); a.Push(3)

doubled := a.Map(func(_ int, item interface{}) Result.Interface {
    return Result.New(Result.WithPayload(item.(int) * 2))
}) // payload: Array [2 4 6]

sum := a.Reduce(0, func(acc, item interface{}) Result.Interface {
    return Result.New(Result.WithPayload(acc.(int) + item.(int)))
}) // payload: 6
```

Each method has a package-level twin (`Array.Map(items, fn)`, …) operating on a
raw `[]interface{}`.

## License

BSD-3-Clause © the go-composites/array authors.
