package lbs

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/xuanlingzi/go-admin-core/tools/utils"
	"net/http"
	"strings"
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
		_amap = client
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

func (m *Amap) GetAddress(latitude, longitude, radius float64) (map[string]string, error) {
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

	result := make(map[string]string)

	addressInfo, err := utils.HttpGet(url)
	if err != nil {
		return result, err
	}

	if gjson.GetBytes(addressInfo, "status").Exists() && gjson.GetBytes(addressInfo, "status").Int() != 1 {
		return result, errors.New(gjson.GetBytes(addressInfo, "info").String())
	}

	addr := ""
	if gjson.GetBytes(addressInfo, "regeocode.addressComponent").Exists() {
		result["country"] = gjson.GetBytes(addressInfo, "regeocode.addressComponent.country").String()
		result["province"] = gjson.GetBytes(addressInfo, "regeocode.addressComponent.province").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.country").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.province").String()
		if gjson.GetBytes(addressInfo, "regeocode.addressComponent.city").Exists() {
			city := gjson.GetBytes(addressInfo, "regeocode.addressComponent.city").String()
			if city != "[]" {
				result["city"] = city
				addr += city
			}
		}
		result["district"] = gjson.GetBytes(addressInfo, "regeocode.addressComponent.district").String()
		result["township"] = gjson.GetBytes(addressInfo, "regeocode.addressComponent.township").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.district").String()
		addr += gjson.GetBytes(addressInfo, "regeocode.addressComponent.township").String()
		if gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber").Exists() {
			streetNumber := fmt.Sprintf("%v%v",
				gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber.street").String(),
				gjson.GetBytes(addressInfo, "regeocode.addressComponent.streetNumber.number").String(),
			)
			addr += streetNumber
			result["street"] = streetNumber
		}
	}

	if gjson.GetBytes(addressInfo, "regeocode.formatted_address").Exists() {
		result["address"] = gjson.GetBytes(addressInfo, "regeocode.formatted_address").String()
	} else {
		result["address"] = addr
	}

	return result, nil
}

func (m *Amap) GetCoordinate(keyword string) (longitude float32, latitude float32, address string, err error) {
	/*
		http://restapi.amap.com/v3/geocode/geo?output=JSON
		&address=" + address;
	*/
	url := fmt.Sprintf("%v/v3/geocode/geo?output=JSON&key=%v&address=%v",
		m.addr,
		m.secretKey,
		keyword,
	)

	addressInfo, err := utils.HttpGet(url)
	if err != nil {
		return
	}

	if gjson.GetBytes(addressInfo, "status").Exists() && gjson.GetBytes(addressInfo, "status").Int() != 1 {
		err = errors.New(gjson.GetBytes(addressInfo, "info").String())
		return
	}

	if gjson.GetBytes(addressInfo, "geocodes.#").Int() > 0 {
		address = gjson.GetBytes(addressInfo, "geocodes.0.formatted_address").String()

		loc := gjson.GetBytes(addressInfo, "geocodes.0.location").String()
		longitude = cast.ToFloat32(loc[:strings.Index(loc, ",")])
		latitude = cast.ToFloat32(loc[strings.Index(loc, ",")+1:])
	}

	return
}

func (m *Amap) GetPosition(imei, network, ac, ci, snr string, result *map[string]interface{}) error {
	/*
		http://apilocate.amap.com/position?output=JSON
		&key=%v
		&accesstype=0
		&imei=" + imei + "
		&cdma=0
		&network=" + network + "
		&bts=460,0," + tac + "," + cellId + "," + snr
	*/

	url := fmt.Sprintf("%v/position?output=JSON&key=%v&accesstype=0&imei=%v&cdma=0&network=%v&bts=460,0,%v,%v,%v",
		m.addr,
		m.secretKey,
		imei,
		network,
		ac,
		ci,
		snr,
	)

	addressInfo, err := utils.HttpGet(url)
	if err != nil {
		return fmt.Errorf("GetPosition error %v", err)
	}

	if gjson.GetBytes(addressInfo, "info").Exists() {
		status := gjson.GetBytes(addressInfo, "info").String()
		if strings.EqualFold(status, "OK") == false {
			return fmt.Errorf("GetPosition error %v", addressInfo)
		}
	}

	if gjson.GetBytes(addressInfo, "geocodes.#").Int() > 0 {
		(*result)["address"] = gjson.GetBytes(addressInfo, "result.desc").String()

		loc := gjson.GetBytes(addressInfo, "result.location").String()
		(*result)["longitude"] = cast.ToFloat32(loc[:strings.Index(loc, ",")])
		(*result)["latitude"] = cast.ToFloat32(loc[strings.Index(loc, ",")+1:])
	}

	return nil
}

// GetClient 暴露原生client
func (m *Amap) GetClient() interface{} {
	return m.conn
}
