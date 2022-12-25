package config

type Amqp struct {
	Pulsar *Pulsar `json:"pulsar,omitempty" yaml:"pulsar"`
	Rocket *Rocket `json:"rocket,omitempty" yaml:"rocket"`
}

var AmqpConfig = new(Amqp)
