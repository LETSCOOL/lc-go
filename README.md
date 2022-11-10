## lc-go
A golang library, include following packages:
1. lg
2. mongobj
3. dij

### lg (Language)

- Ife
    ```go
    // aka. text = (len(text) != 0) ? text : "some text"
    text = Ife(len(text) != 0, text, "some text")
    ```
    ```go
    // aka i = (i != 0) ? i : 123
    i = Ife(i != 0, i, 123)
    ```
- MakeIterator, MakeInverseIterator
    ```go 
    ints := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}
    iter := MakeIterator(ints)
    totalValues := 0
    for v, b, i := iter(); b; v, b, i = iter() {
        totalValues += v
    }
    ```
- IterateFunc, IterateFuncInversely
    ```go
    totalValues := 0
    IterateFuncInversely(ints, func(v int, i int) (stop bool) {
        totalValues += v
        return
    })
    ```
- map-reduce
    ```go
    array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
    floatArray := Map(array, func(in int) float64 {
        return float64(in) * 0.3
    })
    result := Reduce(floatArray, "", func(v float64, lastResult string) string {
        return fmt.Sprintf("%s-%1.2f", lastResult, v)
    })
    ```

### mongobj (MongoDb Object)
(omit)

### dij (Dependency Injection) - **draft**
Sample code
```go
type SampleApp struct {
	lib1 *SampleLib1 `di:"lib1"`
	lib2 *SampleLib2 `di:"lib2"`
}

type SampleLib1 struct {
	lib2 *SampleLib2 `di:"lib2"`
}

type SampleLib2 struct {
	val int `di:"val"`
}

// go test ./dij -v -run TestSample
func TestSample(t *testing.T) {
	t.Run("sample", func(t *testing.T) {
		appTyp := reflect.TypeOf(SampleApp{})
		ref := DependencyReference{"val": 123}
		inst, err := CreateInstance(appTyp, &ref, "^")
		if err != nil {
			t.Fatal(err)
		}
		if app, ok := inst.(*SampleApp); ok {
			if app.lib2 != app.lib1.lib2 {
				t.Errorf("incorrect injection, app.lib2(%v) != app.lib1.lib2(%v)\n", app.lib2, app.lib1.lib2)
			}
			if app.lib2.val != 123 {
				t.Errorf("incorrect injection, app.lib2.val(%d) != 123\n", app.lib2.val)
			}
		} else {
			t.Fatal("didn't create a correct instance, ", reflect.TypeOf(inst))
		}
	})
}
```
