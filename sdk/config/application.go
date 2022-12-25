package config

type Application struct {
	Host          string `json:"host,omitempty" yaml:"host"`
	Port          int64  `json:"port,omitempty" yaml:"port"`
	Name          string `json:"name,omitempty" yaml:"name"`
	AesSecret     string `json:"aes_secret,omitempty" yaml:"aes_secret"`
	Mode          string `json:"mode,omitempty" yaml:"mode"`
	ReadTimeout   int    `json:"read_timeout,omitempty" yaml:"read_timeout"`
	WriterTimeout int    `json:"writer_timeout,omitempty" yaml:"writer_timeout"`
	EnableDP      bool   `json:"enabled_data_permission,omitempty" yaml:"enabled_data_permission"`
	DemoMessage   string `json:"demo_message,omitempty" yaml:"demo_message"`
}

var ApplicationConfig = new(Application)
