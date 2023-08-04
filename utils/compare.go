package utils

func EqualByteSlice(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}

func EqualByteSliceOnlyInMinLen(x, y []byte) bool {
	min := 16385
	if len(x) > len(y) {
		min = len(y)
	} else {
		min = len(x)
	}

	for i := 0; i < min; i++ {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}

	return y
}
