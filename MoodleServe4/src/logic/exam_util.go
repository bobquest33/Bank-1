package logic

import (
	"math"
)

func Analyze(dig int) []int {
	res := []int{}
	if dig < 1 {
		res = append(res, 16, 8, 4, 2, 1, 0)
		return res
	}
	var digit int = dig
	m := math.Log2(float64(dig))
	max := int(m)
	for i := max; i > -1; i -- {
		d := math.Pow(2, float64(i))
		if digit >= int(d) {
			res = append(res, int(d))
			digit -= int(d)
		}
	}
	return res
}