package lbs

type AdapterLocationBasedService interface {
	String() string
	Close()
	GetAddress(latitude, longitude, radius float64) (map[string]string, error)
	GetCoordinate(keyword string) (longitude float32, latitude float32, address string, err error)
	GetPosition(imei, network, ac, ci, snr string, result *map[string]interface{}) error
	GetClient() interface{}
}
