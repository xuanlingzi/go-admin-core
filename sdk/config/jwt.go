package config

type Jwt struct {
	Secret  string `json:"secret,omitempty" yaml:"secret"`
	Timeout int64  `json:"timeout,omitempty" yaml:"timeout"`
}

var JwtConfig = new(Jwt)
