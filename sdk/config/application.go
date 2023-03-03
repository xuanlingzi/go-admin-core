package config

type Application struct {
	Mode          string `json:"mode,omitempty" yaml:"mode"`
	Name          string `json:"name,omitempty" yaml:"name"`
	Host          string `json:"host,omitempty" yaml:"host"`
	Port          int    `json:"port,omitempty" yaml:"port"`
	JwtSecret     string `json:"jwt_secret,omitempty" yaml:"jwt_secret"`
	ReadTimeout   int    `json:"read_timeout,omitempty" yaml:"read_timeout"`
	WriterTimeout int    `json:"writer_timeout,omitempty" yaml:"writer_timeout"`
	EnableDP      bool   `json:"enabled_data_permission,omitempty" yaml:"enabled_data_permission"`
	DemoMessage   string `json:"demo_message,omitempty" yaml:"demo_message"`
}

var ApplicationConfig = new(Application)
