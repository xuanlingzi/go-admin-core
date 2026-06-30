package config

type Lbs struct {
	Amap *Amap `json:"amap"`
}

var LbsConfig = new(Lbs)
