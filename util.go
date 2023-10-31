package parser

func min[T int | uint](a T, b T) T {
	if a <= b {
		return a
	}
	return b
}

func max[T int | uint](a T, b T) T {
	if a >= b {
		return a
	}
	return b
}
