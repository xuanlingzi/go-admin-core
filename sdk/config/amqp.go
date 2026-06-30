package config

type Amqp struct {
	Rabbit *Rabbit `json:"rabbit,omitempty" yaml:"rabbit"`
}

var AmqpConfig = new(Amqp)
