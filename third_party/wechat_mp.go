package third_party

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"net/url"
)

var _wechatMp *http.Client

type WeChatMp struct {
	conn         *http.Client
	appId        string
	appSecret    string
	callbackAddr string
}

func GetWeChatMpClient() *http.Client {
	return _wechatMp
}

func NewWeChatMp(client *http.Client, appId, appSecret, callbackAddr string) *WeChatMp {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{},
		}
	}
	c := &WeChatMp{
		conn:         client,
		appId:        appId,
		appSecret:    appSecret,
		callbackAddr: callbackAddr,
	}
	return c
}

// Close 关闭连接
func (m *WeChatMp) Close() {
	if m.conn != nil {
		m.conn.CloseIdleConnections()
		m.conn = nil
	}
}

func (*WeChatMp) String() string {
	return "wechat_mp"
}

func (m *WeChatMp) GetAccessToken() (string, int, error) {

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

func (m *WeChatMp) GetJSApiTicket(accessToken string) (string, int, error) {

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

func (m *WeChatMp) GetConnectUrl(state, scope string, popUp bool) (string, error) {
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
		WeChatAppConnectAddr,
		m.appId,
		url.QueryEscape(m.callbackAddr),
		scope,
		state,
		popUp,
	)
	return url, nil
}

func (m *WeChatMp) GetUserAccessToken(code, scope string) (string, error) {
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

func (m *WeChatMp) RefreshUserToken(refreshToken string, appId string) (string, error) {
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

func (m *WeChatMp) GetUserInfo(accessToken, openId string) (string, error) {
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

func (m *WeChatMp) SendTemplateMessage(accessToken, openId, templateId, redirectUrl string, data []byte) (string, error) {
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
	url := fmt.Sprintf("%v?access_token=%v", WeChatTemplateMessageAddr, accessToken)

	var body []byte
	body, _ = sjson.SetBytes(body, "touser", openId)
	body, _ = sjson.SetBytes(body, "template_id", templateId)
	body, _ = sjson.SetBytes(body, "url", redirectUrl)
	body, _ = sjson.SetRawBytes(body, "data", data)

	resp, err := httpPost(url, string(body))
	if err != nil {
		return "", err
	}

	return resp, nil
}

// GetClient 暴露原生client
func (m *WeChatMp) GetClient() *http.Client {
	return m.conn
}
