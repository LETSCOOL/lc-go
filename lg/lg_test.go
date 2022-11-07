package lg

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// go test ./lg -v -run TestCallWithParameter
func TestCallWithParameter(t *testing.T) {
	t.Run("directly", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			go func() {
				time.Sleep(time.Microsecond * 100)
				fmt.Println(i)
			}()
		}
		time.Sleep(time.Microsecond * 200)
	})

	t.Run("pass by value", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			go func(i int) {
				time.Sleep(time.Microsecond * 100)
				fmt.Println(i)
			}(i)
		}
		time.Sleep(time.Microsecond * 200)
	})

	t.Run("pass by reference", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			go func(pI *int) {
				time.Sleep(time.Microsecond * 100)
				fmt.Println(*pI)
			}(&i)
		}
		time.Sleep(time.Microsecond * 200)
	})
}

type B01 struct {
	counter int
}

func (b B01) PrintTypeByValue() {
	b.counter = b.counter + 1
	fmt.Printf("PrintTypeByValue: %v, counter: %d\n", reflect.TypeOf(b), b.counter)
}

func (b *B01) PrintTypeByPointer() {
	b.counter = b.counter + 1
	fmt.Printf("PrintTypeByPointer: %v, counter: %d\n", reflect.TypeOf(b), b.counter)
}

type E01 struct {
	B01
}

type E02 struct {
	*B01
}

// go test ./lg -v -run TestEmbeddedField
func TestEmbeddedField(t *testing.T) {
	t.Run("E01", func(t *testing.T) {
		e01 := E01{}
		fmt.Printf("===e01 by value=== (counter: %d)\n", e01.counter)
		e01.PrintTypeByValue()
		fmt.Printf("\tcounter: %d\n", e01.counter)
		e01.PrintTypeByPointer()
		fmt.Printf("\tcounter: %d\n", e01.counter)
		fmt.Printf("===e01 by pointer=== (counter: %d)\n", e01.counter)
		pB01 := &e01
		pB01.PrintTypeByValue()
		fmt.Printf("\tcounter: %d\n", e01.counter)
		pB01.PrintTypeByPointer()
		fmt.Printf("\tcounter: %d\n", e01.counter)
	})

	t.Run("E02", func(t *testing.T) {
		e02 := E02{
			&B01{},
		}
		fmt.Printf("===e02 by value=== (counter: %d)\n", e02.counter)
		e02.PrintTypeByValue()
		fmt.Printf("\tcounter: %d\n", e02.counter)
		e02.PrintTypeByPointer()
		fmt.Printf("\tcounter: %d\n", e02.counter)
		fmt.Printf("===e02 by pointer=== (counter: %d)\n", e02.counter)
		pB01 := &e02
		pB01.PrintTypeByValue()
		fmt.Printf("\tcounter: %d\n", e02.counter)
		pB01.PrintTypeByPointer()
		fmt.Printf("\tcounter: %d\n", e02.counter)
	})
}
