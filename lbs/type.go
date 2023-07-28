package lbs

type AdapterLocationBasedService interface {
	String() string
	Close()
	GetAddress(latitude, longitude, radius float64) (string, error)
	GetClient() interface{}
}
