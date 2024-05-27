package Array

import Result "github.com/golang-oop/result/src"

type Interface interface {
	Each(fn func(int, interface{}) Result.Interface) Result.Interface
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

func (d *data) Push(item interface{}) Result.Interface {
	temp := append(d.value, item)
	d.value = temp
	return Result.New(
		Result.WithPayload(d),
	)
}
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
		Result.WithPayload(d.value[len(d.value)]),
	)
}
func (d *data) Clear() Result.Interface {
	d.value = make([]interface{}, 0)
	return Result.New()
}
func (d data) Copy() Result.Interface {
	return Result.New()
}
