package config

type SmtpConnectOptions struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
}
