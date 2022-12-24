package lbs

type AdapterLocationBasedService interface {
	String() string
	GetAddress(latitude, longitude, radius float64) (string, error)
}
