package lg

import "testing"

// go test ./lg -v
func TestIterator(t *testing.T) {
	ints := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}
	t.Run("current order", func(t *testing.T) {
		iter := MakeIterator(ints)
		totalValues := 0
		totalIndexes := 0
		for v, b, i := iter(); b; v, b, i = iter() {
			totalValues += v
			totalIndexes += i
			//t.Logf("%d. %d", i, v)
		}
		if totalIndexes*10 != totalValues {
			t.Errorf("error: %d * 10 != %d", totalIndexes, totalValues)
		}
		if totalValues != 90*10/2 {
			t.Errorf("error: %d != %d", totalValues, 90*10/2)
		}
	})

	t.Run("inverse order", func(t *testing.T) {
		iter := MakeInverseIterator(ints)
		totalValues := 0
		totalIndexes := 0
		for v, b, i := iter(); b; v, b, i = iter() {
			totalValues += v
			totalIndexes += i
			//t.Logf("%d. %d", i, v)
		}
		if totalIndexes*10 != totalValues {
			t.Errorf("error: %d * 10 != %d", totalIndexes, totalValues)
		}
		if totalValues != 90*10/2 {
			t.Errorf("error: %d != %d", totalValues, 90*10/2)
		}
	})
}

func TestIterate(t *testing.T) {
	ints := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}
	t.Run("current order", func(t *testing.T) {
		totalValues := 0
		totalIndexes := 0
		IterateFunc(ints, func(v int, i int) (stop bool) {
			totalValues += v
			totalIndexes += i
			//t.Logf("%d * 10 ?= %d", totalIndexes, totalValues)
			return
		})
		if totalIndexes*10 != totalValues {
			t.Errorf("error: %d * 10 != %d", totalIndexes, totalValues)
		}
		if totalValues != 90*10/2 {
			t.Errorf("error: %d != %d", totalValues, 90*10/2)
		}
	})

	t.Run("current order", func(t *testing.T) {
		totalValues := 0
		totalIndexes := 0
		IterateFuncInversely(ints, func(v int, i int) (stop bool) {
			totalValues += v
			totalIndexes += i
			//t.Logf("%d * 10 ?= %d", totalIndexes, totalValues)
			return
		})
		if totalIndexes*10 != totalValues {
			t.Errorf("error: %d * 10 != %d", totalIndexes, totalValues)
		}
		if totalValues != 90*10/2 {
			t.Errorf("error: %d != %d", totalValues, 90*10/2)
		}
	})
}
