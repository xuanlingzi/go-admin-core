package lbs

import "net/http"

type AdapterLocationBasedService interface {
	String() string
	GetAddress(latitude, longitude, radius float64) (string, error)
	GetClient() *http.Client
}
