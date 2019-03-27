package audioop

func MinInt32(x, y int32) int32 {
	if x < y {
		return x
	} else {
		return y
	}
}

func MaxInt32(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}

func AbsInt32(x int32) int32 {
	if x >= 0 {
		return x
	} else {
		return 0 - x
	}
}

func SumInt32(i ...int32) int {
	var sum int
	for _, x := range i {
		sum += int(x)
	}
	return sum
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
