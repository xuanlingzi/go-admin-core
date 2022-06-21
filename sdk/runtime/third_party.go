package runtime

import (
	"errors"
	"github.com/xuanlingzi/go-admin-core/message"
)

type ThirdParty struct {
	prefix     string
	thirdParty message.AdapterThirdParty
}

// String string输出
func (e *ThirdParty) String() string {
	if e.thirdParty == nil {
		return ""
	}
	return e.thirdParty.String()
}

func (e *ThirdParty) GetConnectUrl(state, scope, redirectUrl string, popUp bool) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetConnectUrl(state, scope, redirectUrl, popUp)
}

func (e *ThirdParty) GetUserAccessToken(code, state string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetUserAccessToken(code, state)
}

func (e *ThirdParty) GetUserInfo(accessToken, openId string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetUserInfo(accessToken, openId)
}

// GetAccessToken 获取access token
func (e *ThirdParty) GetAccessToken(force bool) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetAccessToken(force)
}

// SendTemplateMessage 发送模板消息
func (e *ThirdParty) SendTemplateMessage(openId, templateId, url string, data []byte) (string, error) {
	if e.thirdParty == nil {
		return "", nil
	}
	return e.thirdParty.SendTemplateMessage(openId, templateId, url, data)
}
