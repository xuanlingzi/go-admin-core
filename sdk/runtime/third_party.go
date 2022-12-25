package runtime

import (
	"errors"
	"github.com/xuanlingzi/go-admin-core/third_party"
)

type ThirdParty struct {
	prefix     string
	thirdParty third_party.AdapterThirdParty
}

// String string输出
func (e *ThirdParty) String() string {
	if e.thirdParty == nil {
		return ""
	}
	return e.thirdParty.String()
}

func (e *ThirdParty) GetAccessToken() (string, int, error) {
	if e.thirdParty == nil {
		return "", 0, errors.New("third party not initialized")
	}
	return e.thirdParty.GetAccessToken()
}

func (e *ThirdParty) GetJSApiTicket(accessToken string) (string, int, error) {
	if e.thirdParty == nil {
		return "", 0, errors.New("third party not initialized")
	}
	return e.thirdParty.GetJSApiTicket(accessToken)
}

func (e *ThirdParty) GetConnectUrl(state, scope string, popUp bool) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetConnectUrl(state, scope, popUp)
}

func (e *ThirdParty) GetUserAccessToken(code, state string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetUserAccessToken(code, state)
}

func (e *ThirdParty) RefreshUserToken(refreshToken string, appId string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.RefreshUserToken(refreshToken, appId)
}

func (e *ThirdParty) GetUserInfo(accessToken, openId string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	return e.thirdParty.GetUserInfo(accessToken, openId)
}

// SendTemplateMessage 发送模板消息
func (e *ThirdParty) SendTemplateMessage(accessToken, openId, templateId, url string, data []byte) (string, error) {
	if e.thirdParty == nil {
		return "", nil
	}
	return e.thirdParty.SendTemplateMessage(accessToken, openId, templateId, url, data)
}
