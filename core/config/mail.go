package config

type Mail struct {
	Smtp *SmtpConnectOptions `json:"smtp" yaml:"smtp"`
}

var MailConfig = new(Mail)
