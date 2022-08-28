package lg

// IfeS 用來替代 golang 不支援 "?:" 的語法。
func IfeS(condition bool, ifValue string, elseValue string) string {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

func IfeB(condition bool, ifValue bool, elseValue bool) bool {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

func IfeI(condition bool, ifValue int, elseValue int) int {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}