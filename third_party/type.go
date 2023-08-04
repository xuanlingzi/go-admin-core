package third_party

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/xuanlingzi/go-admin-core/tools/utils"
)

// Wechat
const (
	WECHAT_ACCESS_TOKEN_KEY      = "WECHAT_ACCESS_TOKEN:%s"
	WECHAT_JSAPI_TICKET_KEY      = "WECHAT_JSAPI_TICKET:%s"
	WECHAT_STATE_KEY             = "WECHAT_STATE:%v"
	WECHAT_PLATFORM_KEY          = "WECHAT_PLATFORM:%v"
	WECHAT_REDIRECT_KEY          = "WECHAT_REDIRECT:%v"
	WECHAT_REDIRECT_EXP_KEY      = "WECHAT_REDIRECT_EXP:%v"
	WECHAT_USER_ACCESS_TOKEN_KEY = "WECHAT_USER_ACCESS_TOKEN:%v"
)

var (
	WeChatQRLogin  = "snsapi_login"
	WeChatUserInfo = "snsapi_userinfo"
	WeChatBase     = "snsapi_base"
)

var (
	WeChatAccessTokenAddr = "https://api.weixin.qq.com/cgi-bin/token"
	WeChatJSApiTicketAddr = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"

	WeChatQRConnectAddr        = "https://open.weixin.qq.com/connect/qrconnect"
	WeChatAppConnectAddr       = "https://open.weixin.qq.com/connect/oauth2/authorize"
	WeChatUserAccessTokenAddr  = "https://api.weixin.qq.com/sns/oauth2/access_token"
	WeChatRefreshUserTokenAddr = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	WeChatUserInfoAddr         = "https://api.weixin.qq.com/sns/userinfo"

	WeChatTemplateMessageAddr = "https://api.weixin.qq.com/cgi-bin/message/template/send"
)

type AdapterThirdParty interface {
	String() string
	GetConnectUrl(state, scope string, popUp bool) (string, error)
	GetAccessToken() (string, int, error)
	GetJSApiTicket(accessToken string) (string, int, error)
	GetUserAccessToken(code, state string) (string, error)
	RefreshUserToken(refreshToken string, appId string) (string, error)
	GetUserInfo(userAccessToken, openId string) (string, error)
	SendTemplateMessage(accessToken, openId, templateId, url string, data []byte) (string, error)
}

func httpGet(url string) (string, error) {
	body, err := utils.HttpGet(url)
	if err != nil {
		return "", err
	}

	if gjson.GetBytes(body, "errcode").Exists() && gjson.GetBytes(body, "errcode").Int() != 0 {
		return "", errors.New(gjson.GetBytes(body, "errmsg").String())
	}

	return string(body), nil
}

func httpPost(url, content string) (string, error) {
	body, err := utils.HttpPost(url, content)
	if err != nil {
		return "", err
	}

	if gjson.GetBytes(body, "errcode").Exists() && gjson.GetBytes(body, "errcode").Int() != 0 {
		return "", errors.New(gjson.GetBytes(body, "errmsg").String())
	}

	return string(body), nil
}
