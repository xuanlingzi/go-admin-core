package runtime

import (
	"errors"
	"github.com/xuanlingzi/go-admin-core/sdk"
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

func (e *ThirdParty) GetAccessToken() (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	accessToken, err := sdk.Runtime.GetCacheAdapter().Get(third_party.WECHAT_ACCESS_TOKEN)
	if err != nil {
		return accessToken, nil
	}

	accessToken, expireAt, err := e.thirdParty.GetAccessToken()
	if err = sdk.Runtime.GetCacheAdapter().Set(third_party.WECHAT_ACCESS_TOKEN, accessToken, expireAt); err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func (e *ThirdParty) GetJSApiTicket() (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}
	ticket, err := sdk.Runtime.GetCacheAdapter().Get(third_party.WECHAT_JSAPI_TICKET)
	if err != nil {
		return ticket, nil
	}

	accessToken, err := e.GetAccessToken()
	if err != nil {
		return ticket, err
	}
	ticket, expireAt, err := e.thirdParty.GetJSApiTicket(accessToken)
	if err = sdk.Runtime.GetCacheAdapter().Set(third_party.WECHAT_JSAPI_TICKET, ticket, expireAt); err != nil {
		return ticket, err
	}

	return ticket, nil
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
func (e *ThirdParty) SendTemplateMessage(openId, templateId, url string, data []byte) (string, error) {
	if e.thirdParty == nil {
		return "", nil
	}
	return e.thirdParty.SendTemplateMessage(openId, templateId, url, data)
}
