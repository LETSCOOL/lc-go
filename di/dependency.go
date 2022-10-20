package di

import (
	"fmt"
	. "lc-go/lg"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

type DependencyKey string
type Lifecycle int

const (
	TagName = "di"
)

//const (
//	Transient Lifecycle = iota // per request
//	Singleton                  // per device
//	Scoped                     // per goroutine
//)

type Injectable interface {
	GetInjection() map[DependencyKey]any
	//RunInitialization()
}

func MethodByName(typ reflect.Type, fnName string) (reflect.Method, bool) {
	method, ok := typ.MethodByName(fnName)
	if ok {
		return method, ok
	}
	ptrTyp := reflect.PointerTo(typ)
	return ptrTyp.MethodByName(fnName)
}

func ToInstance(typ reflect.Type) (any, error) {
	if typ.Kind() == reflect.Pointer {
		log.Fatal("Don't use Pointer type")
	}
	instValue := reflect.New(typ)
	instIf := instValue.Interface()
	/*dep, ok := instIf.(Injectable)
	if !ok {
		return dep, errors.New("the type does not implement Injectable interface")
	}*/

	if _, ok := MethodByName(typ, "RunInitialization"); !ok {
		return nil, fmt.Errorf("the type(%v) does not implement a method with name 'RunInitialization'", typ)
	}

	return instIf, nil
}

func Invoke(inst interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, 0, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(inst).MethodByName(methodName).Call(inputs)
}

// parseDiTag parses the tag with 'di' key.
//
// TODO: complete the algorithm.
func parseDiTag(tag string) (name string) {
	segments := Map(strings.Split(tag, ","), func(in string) string {
		return strings.TrimSpace(in)
	})

	if len(segments) > 0 {
		name = segments[0]
	} else {
		name = ""
	}
	return
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func InitialDependency(typ reflect.Type, inst any, reference *map[DependencyKey]any) error {
	TypeOfType := reflect.TypeOf(reflect.TypeOf(struct{}{}))
	initMethod, _ := MethodByName(typ, "RunInitialization")
	methodTyp := initMethod.Type
	//log.Println(methodTyp.NumIn())
	if methodTyp.NumIn() <= 0 {
		// only receiver
		//log.Println("Only receiver", reflect.ValueOf(inst).MethodByName("RunInitialization").IsValid())
		reflect.ValueOf(inst).MethodByName("RunInitialization").Call(nil)
	} else {
		numMethodIn := methodTyp.NumIn()
		args := make([]reflect.Value, 0, numMethodIn-1)
		//log.Println("numMethodIn", numMethodIn)
		for i := 1; i < methodTyp.NumIn(); i++ {
			inTyp := methodTyp.In(i)
			if inTyp.Kind() != reflect.Struct {
				//log.Println(inTyp.Kind())
				return fmt.Errorf("initialization function support struct parameters(%s) only", inTyp.Name())
			}
			//log.Println("In", i, inTyp, inTyp.Kind())
			inValuePtr := reflect.New(inTyp)
			inValue := inValuePtr.Elem()
			//log.Println(reflect.TypeOf(inValue.Interface()))
			// TODO: init inValue
			for j := 0; j < inTyp.NumField(); j++ {
				fieldSpec := inTyp.Field(j)
				name := parseDiTag(fieldSpec.Tag.Get(TagName))
				//log.Println("Field name: ", name, fieldSpec)
				name = Ife(name == "", fieldSpec.Name, name)
				//log.Println("Field name: ", name)
				var refValue any
				if v, existing := (*reference)[DependencyKey(name)]; existing {
					if reflect.TypeOf(v) == TypeOfType {
						inst, err := CreateDependency(v.(reflect.Type), reference)
						if err != nil {
							return err
						}
						(*reference)[DependencyKey(name)] = inst
						refValue = inst
					} else {
						refValue = v
					}
				} else {
					if fieldSpec.Type.Kind() == reflect.Pointer && fieldSpec.Type.Elem().Kind() == reflect.Struct {
						underlyingType := fieldSpec.Type.Elem()
						inst, err := CreateDependency(underlyingType, reference)
						if err != nil {
							return err
						}
						(*reference)[DependencyKey(name)] = inst
						refValue = inst
					} else {
						return fmt.Errorf("non available reference(%s) for injection", name)
					}
				}
				if fieldSpec.Type == reflect.TypeOf(refValue) {
					field := inValue.FieldByName(fieldSpec.Name)
					firstChar := fieldSpec.Name[0]
					if firstChar >= 'A' && firstChar <= 'Z' {
						field.Set(reflect.ValueOf(refValue))
					} else {
						SetUnexportedField(field, refValue)
					}
				} else {
					return fmt.Errorf("filed(%s)'s declaration(%v) and value(%v) are not same type", fieldSpec.Name, fieldSpec.Type, reflect.TypeOf(refValue))
				}
			}

			args = append(args, inValue)
		}

		//Invoke(inst, "RunInitialization", args...)
		//log.Println("args num", len(args))
		reflect.ValueOf(inst).MethodByName("RunInitialization").Call(args)
	}

	return nil
}

// CreateDependency create new instance for dependency injection
//
// TODO: Before injecting, should the instance be initialized completely?
func CreateDependency(rootTyp reflect.Type, reference *map[DependencyKey]any) (any, error) {
	// fmt.Println("***", rootTyp, "***", "start")
	inst, err := ToInstance(rootTyp)
	if err != nil {
		return nil, err
	}

	if dep, ok := inst.(Injectable); ok {
		injection := dep.GetInjection()
		for injKey, injValue := range injection {
			if injTyp, ok := injValue.(reflect.Type); ok {
				// type, create and inject
				existedRef, isExisted := (*reference)[injKey]
				if isExisted {
					if reflect.TypeOf(existedRef) == injTyp {
						// same type and already created, do nothing
						// fmt.Println("***", rootTyp, "***", "CreateDependency: exist injKey", injKey, existedRef)
					} else {
						return nil, fmt.Errorf("different values(%v)/types(%v) refer to same key(%v)", existedRef, injTyp, injKey)
					}
				} else if injTyp.Kind() == reflect.Pointer && injTyp.Elem().Kind() == reflect.Struct {
					// fmt.Println("***", rootTyp, "***", "CreateDependency: create injTyp", injKey, injTyp)
					value, err := CreateDependency(injTyp.Elem(), reference)
					if err != nil {
						return nil, err
					}
					(*reference)[injKey] = value
				} else {
					// any other value, pass to reference
					fmt.Println("***", rootTyp, "***", "CreateDependency: pass injValue", injKey, injTyp)
					(*reference)[injKey] = injValue
				}
			} else {
				// any other value, pass to reference
				// fmt.Println("***", rootTyp, "***", "CreateDependency: pass injValue", injKey, injTyp)
				(*reference)[injKey] = injValue
			}
		}
	}

	err = InitialDependency(rootTyp, inst, reference)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

//func inspectFuncType(f interface{}) {
//	funcType := reflect.TypeOf(f)
//	if funcType.Kind() == reflect.Func {
//		Printfln("Function parameters: %v", funcType.NumIn())
//		for i := 0; i < funcType.NumIn(); i++ {
//			paramType := funcType.In(i)
//			if i < funcType.NumIn()-1 {
//				Printfln("Parameter #%v, Type: %v", i, paramType)
//			} else {
//				Printfln("Parameter #%v, Type: %v, Variadic: %v", i, paramType,
//					funcType.IsVariadic())
//			}
//		}
//		Printfln("Function results: %v", funcType.NumOut())
//		for i := 0; i < funcType.NumOut(); i++ {
//			resultType := funcType.Out(i)
//			Printfln("Result #%v, Type: %v", i, resultType)
//		}
//	}
//}
//
//func invokeFunction(f interface{}, params ...interface{}) {
//	paramVals := []reflect.Value{}
//	for _, p := range params {
//		paramVals = append(paramVals, reflect.ValueOf(p))
//	}
//	funcVal := reflect.ValueOf(f)
//	if funcVal.Kind() == reflect.Func {
//		results := funcVal.Call(paramVals)
//		for i, r := range results {
//			Printfln("Result #%v: %v", i, r)
//		}
//	}
//}
