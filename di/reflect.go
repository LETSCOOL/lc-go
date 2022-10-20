package di

import "reflect"

func GetFuncPramDependency(f interface{}) {
	funcTyp := reflect.TypeOf(f)
	if funcTyp.Kind() != reflect.Func {

	}

}
