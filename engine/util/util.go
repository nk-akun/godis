package util

// Max return max(x,y)
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Min return min(x,y)
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// NearLargeUnsignedBinary return the first number larger than given num
func NearLargeUnsignedBinary(num uint32) uint32 {
	var i uint32
	for i = 2; i < num; i <<= 1 {
	}
	return i
}
