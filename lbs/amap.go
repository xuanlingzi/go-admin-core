package lbs

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xuanlingzi/go-admin-core/tools/utils"
	"net/http"
)

var _amap *http.Client

type Amap struct {
	conn      *http.Client
	addr      string
	secretKey string
}

func GetAmapClient() *http.Client {
	return _amap
}

func NewAmap(client *http.Client, addr, secretKey string) *Amap {
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
	return c
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

	addressInfo, err := utils.HttpGet(url)
	if err != nil {
		return "", err
	}

	if gjson.GetBytes(addressInfo, "status").Exists() && gjson.GetBytes(addressInfo, "status").Int() != 1 {
		return "", errors.New(gjson.GetBytes(addressInfo, "info").String())
	}

	if gjson.GetBytes(addressInfo, "regeocode.formatted_address").Exists() {
		return gjson.GetBytes(addressInfo, "regeocode.formatted_address").String(), nil
	}
	addr := ""
	if gjson.GetBytes(addressInfo, "regeocode.addressComponent").Exists() {
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.country").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.province").String()
		if gjson.GetBytes(addressInfo, "regeocode.addressComponent.city").Exists() {
			city := gjson.GetBytes(addressInfo, "regeocode.addressComponent.city").String()
			if city != "[]" {
				addr += city
			}
		}
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.district").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.township").String()
		if gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber").Exists() {
			addr += fmt.Sprintf("%v%v",
				gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber.street").String(),
				gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber.number").String(),
			)
		}
	}

	return addr, nil
}

// GetClient 暴露原生client
func (m *Amap) GetClient() interface{} {
	return m.conn
}
