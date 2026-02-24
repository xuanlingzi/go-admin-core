package utils

import "math"

const (
	a = 6378137.0
	f = 1 / 298.257223563
	b = a * (1 - f)
)

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// GetDistance 使用文森特公式（Vincenty's formulae），它考虑了地球的椭球形状，因此可以提供更精确的结果。但是，这种方法的计算复杂度更高，可能需要更多的计算资源
func GetDistance(lat1, lon1, lat2, lon2 float64) float64 {
	L := toRad(lon2 - lon1)
	U1 := math.Atan((1 - f) * math.Tan(toRad(lat1)))
	U2 := math.Atan((1 - f) * math.Tan(toRad(lat2)))
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	lambda := L
	lambdaP := 2 * math.Pi
	iterLimit := 100
	var cosSqAlpha, sigma, sinSigma, cos2SigmaM, cosSigma, sinLambda, cosLambda float64

	for math.Abs(lambda-lambdaP) > 1e-12 && iterLimit > 0 {
		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)
		sinSigma = math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) + (cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda))
		if sinSigma == 0 {
			return 0
		}
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha := cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		if cosSqAlpha != 0 {
			cos2SigmaM = cosSigma - 2*sinU1*sinU2/cosSqAlpha
		} else {
			cos2SigmaM = 0
		}
		C := f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))
		lambdaP = lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
		iterLimit--
	}

	if iterLimit == 0 {
		return 0
	}

	uSq := cosSqAlpha * (a*a - b*b) / (b * b)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))
	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))

	return b * A * (sigma - deltaSigma)
}
