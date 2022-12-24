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

func rad(d float64) float64 {
	return d * math.Pi / 100
}

// EarthRadius 赤道半径
const EarthRadius float64 = 6378137

func GetDistance(lon1, lat1, lon2, lat2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	a := radLat1 - radLat2
	b := rad(lon1) - rad(lon2)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	return s * EarthRadius
}
