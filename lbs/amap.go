package lbs

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

type Amap struct {
	conn      *http.Client
	addr      string
	secretKey string
}

func NewAmap(client *http.Client, addr, secretKey string) (*Amap, error) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{},
		}
	}
	c := &Amap{
		conn:      client,
		addr:      addr,
		secretKey: secretKey,
	}
	return c, nil
}

// Close 关闭连接
func (m *Amap) Close() {
	if m.conn != nil {
		m.conn.CloseIdleConnections()
		m.conn = nil
	}
}

func (*Amap) String() string {
	return "amap"
}

func (m *Amap) GetAddress(latitude, longitude, radius float64) (string, error) {
	/*
		http://restapi.amap.com/v3/geocode/regeo?output=JSON
		&location=" + longitude + "," + latitude + "
		&key=
		&radius=
		&extensions=all
		&batch=false
		&roadlevel=0";
	*/
	url := fmt.Sprintf("%v/v3/geocode/regeo?output=JSON&key=%v&location=%v,%v&radius=%v&extensions=all&batch=false&roadlevel=0",
		m.addr,
		m.secretKey,
		longitude,
		latitude,
		radius,
	)

	addressInfo, err := httpGet(url)
	if err != nil {
		return "", err
	}

	if gjson.Get(addressInfo, "regeocode.formatted_address").Exists() {
		return gjson.Get(addressInfo, "regeocode.formatted_address").String(), nil
	}
	addr := ""
	if gjson.Get(addressInfo, "regeocode.addressComponent").Exists() {
		addr += gjson.Get(addressInfo, "regeocode.addressComponent.country").String()
		addr += gjson.Get(addressInfo, "regeocode.addressComponent.province").String()
		if gjson.Get(addressInfo, "regeocode.addressComponent.city").Exists() {
			city := gjson.Get(addressInfo, "regeocode.addressComponent.city").String()
			if city != "[]" {
				addr += city
			}
		}
		addr += gjson.Get(addressInfo, "regeocode.addressComponent.district").String()
		addr += gjson.Get(addressInfo, "regeocode.addressComponent.township").String()
		if gjson.Get(addressInfo, "regeocode.addressComponent.streetNumber").Exists() {
			addr += fmt.Sprintf("%v%v",
				gjson.Get(addressInfo, "regeocode.addressComponent.streetNumber.street").String(),
				gjson.Get(addressInfo, "regeocode.addressComponent.streetNumber.number").String(),
			)
		}
	}

	return addr, nil
}

func httpGet(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if gjson.GetBytes(body, "status").Exists() && gjson.GetBytes(body, "status").Int() != 1 {
		return "", errors.New(gjson.GetBytes(body, "info").String())
	}

	return string(body), nil
}

func httpPost(url, content string) (string, error) {
	response, err := http.Post(url, "application/json", strings.NewReader(content))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if gjson.GetBytes(body, "status").Exists() && gjson.GetBytes(body, "status").Int() != 1 {
		return "", errors.New(gjson.GetBytes(body, "info").String())
	}

	return string(body), nil
}

// GetClient 暴露原生client
func (m *Amap) GetClient() *http.Client {
	return m.conn
}
