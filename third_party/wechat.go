package third_party

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"net/url"
	"strings"
)

var _wechatClient = make(map[string]*http.Client)

type WeChatClient struct {
	conn         *http.Client
	appId        string
	appSecret    string
	appType      string
	callbackAddr string
}

func GetWeChatClient(appId string) *http.Client {
	return _wechatClient[appId]
}

func NewWeChatClient(client *http.Client, appId, appSecret, callbackAddr, appType string) *WeChatClient {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{},
		}
		_wechatClient[appId] = client
	}
	c := &WeChatClient{
		conn:         client,
		appId:        appId,
		appSecret:    appSecret,
		appType:      appType,
		callbackAddr: callbackAddr,
	}
	return c
}

// Close 关闭连接
func (m *WeChatClient) Close() {
	if m.conn != nil {
		m.conn.CloseIdleConnections()
		m.conn = nil
	}
}

func (m *WeChatClient) String() string {
	return m.appId
}

func (m *WeChatClient) GetAccessToken() (string, int, error) {

	var err error
	var accessToken string
	/*
		https://api.weixin.qq.com/cgi-bin/token?
		grant_type=client_credential
		&appid=APPID
		&secret=APPSECRET
	*/
	accessTokenUrl := fmt.Sprintf("%v?appid=%v&secret=%v&grant_type=client_credential",
		WeChatAccessTokenAddr,
		m.appId,
		m.appSecret,
	)
	body, err := httpGet(accessTokenUrl)
	if err != nil {
		return "", 0, err
	}

	accessToken = gjson.Get(body, "access_token").String()
	expiresIn := cast.ToInt(gjson.Get(body, "expires_in").Int())

	return accessToken, expiresIn, nil
}

func (m *WeChatClient) GetJSApiTicket(accessToken string) (string, int, error) {

	/*
		https://api.weixin.qq.com/cgi-bin/ticket/getticket?
		access_token=ACCESS_TOKEN
		&type=jsapi
	*/
	ticketUrl := fmt.Sprintf("%v?type=%v&access_token=%v",
		WeChatJSApiTicketAddr,
		"jsapi",
		accessToken,
	)
	body, err := httpGet(ticketUrl)
	if err != nil {
		return "", 0, err
	}

	ticket := gjson.Get(body, "ticket").String()
	expiresIn := cast.ToInt(gjson.Get(body, "expires_in").Int())

	return ticket, expiresIn, nil
}

func (m *WeChatClient) GetConnectUrl(state, scope string, popUp bool) (string, error) {

	var connectUrl string
	if strings.EqualFold(m.appType, "open") {
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
		connectUrl = fmt.Sprintf("%v?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v&forcePopup=%v#wechat_redirect",
			WeChatQRConnectAddr,
			m.appId,
			m.callbackAddr,
			scope,
			state,
			popUp,
		)
	} else {
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
		connectUrl = fmt.Sprintf("%v?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v&forcePopup=%v#wechat_redirect",
			WeChatAppConnectAddr,
			m.appId,
			url.QueryEscape(m.callbackAddr),
			scope,
			state,
			popUp,
		)
	}
	return connectUrl, nil
}

func (m *WeChatClient) GetUserAccessToken(code, scope string) (string, error) {
	/*
		https://api.weixin.qq.com/sns/oauth2/access_token?
		appid=APPID
		&secret=SECRET
		&code=CODE
		&grant_type=authorization_code
	*/
	userAccessTokenUrl := fmt.Sprintf("%v?appid=%v&secret=%v&code=%v&grant_type=authorization_code",
		WeChatUserAccessTokenAddr,
		m.appId,
		m.appSecret,
		code,
	)

	body, err := httpGet(userAccessTokenUrl)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatClient) RefreshUserToken(refreshToken string, appId string) (string, error) {
	/*
		https://api.weixin.qq.com/sns/oauth2/refresh_token?
		appid=APPID
		&grant_type=refresh_token
		&refresh_token=REFRESH_TOKEN
	*/
	refreshUserTokenUrl := fmt.Sprintf("%v?appid=%v&refresh_token=%v&grant_type=refresh_token",
		WeChatRefreshUserTokenAddr,
		appId,
		refreshToken,
	)

	body, err := httpGet(refreshUserTokenUrl)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatClient) GetUserInfo(userAccessToken, openId string) (string, error) {

	/*
		https://api.weixin.qq.com/sns/userinfo?
		access_token=ACCESS_TOKEN
		&openid=OPENID
		&lang=zh_CN
	*/
	userInfoUrl := fmt.Sprintf("%v?access_token=%v&openid=%v&lang=zh_CN",
		WeChatUserInfoAddr,
		userAccessToken,
		openId,
	)

	body, err := httpGet(userInfoUrl)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatClient) GetSubscribeUserInfo(accessToken, openId string) (string, error) {

	/*
		https://api.weixin.qq.com/cgi-bin/user/info
		?access_token=%s
		&openid=%s
		&lang=zh_CN
	*/
	userInfoUrl := fmt.Sprintf("%v?access_token=%v&openid=%v&lang=zh_CN",
		WeChatSubscribeUserInfoAddr,
		accessToken,
		openId,
	)

	body, err := httpGet(userInfoUrl)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (m *WeChatClient) SendTemplateMessage(accessToken, openId, templateId, redirectUrl string, data []byte, rootData []byte) (string, error) {

	if strings.EqualFold(m.appType, "open") {
		return "", errors.New("WeChat open app not support send template message")
	}

	/*
			https://api.weixin.qq.com/cgi-bin/message/template/send?
			access_token=ACCESS_TOKEN
			{
			   "touser":"OPENID",
			   "template_id":"ngqIpbwh8bUfcSsECmogfXcV14J0tQlEpBO27izEYtY",
			   "url":"http://weixin.qq.com/download",
			   "miniprogram":{
				 "appid":"xiaochengxuappid12345",
				 "pagepath":"index?foo=bar"
			   },
			   "data":{
					   "first": {
						   "value":"恭喜你购买成功！",
						   "color":"#173177"
					   },
					   "keyword1":{
						   "value":"巧克力",
						   "color":"#173177"
					   },
					   "keyword2": {
						   "value":"39.8元",
						   "color":"#173177"
					   },
					   "keyword3": {
						   "value":"2014年9月22日",
						   "color":"#173177"
					   },
					   "remark":{
						   "value":"欢迎再次购买！",
						   "color":"#173177"
					   }
			   }
		   }
	*/
	sendTemplateUrl := fmt.Sprintf("%v?access_token=%v", WeChatTemplateMessageAddr, accessToken)

	var body []byte
	if rootData != nil {
		body = rootData
	}
	body, _ = sjson.SetBytes(body, "touser", openId)
	body, _ = sjson.SetBytes(body, "template_id", templateId)
	body, _ = sjson.SetBytes(body, "url", redirectUrl)
	body, _ = sjson.SetRawBytes(body, "data", data)

	resp, err := httpPost(sendTemplateUrl, string(body))
	if err != nil {
		return "", err
	}

	return resp, nil
}

// GetClient 暴露原生client
func (m *WeChatClient) GetClient() *http.Client {
	return m.conn
}
