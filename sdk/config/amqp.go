package config

type Amqp struct {
	Pulsar *Pulsar
	Rocket *Rocket
}

var AmqpConfig = new(Amqp)
