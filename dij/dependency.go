package dij

import (
	"fmt"
	. "github.com/letscool/lc-go/lg"
	"log"
	"reflect"
	"strings"
)

type DependencyKey = string
type DependencyReference *map[DependencyKey]any
type DependencyStack []*DependencyStackRecord

func (s DependencyStack) NumOfRecords() int {
	return len(s)
}

func (s DependencyStack) GetRecord(index int) *DependencyStackRecord {
	return s[index]
}

const (
	TagName      = "di"
	StackDeepKey = "__stack_deep__*"
	StackKey     = "__stack__*"
)

type DependencyStackRecord struct {
	Inst     any
	Deep     int
	Fullname string
}

func (r DependencyStackRecord) InstType() reflect.Type {
	return reflect.TypeOf(r.Inst)
}

func (r DependencyStackRecord) NameOfInstType() string {
	return fmt.Sprintf("%v", reflect.TypeOf(r.Inst))
}

type InjectionHandler interface {
	// DidDependencyInjection will be called after injection is completed.
	DidDependencyInjection()
}

type InitializationHandler interface {
	// DidDependencyInitialization will be called after initialization is completed. If receiver implements InjectionHandler
	// interface, InjectionHandler.DidDependencyInjection will be called first.
	DidDependencyInitialization()
}

var LogEnabled = false

func EnableLog() {
	LogEnabled = true
}

func CallDependencyInjection(initMethod reflect.Method, inst any, reference DependencyReference) error {
	methodTyp := initMethod.Type
	//log.Println(methodTyp.NumIn())
	if methodTyp.NumIn() <= 0 {
		// only receiver
		//log.Println("Only receiver", reflect.ValueOf(inst).MethodByName("InjectDependency").IsValid())
		reflect.ValueOf(inst).MethodByName("InjectDependency").Call(nil)
	} else {
		numMethodIn := methodTyp.NumIn()
		args := make([]reflect.Value, 0, numMethodIn-1)
		//log.Println("numMethodIn", numMethodIn)
		for i := 1; i < methodTyp.NumIn(); i++ {
			inTyp := methodTyp.In(i)
			instPtrValue, err := createAndInitializeInstance(inTyp, reference, true, "")
			if err != nil {
				return err
			}
			instValue := instPtrValue.Elem()
			args = append(args, instValue)
		}

		//Invoke(inst, "InjectDependency", args...)
		//log.Println("args num", len(args))
		reflect.ValueOf(inst).MethodByName("InjectDependency").Call(args)
	}

	return nil
}

func createAndInitializeInstance(insTyp reflect.Type, reference DependencyReference, forParameter bool, applyingName string) (reflect.Value, error) {
	// ================================
	// save stack
	stack := &DependencyStackRecord{
		Fullname: applyingName,
	}
	stackDeepCount := 0
	if count, existing := (*reference)[StackDeepKey]; existing {
		stackDeepCount = count.(int) + 1
		if stackDeepCount > 20 {
			return reflect.ValueOf(nil), fmt.Errorf("statck go to deep (%d)", stackDeepCount)
		}
	} else {
		stackDeepCount = 1
	}
	(*reference)[StackDeepKey] = stackDeepCount
	stack.Deep = stackDeepCount
	defer func() {
		(*reference)[StackDeepKey] = (*reference)[StackDeepKey].(int) - 1
	}()
	if stackSlice, existing := (*reference)[StackKey]; existing {
		slice := stackSlice.(DependencyStack)
		slice = append(slice, stack)
		(*reference)[StackKey] = slice
	} else {
		slice := DependencyStack{stack}
		(*reference)[StackKey] = slice
	}
	// == end of saving stack deep count
	// =================================

	if insTyp.Kind() != reflect.Struct {
		//log.Println(insTyp.Kind())
		return reflect.ValueOf(nil), fmt.Errorf("initialization function support struct parameters(%s) only", insTyp.Name())
	}

	// create new instance and save for reference
	//log.Println("In", i, insTyp, insTyp.Kind())
	instPtrValue := reflect.New(insTyp)
	instValue := instPtrValue.Elem()
	instPtrIf := instPtrValue.Interface()
	if applyingName != "" {
		if LogEnabled {
			log.Printf("*** Set %v by type %v\n", applyingName, insTyp)
		}
		(*reference)[applyingName] = instPtrIf
	}
	stack.Inst = instPtrIf

	if err := initializeInstance(insTyp, instValue, reference, forParameter); err != nil {
		return reflect.ValueOf(nil), err
	}

	if injectMethod, ok := GetMethodForReceiverType(insTyp, "InjectDependency"); ok {
		if err := CallDependencyInjection(injectMethod, instPtrIf, reference); err != nil {
			return reflect.ValueOf(nil), err
		}
	}

	return instPtrValue, nil

}

func initializeInstance(insTyp reflect.Type, instValue reflect.Value, reference DependencyReference, forParameter bool) error {
	TypeOfType := reflect.TypeOf(reflect.TypeOf(struct{}{}))

	for j := 0; j < insTyp.NumField(); j++ {
		fieldSpec := insTyp.Field(j)
		diTag, existingDiTag := fieldSpec.Tag.Lookup(TagName)
		if fieldSpec.Type.Kind() == reflect.Struct {
			if existingDiTag {
				return fmt.Errorf("embedded/extended field can't do dependency injection, (%v)", fieldSpec)
			}
			// embedded/extended field with struct kind may contain dependency injection tag, initialize it
			fmt.Printf("***** %v %d from %v \n", fieldSpec, j, insTyp)
			instValForField := instValue.Field(j)
			if err := initializeInstance(fieldSpec.Type, instValForField, reference, false); err != nil {
				return err
			}
			continue
		}
		name := parseDiTag(diTag)
		if name == "" {
			if existingDiTag {
				name = fieldSpec.Name
				if name == "" || name == "_" {
					name = FullnameOfType(fieldSpec.Type)
				}
			} else {
				// do not be referred for dependency injection
				continue
			}
		} else if name == "-" {
			// do not be referred for dependency injection
			continue
		} else if name == "^" {
			name = FullnameOfType(fieldSpec.Type)
		}
		//log.Printf("Field name: %s (%d) %v ", name, len(name), fieldSpec)
		if l := len(name); l == 0 {
			return fmt.Errorf("not support anonymous name for dependency reference")
		} else if l == 1 {
			c := name[0]
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				return fmt.Errorf("not support symbol '%c' as dependency reference", c)
			}
		}
		if LogEnabled {
			log.Printf("Field name: %s for %v", name, fieldSpec)
		}
		var refValue any
		if v, existing := (*reference)[name]; existing {
			if reflect.TypeOf(v) == TypeOfType {
				insTyp := v.(reflect.Type)
				if insTyp.Kind() == reflect.Struct {
					instPtrValForField, err := createAndInitializeInstance(insTyp, reference, false, name)
					if err != nil {
						return err
					}
					instForField := instPtrValForField.Interface()
					//(*reference)[name] = inst
					refValue = instForField
				} else {
					// TODO: any good way?
					refValue = v
				}
			} else {
				refValue = v
			}
		} else {
			if fieldSpec.Type.Kind() == reflect.Pointer && fieldSpec.Type.Elem().Kind() == reflect.Struct {
				// create and initialize instance
				underlyingType := fieldSpec.Type.Elem()
				instPtrValForField, err := createAndInitializeInstance(underlyingType, reference, false, name)
				if err != nil {
					return err
				}
				instForField := instPtrValForField.Interface()
				//(*reference)[name] = inst
				refValue = instForField
			} else {
				if forParameter {
					return fmt.Errorf("non available reference(%s) for injection", name)
				} else {
					refValue = nil
				}
			}
		}
		if refValue == nil {
			// ignored
		} else if fieldSpec.Type == reflect.TypeOf(refValue) {
			field := instValue.Field(j)
			if field.Type() != fieldSpec.Type {
				log.Fatalf("Struct instance and type have different type for field index: %d, %v != %v", j, field.Type(), fieldSpec.Type)
			}
			//field := instValue.FieldByName(fieldSpec.Name)
			firstChar := fieldSpec.Name[0]
			if firstChar == '_' {
				// don't assign
			} else if firstChar >= 'A' && firstChar <= 'Z' {
				field.Set(reflect.ValueOf(refValue))
			} else {
				SetUnexportedField(field, refValue)
			}
		} else {
			return fmt.Errorf("filed(%s)'s declaration(%v) and value(%v) are not same type", fieldSpec.Name, fieldSpec.Type, reflect.TypeOf(refValue))
		}
	}
	return nil
}

// CreateInstance create an instance of rootTyp, the rootTyp should be kind of struct.
//
// A pointer of an instance for rootTyp will be returned if success.
func CreateInstance(rootTyp reflect.Type, reference DependencyReference, instName string) (any, error) {
	if reference == nil {
		reference = &map[DependencyKey]any{}
	}

	if instName == "^" || instName == "" || instName == "_" {
		instName = FullnameOfType(reflect.PointerTo(rootTyp))
	} // else if instName == "" {
	//	log.Printf("Root instance doesn't support empty name. If you use empty name, it will not be referred in dependency-injection flow.")
	//}

	instPtrValue, err := createAndInitializeInstance(rootTyp, reference, false, instName)
	if err != nil {
		return nil, err
	}
	instPtrIf := instPtrValue.Interface()

	if deepCount := (*reference)[StackDeepKey].(int); deepCount != 0 {
		return nil, fmt.Errorf("incorrect final stack deep count(%d)", deepCount)
	}

	if stackSlice, existing := (*reference)[StackKey]; existing {
		slice := stackSlice.(DependencyStack)
		if LogEnabled {
			for _, stack := range slice {
				log.Printf("%2d. %s => %v\n", stack.Deep, Ife(stack.Fullname == "", "(NAV)", stack.Fullname), reflect.TypeOf(stack.Inst))
			}
		}
		sliceLen := len(slice)
		for j := sliceLen - 1; j >= 0; j-- {
			stack := slice[j]
			if handler, ok := stack.Inst.(InjectionHandler); ok {
				handler.DidDependencyInjection()
			}
		}
		for j := sliceLen - 1; j >= 0; j-- {
			stack := slice[j]
			if handler, ok := stack.Inst.(InitializationHandler); ok {
				handler.DidDependencyInitialization()
			}
		}
	} else {
		return nil, fmt.Errorf("missing stack")
	}

	return instPtrIf, nil
}

func GetCountOfDependencyStack(ref DependencyReference) int {
	return (*ref)[StackDeepKey].(int)
}

func GetHistoryOfDependencyStack(ref DependencyReference) (stack DependencyStack) {
	if stackSlice, existing := (*ref)[StackKey]; !existing {
		return nil
	} else {
		slice := stackSlice.(DependencyStack)
		return slice
	}
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
