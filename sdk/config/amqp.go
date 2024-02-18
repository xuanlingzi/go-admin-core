package config

type Amqp struct {
	Rocket *Rocket `json:"rocket,omitempty" yaml:"rocket"`
	Rabbit *Rabbit `json:"rabbit,omitempty" yaml:"rabbit"`
}

var AmqpConfig = new(Amqp)
