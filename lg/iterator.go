package lg

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
