// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

func Map[IN string | bool | int | any, OUT string | bool | int | any](list []IN, f func(v IN) OUT) []OUT {
	out := make([]OUT, 0, len(list))
	for _, i := range list {
		o := f(i)
		out = append(out, o)
	}
	return out
}

func Reduce[IN any, OUT any](list []IN, initialOutputValue OUT, f func(v IN, lastResult OUT) OUT) OUT {
	out := initialOutputValue
	for _, v := range list {
		out = f(v, out)
	}
	return out
}

func Filter[T any](items []T, f func(v T) (is bool)) []T {
	filteredItems := make([]T, 0)

	for _, item := range items {
		if f(item) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func FilterFirst[T any](items []T, f func(v T) (is bool)) (val T, exists bool) {
	for _, item := range items {
		if f(item) {
			return item, true
		}
	}
	return
}

func MakeIterator[T any](items []T) func() (T, bool, int) {
	index := 0
	return func() (val T, ok bool, i int) {
		i = index
		if index >= len(items) {
			return
		}
		val, ok = items[index], true
		index++
		return
	}
}

func MakeInverseIterator[T any](items []T) func() (T, bool, int) {
	index := len(items) - 1
	return func() (val T, ok bool, i int) {
		i = index
		if index < 0 {
			return
		}
		val, ok = items[index], true
		index--
		return
	}
}

func IterateFunc[T any](items []T, f func(v T, index int) (stop bool)) {
	for index, value := range items {
		stop := f(value, index)
		if stop {
			break
		}
	}
}

func IterateFuncInversely[T any](items []T, f func(v T, index int) (stop bool)) {
	for index := len(items) - 1; index >= 0; index-- {
		stop := f(items[index], index)
		if stop {
			break
		}
	}
}
