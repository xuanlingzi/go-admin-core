package third_party

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/http"
)

var _wechatOpen *http.Client

type WeChatOpen struct {
	conn         *http.Client
	appId        string
	appSecret    string
	callbackAddr string
}

func GetWeChatOpenClient() *http.Client {
	return _wechatMp
}

func NewWeChatOpen(client *http.Client, appId, appSecret, callbackAddr string) *WeChatOpen {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{},
		}
	}
	c := &WeChatOpen{
		conn:         client,
		appId:        appId,
		appSecret:    appSecret,
		callbackAddr: callbackAddr,
	}
	return c
}

// Close 关闭连接
func (m *WeChatOpen) Close() {
	if m.conn != nil {
		m.conn.CloseIdleConnections()
		m.conn = nil
	}
}

func (*WeChatOpen) String() string {
	return "wechat_open"
}

func (m *WeChatOpen) GetAccessToken() (string, int, error) {

	var err error
	var accessToken string
	/*
		https://api.weixin.qq.com/cgi-bin/token?
		grant_type=client_credential
		&appid=APPID
		&secret=APPSECRET
	*/
	url := fmt.Sprintf("%v?appid=%v&secret=%v&grant_type=client_credential",
		WeChatAccessTokenAddr,
		m.appId,
		m.appSecret,
	)
	body, err := httpGet(url)
	if err != nil {
		return "", 0, err
	}

	accessToken = gjson.Get(body, "access_token").String()
	expiresIn := cast.ToInt(gjson.Get(body, "expires_in").Int())

	return accessToken, expiresIn, nil
}

func (m *WeChatOpen) GetJSApiTicket(accessToken string) (string, int, error) {

	var err error
	var ticket string
	/*
		https://api.weixin.qq.com/cgi-bin/ticket/getticket?
		access_token=ACCESS_TOKEN
		&type=jsapi
	*/
	url := fmt.Sprintf("%v?type=%v&access_token=%v",
		WeChatJSApiTicketAddr,
		"jsapi",
		accessToken,
	)
	body, err := httpGet(url)
	if err != nil {
		return "", 0, err
	}

	ticket = gjson.Get(body, "ticket").String()
	expiresIn := cast.ToInt(gjson.Get(body, "expires_in").Int())

	return ticket, expiresIn, nil
}

func (m *WeChatOpen) GetConnectUrl(state, scope string, popUp bool) (string, error) {
	/*
		url?
		appid=APPID
		&redirect_uri=REDIRECT_URI
		&response_type=code
		&scope=SCOPE
		&state=STATE
		&forcePopup=FORCE_POPUP
		#wechat_redirect
	*/
	url := fmt.Sprintf("%v?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v&forcePopup=%v#wechat_redirect",
		WeChatQRConnectAddr,
		m.appId,
		m.callbackAddr,
		scope,
		state,
		popUp,
	)
	return url, nil
}

func (m *WeChatOpen) GetUserAccessToken(code, scope string) (string, error) {
	/*
		https://api.weixin.qq.com/sns/oauth2/access_token?
		appid=APPID
		&secret=SECRET
		&code=CODE
		&grant_type=authorization_code
	*/
	url := fmt.Sprintf("%v?appid=%v&secret=%v&code=%v&grant_type=authorization_code",
		WeChatUserAccessTokenAddr,
		m.appId,
		m.appSecret,
		code,
	)

	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatOpen) RefreshUserToken(refreshToken string, appId string) (string, error) {
	/*
		https://api.weixin.qq.com/sns/oauth2/refresh_token?
		appid=APPID
		&grant_type=refresh_token
		&refresh_token=REFRESH_TOKEN
	*/
	url := fmt.Sprintf("%v?appid=%v&refresh_token=%v&grant_type=refresh_token",
		WeChatRefreshUserTokenAddr,
		appId,
		refreshToken,
	)

	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatOpen) GetUserInfo(accessToken, openId string) (string, error) {
	/*
		https://api.weixin.qq.com/sns/userinfo?
		access_token=ACCESS_TOKEN
		&openid=OPENID
		&lang=zh_CN
	*/
	url := fmt.Sprintf("%v?access_token=%v&openid=%v&lang=zh_CN",
		WeChatUserInfoAddr,
		accessToken,
		openId,
	)

	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatOpen) SendTemplateMessage(accessToken, openId, templateId, redirectUrl string, data []byte) (string, error) {
	return "", errors.New("开放平台无法发送消息")
}

// GetClient 暴露原生client
func (m *WeChatOpen) GetClient() *http.Client {
	return m.conn
}
