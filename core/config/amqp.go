package config

type Amqp struct {
	Rabbit *Rabbit `json:"rabbit,omitempty" yaml:"rabbit"`
	Mqtt   *Mqtt   `json:"mqtt,omitempty" yaml:"mqtt"`
}

var AmqpConfig = new(Amqp)
