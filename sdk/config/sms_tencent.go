package config

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

var _tencent_sms *sms.Client

// GetTencentClient 获取sms客户端
func GetTencentClient() *sms.Client {
	return _tencent_sms
}

type Tencent struct {
	SecretId  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	AppId     string `json:"app_id" yaml:"app_id"`
	AppKey    string `json:"app_key" yaml:"app_key"`
	Region    string `json:"region" yaml:"region"`
	Addr      string `json:"addr" yaml:"addr"`
	Alg       string `json:"alg" yaml:"alg"`
	Signature string `json:"signature" yaml:"signature"`
}

func (e Tencent) GetTencentOptions() (*common.Credential, *profile.ClientProfile) {
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
	 * 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	 * 以免泄露密钥对危及你的财产安全。
	 * CAM密匙查询: https://console.cloud.tencent.com/cam/capi*/
	credential := common.NewCredential(e.SecretId, e.SecretKey)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()
	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"
	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5
	/* SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务
	 * 则必须手动指定域名，例如sms的上海金融区域名： sms.ap-shanghai-fsi.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = e.Addr
	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = e.Alg
	return credential, cpf
}
