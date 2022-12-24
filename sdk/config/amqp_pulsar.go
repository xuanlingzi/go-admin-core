package config

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"time"
)

type Pulsar struct {
	Addr      string `yaml:"addr"`
	AppId     string `yaml:"appId"`
	Namespace string `yaml:"namespace"`
	Role      string `yaml:"role"`
	SecretKey string `yaml:"secretKey"`
}

func (p *Pulsar) GetClientOptions() *pulsar.ClientOptions {
	return &pulsar.ClientOptions{
		URL:               fmt.Sprintf("http://%v", p.Addr),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
		Authentication:    pulsar.NewAuthenticationToken(p.SecretKey),
	}
}
