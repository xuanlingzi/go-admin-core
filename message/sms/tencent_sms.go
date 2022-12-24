package sms

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentCloudSMS struct {
	client    *sms.Client
	appId     string
	signature string
}

func NewTencentSms(client *sms.Client, credential *common.Credential, cpf *profile.ClientProfile, appId, signature, region string) *TencentCloudSMS {
	if client == nil {
		var err error

		/* 实例化要请求产品(以sms为例)的client对象
		 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，或者引用预设的常量 */
		client, err = sms.NewClient(credential, region, cpf)
		if err != nil {
			panic("Failed to setup TencentCloudSMS client.")
		}
	}
	c := &TencentCloudSMS{
		client:    client,
		appId:     appId,
		signature: signature,
	}
	return c
}

func (*TencentCloudSMS) Setup() error {
	return nil
}

func (*TencentCloudSMS) String() string {
	return "tencent_sms"
}

func (m *TencentCloudSMS) Send(phones []string, templateId string, params map[string]string) error {

	if len(phones) == 0 {
		return nil
	}

	var values []string
	for _, v := range params {
		values = append(values, v)
	}

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	* 你可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/sms/smslist
	 * sms helper: https://cloud.tencent.com/document/product/382/3773 */

	/* 短信应用ID: 短信SdkAppid在 [短信控制台] 添加应用后生成的实际SdkAppid，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr(m.appId)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
	request.SignName = common.StringPtr(m.signature)
	/* 国际/港澳台短信 senderid: 国内短信填空，默认未开通，如需开通请联系 [sms helper] */
	//request.SenderId = common.StringPtr("xxx")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	//request.SessionContext = common.StringPtr("xxx")
	/* 短信码号扩展号: 默认未开通，如需开通请联系 [sms helper] */
	//request.ExtendCode = common.StringPtr("0")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(values)
	/* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
	request.TemplateId = common.StringPtr(templateId)
	/* 下发手机号码，采用 e.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	var formatPhones []string
	for _, phone := range phones {
		formatPhones = append(formatPhones, fmt.Sprintf("+86%v", phone))
	}
	request.PhoneNumberSet = common.StringPtrs(formatPhones)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	_, err := m.client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}

	//b, _ := json.Marshal(response.Response)
	//// 打印返回的json字符串
	//logger.Infof("%s", b)

	return nil
}

func (m *TencentCloudSMS) Close() {
	if m.client != nil {
		m.Close()
		m.client = nil
	}
}

// GetClient 暴露原生client
func (m *TencentCloudSMS) GetClient() *sms.Client {
	return m.client
}