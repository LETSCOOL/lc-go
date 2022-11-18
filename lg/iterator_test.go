// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

import (
	"fmt"
	"testing"
)

// go test ./lg -v -run TestIterator
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

// go test ./lg -v -run TestIterateFunc
func TestIterateFunc(t *testing.T) {
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

// go test ./lg -v -run TestMapReduce
func TestMapReduce(t *testing.T) {
	t.Run("map-reduce", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		floatArray := Map(array, func(in int) float64 {
			return float64(in) * 0.3
		})
		result := Reduce(floatArray, "", func(v float64, lastResult string) string {
			return fmt.Sprintf("%s-%1.2f", lastResult, v)
		})
		if result != "-0.30-0.60-0.90-1.20-1.50-1.80-2.10-2.40-2.70" {
			t.Errorf("incorrect value: %s", result)
		}
	})
}
