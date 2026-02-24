package config

type Database struct {
	Driver          string             `json:"driver,omitempty" yaml:"driver"`
	Source          string             `json:"source,omitempty" yaml:"source"`
	ConnMaxIdleTime int                `json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time"`
	ConnMaxLifetime int                `json:"conn_max_lifetime,omitempty" yaml:"conn_max_lifetime"`
	MaxIdleConns    int                `json:"max_idle_conns,omitempty" yaml:"max_idle_conns"`
	MaxOpenConns    int                `json:"max_open_conns,omitempty" yaml:"max_open_conns"`
	Registers       []DBResolverConfig `json:"registers,omitempty" yaml:"registers"`
}

type DBResolverConfig struct {
	Sources  []string
	Replicas []string
	Policy   string
	Tables   []string
}

var (
	DatabaseConfig  = new(Database)
	DatabasesConfig = make(map[string]*Database)
)
