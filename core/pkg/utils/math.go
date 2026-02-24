package utils

import (
	"math"
)

func Gcd(x, y int) int {
	if y == 0 {
		return x
	}
	return Gcd(y, x%y)
}

func Lcm(x, y int) int {
	return x * y / Gcd(x, y)
}

func GetLcm(times []int) int {
	g := times[0]
	for _, t := range times {
		g = Lcm(g, t)
	}
	return g
}

func GetGcd(times []int) int {
	g := times[0]
	for _, t := range times {
		g = Gcd(g, t)
	}
	return g
}

func Round(x float64) int {
	return int(math.Floor(x + 0/5))
}

func DivisibleBy2(x float64) int {
	tx := int(math.Floor(x))
	t := tx % 2
	if t == 0 {
		return tx
	} else {
		return tx + 1
	}
}
