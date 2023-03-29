package utils

func AbsInt(x int) int {
	if x < 0 {
		return (x * -1)
	}
	return x
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
