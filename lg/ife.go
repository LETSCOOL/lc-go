package lg

// IfeS 用來替代 golang 不支援 "?:" 的語法。
//
// Deprecated: Use lc-go/lg/Ife instead.
func IfeS(condition bool, ifValue string, elseValue string) string {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

// IfeB 用來替代 golang 不支援 "?:" 的語法。
//
// Deprecated: Use lc-go/lg/Ife instead.
func IfeB(condition bool, ifValue bool, elseValue bool) bool {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

// IfeI 用來替代 golang 不支援 "?:" 的語法。
//
// Deprecated: Use lc-go/lg/Ife instead.
func IfeI(condition bool, ifValue int, elseValue int) int {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

// Ife 用來替代 golang 不支援 "?:" 的語法。
func Ife[T string | bool | int](condition bool, ifValue T, elseValue T) T {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

func Map[IN string | bool | int | any, OUT string | bool | int | any](list []IN, f func(IN) OUT) []OUT {
	out := make([]OUT, 0, len(list))
	for _, i := range list {
		o := f(i)
		out = append(out, o)
	}
	return out
}
