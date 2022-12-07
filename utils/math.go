package utils

// AbsInt 返回给出数据的绝对值
// 传入数据为int
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
