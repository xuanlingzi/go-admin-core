package config

import "strings"

// WeChatPayMerchant 一个微信支付主体。
//
// 商户号是证书与密钥的边界：merchant_id / api_v3_key / serial_no / 私钥 / 平台证书
// 全都是商户级的，不同商户之间不可复用。同一个商户号下可以关联多个 appId。
//
// Key 不从 YAML 读，由所在字段名决定（rktv / activity）。
//
// 验签方式二选一，由 PublicKeyId / PublicKeyPath 是否配置决定：
//   - 都留空：平台证书模式，下载并缓存微信支付平台证书（老模式）
//   - 都配上：微信支付公钥模式，用固定的公钥验签，不再下载平台证书
//
// 新入驻的商户号只发公钥、不再发平台证书，所以两种模式必须共存。
type WeChatPayMerchant struct {
	Key            string   `yaml:"-" json:"-"`
	MerchantId     string   `yaml:"merchant_id" json:"merchant_id"`
	AppIds         []string `yaml:"app_ids" json:"app_ids"`
	PrivateKeyPath string   `yaml:"private_key_path" json:"private_key_path"`
	SerialNumber   string   `yaml:"serial_no" json:"serial_no"`
	ApiV3Key       string   `yaml:"api_v3_key" json:"api_v3_key"`
	CallbackAddr   string   `yaml:"callback_addr" json:"callback_addr"`
	WeChatCertPath string   `yaml:"wechat_cert_path" json:"wechat_cert_path"`
	// PublicKeyId 微信支付公钥 ID（形如 PUB_KEY_ID_0000...），商户平台下载公钥时一并给出
	PublicKeyId string `yaml:"public_key_id" json:"public_key_id"`
	// PublicKeyPath 微信支付公钥文件路径（pub_key.pem）
	PublicKeyPath string `yaml:"public_key_path" json:"public_key_path"`
}

// UseWeChatPayPublicKey 是否启用微信支付公钥模式。
// 必须两项都配齐才算数——只配一半是配置错误，由调用方拦下，
// 不能在这里静默退回平台证书模式，否则回调验签会以「签名错误」的形式远程失败。
func (m *WeChatPayMerchant) UseWeChatPayPublicKey() bool {
	// 与 utils.StringIsEmpty 对齐：YAML 里写成 `public_key_id: " "` 也算没配，
	// 否则 assertPublicKeyConfig 判空、这里判非空，两处会打架
	return m != nil &&
		strings.TrimSpace(m.PublicKeyId) != "" &&
		strings.TrimSpace(m.PublicKeyPath) != ""
}

// WeChatPayOption 微信支付配置。
//
// 目前只有 rktv（B 端）与 activity（C 端）两个支付主体，直接按名字平铺，
// 不做成列表——列表只会诱使人往里堆。将来若要让每个商户配自己的支付，
// 改用数据库管理，而不是继续扩这份 YAML。
type WeChatPayOption struct {
	Rktv     *WeChatPayMerchant `yaml:"rktv" json:"rktv"`
	Activity *WeChatPayMerchant `yaml:"activity" json:"activity"`

	// 以下为改造前的单商户扁平配置，rktv/activity 都没配时按单商户回退，
	// 保证旧配置文件与未改造的服务继续可用。
	MerchantId     string   `yaml:"merchant_id" json:"merchant_id"`
	AppId          string   `yaml:"app_id" json:"app_id"`
	AppIds         []string `yaml:"app_ids" json:"app_ids"`
	PrivateKeyPath string   `yaml:"private_key_path" json:"private_key_path"`
	SerialNumber   string   `yaml:"serial_no" json:"serial_no"`
	ApiV3Key       string   `yaml:"api_v3_key" json:"api_v3_key"`
	CallbackAddr   string   `yaml:"callback_addr" json:"callback_addr"`
	// CallbackAddrs 按 appId 覆盖回调地址（单商户多 appId 的旧写法）
	CallbackAddrs  map[string]string `yaml:"callback_addrs" json:"callback_addrs"`
	WeChatCertPath string            `yaml:"wechat_cert_path" json:"wechat_cert_path"`
	PublicKeyId    string            `yaml:"public_key_id" json:"public_key_id"`
	PublicKeyPath  string            `yaml:"public_key_path" json:"public_key_path"`
}
