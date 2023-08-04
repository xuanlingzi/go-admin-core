package runtime

import (
	"errors"
	"fmt"
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/third_party"
)

type ThirdParty struct {
	prefix     string
	thirdParty third_party.AdapterThirdParty
}

// NewThirdParty 创建对应上下文缓存
func NewThirdParty(prefix string, thirdParty third_party.AdapterThirdParty) *ThirdParty {
	return &ThirdParty{
		prefix:     prefix,
		thirdParty: thirdParty,
	}
}

// String string输出
func (e *ThirdParty) String() string {
	if e.thirdParty == nil {
		return ""
	}
	return e.thirdParty.String()
}

func (e *ThirdParty) GetAccessToken(cache storage.AdapterCache) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}

	var err error
	key := fmt.Sprintf(third_party.WECHAT_ACCESS_TOKEN_KEY, e.prefix)
	accessToken, err := cache.Get(key)
	if err != nil {
		var expireIn int
		accessToken, expireIn, err := e.thirdParty.GetAccessToken()
		if err == nil {
			cache.Set(key, accessToken, expireIn)
		}
	}
	return accessToken, err
}

func (e *ThirdParty) GetJSApiTicket(cache storage.AdapterCache) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}

	var err error
	key := fmt.Sprintf(third_party.WECHAT_JSAPI_TICKET_KEY, e.prefix)
	ticket, err := cache.Get(key)
	if err != nil {

		// 获取accessToken
		var accessToken string
		accessToken, err = e.GetAccessToken(cache)
		if err != nil {
			return "", err
		}

		var expireIn int
		ticket, expireIn, err := e.thirdParty.GetJSApiTicket(accessToken)
		if err == nil {
			cache.Set(key, ticket, expireIn)
		}
	}

	return ticket, err
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

func (e *ThirdParty) GetUserInfo(userAccessToken, openId string) (string, error) {
	if e.thirdParty == nil {
		return "", errors.New("third party not initialized")
	}

	return e.thirdParty.GetUserInfo(userAccessToken, openId)
}

// SendTemplateMessage 发送模板消息
func (e *ThirdParty) SendTemplateMessage(cache storage.AdapterCache, openId, templateId, url string, data []byte) (string, error) {
	if e.thirdParty == nil {
		return "", nil
	}

	accessToken, err := e.GetAccessToken(cache)
	if err != nil {
		return "", err
	}

	return e.thirdParty.SendTemplateMessage(accessToken, openId, templateId, url, data)
}
