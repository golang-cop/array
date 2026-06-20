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

## Access

`New()` → `Push`, `Pop`, `First`, `Last`, `Fetch(i)`, `Clear`, `Copy` — each
returns a `Result.Interface`.

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
