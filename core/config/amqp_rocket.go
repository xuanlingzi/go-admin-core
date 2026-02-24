package config

type Rocket struct {
	Addr       string `json:"addr" yaml:"addr"`
	AccessKey  string `json:"access_key" yaml:"access_key"`
	SecretKey  string `json:"secret_key" yaml:"secret_key"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	InstanceId string `json:"instance_id" yaml:"instance_id"`
	ProducerId string `json:"producer_id" yaml:"producer_id"`
	ConsumerId string `json:"consumer_id" yaml:"consumer_id"`
}
