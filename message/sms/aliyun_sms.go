package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/tidwall/sjson"
)

type AliyunSMS struct {
	client    *dysmsapi20170525.Client
	signature string
}

func NewAliyunSms(client *dysmsapi20170525.Client, accessId, accessSecret, region, signature string) *AliyunSMS {
	if client == nil {
		var err error

		config := &openapi.Config{
			// 必填，您的 AccessKey ID
			AccessKeyId: tea.String(accessId),
			// 必填，您的 AccessKey Secret
			AccessKeySecret: tea.String(accessSecret),
			// 访问的域名
			Endpoint: tea.String(region),
		}

		/* 实例化要请求产品(以sms为例)的client对象
		 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，或者引用预设的常量 */
		client = &dysmsapi20170525.Client{}
		client, err = dysmsapi20170525.NewClient(config)
		if err != nil {
			panic("Failed to setup Aliyun SMS client.")
		}
	}
	c := &AliyunSMS{
		client:    client,
		signature: signature,
	}
	return c
}

func (*AliyunSMS) Setup() error {
	return nil
}

func (*AliyunSMS) String() string {
	return "aliyun_sms"
}

func (m *AliyunSMS) Send(phones []string, templateId string, params map[string]string) error {

	if len(phones) == 0 {
		return nil
	}

	values := ""
	for k, v := range params {
		values, _ = sjson.Set(values, k, v)
	}

	for _, phone := range phones {

		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  tea.String(phone),
			SignName:      tea.String(m.signature),
			TemplateCode:  tea.String(templateId),
			TemplateParam: tea.String(values),
		}
		_, err := m.client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
		if err != nil {
			continue
		}
	}

	return nil
}

func (m *AliyunSMS) Close() {
	if m.client != nil {
		m.Close()
		m.client = nil
	}
}

// GetClient 暴露原生client
func (m *AliyunSMS) GetClient() interface{} {
	return m.client
}
