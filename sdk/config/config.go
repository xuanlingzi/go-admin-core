package config

import (
	"fmt"
	"log"

	"github.com/xuanlingzi/go-admin-core/config"
	"github.com/xuanlingzi/go-admin-core/config/source"
)

var (
	ExtendConfig interface{}
	_cfg         *Settings
)

// Settings 兼容原先的配置结构
type Settings struct {
	Settings  Config `yaml:"settings"`
	callbacks []func()
}

func (e *Settings) runCallback() {
	for i := range e.callbacks {
		e.callbacks[i]()
	}
}

func (e *Settings) OnChange() {
	e.init()
	log.Println("config change and reload")
}

func (e *Settings) Init() {
	e.init()
	log.Println("config init")
}

func (e *Settings) init() {
	e.Settings.Logger.Setup()
	e.Settings.multiDatabase()
	e.runCallback()
}

// Config 配置集合
type Config struct {
	Application *Application          `json:"application" yaml:"application"`
	Secret      *Secret               `json:"secret" yaml:"secret"`
	Ssl         *Ssl                  `json:"ssl" yaml:"ssl"`
	Logger      *Logger               `json:"logger" yaml:"logger"`
	Jwt         *Jwt                  `json:"jwt" yaml:"jwt"`
	Database    *Database             `json:"database" yaml:"database"`
	Databases   *map[string]*Database `json:"databases" yaml:"databases"`
	Gen         *Gen                  `json:"gen" yaml:"gen"`
	Cache       *Cache                `json:"cache" yaml:"cache"`
	Queue       *Queue                `json:"queue" yaml:"queue"`
	Locker      *Locker               `json:"locker" yaml:"locker"`
	Lbs         *Lbs                  `json:"lbs,omitempty" yaml:"lbs"`
	FileStore   *FileStore            `json:"file_store,omitempty" yaml:"file_store"`
	Sms         *Sms                  `json:"sms,omitempty" yaml:"sms"`
	Mail        *Mail                 `json:"mail" yaml:"mail"`
	Amqp        *Amqp                 `json:"mq,omitempty" yaml:"mq"`
	WeChat      *WeChat               `json:"wechat,omitempty" yaml:"wechat"`
	Payment     *Payment              `json:"payment,omitempty" yaml:"payment"`
	Extend      interface{}           `json:"extend" yaml:"extend"`
}

// 多db改造
func (e *Config) multiDatabase() {
	if len(*e.Databases) == 0 {
		*e.Databases = map[string]*Database{
			"*": e.Database,
		}

	}
}

// Setup 载入配置文件
func Setup(s source.Source,
	fs ...func()) {
	_cfg = &Settings{
		Settings: Config{
			Application: ApplicationConfig,
			Secret:      SecretConfig,
			Ssl:         SslConfig,
			Logger:      LoggerConfig,
			Jwt:         JwtConfig,
			Database:    DatabaseConfig,
			Databases:   &DatabasesConfig,
			Gen:         GenConfig,
			Cache:       CacheConfig,
			Queue:       QueueConfig,
			Locker:      LockerConfig,
			Lbs:         LbsConfig,
			FileStore:   FileStoreConfig,
			Sms:         SmsConfig,
			Mail:        MailConfig,
			Amqp:        AmqpConfig,
			WeChat:      WeChatConfig,
			Payment:     PaymentConfig,
			Extend:      ExtendConfig,
		},
		callbacks: fs,
	}
	var err error
	config.DefaultConfig, err = config.NewConfig(
		config.WithSource(s),
		config.WithEntity(_cfg),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("New config object fail: %s", err.Error()))
	}
	_cfg.Init()
}
