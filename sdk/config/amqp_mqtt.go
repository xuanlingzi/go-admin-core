package config

// MqttConfig MQTT配置
type Mqtt struct {
	Broker              string `json:"broker" yaml:"broker"`                                 // MQTT服务器地址
	ClientId            string `json:"client_id" yaml:"client_id"`                           // 客户端ID
	Username            string `json:"username" yaml:"username"`                             // 用户名
	Password            string `json:"password" yaml:"password"`                             // 密码
	ConnectTimeout      int    `json:"connect_timeout" yaml:"connect_timeout"`               // 连接超时(秒)
	Keepalive           int    `json:"keepalive" yaml:"keepalive"`                           // 心跳间隔(秒)
	Qos                 byte   `json:"qos" yaml:"qos"`                                       // 服务质量等级
	DeviceServerTopic   string `json:"device_server_topic" yaml:"device_server_topic"`       // 设备服务端订阅主题
	ClientTopicPrefix   string `json:"client_topic_prefix" yaml:"client_topic_prefix"`       // 客户端订阅主题前缀
	CSServerTopic       string `json:"cs_server_topic" yaml:"cs_server_topic"`               // 客服服务端订阅主题
	CSClientTopicPrefix string `json:"cs_client_topic_prefix" yaml:"cs_client_topic_prefix"` // 客服订阅主题前缀
	CSId                string `json:"cs_id" yaml:"cs_id"`                                   // 客服ID
}
