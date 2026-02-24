package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type VerifierAuth struct {
}

func NewVerifierAuth() *VerifierAuth {
	return &VerifierAuth{}
}

func (v *VerifierAuth) un() int64 {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	return time.Now().UnixNano() / 1000 / 30
}

func (v *VerifierAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (v *VerifierAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (v *VerifierAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (v *VerifierAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (v *VerifierAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (v *VerifierAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := v.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := v.toUint32(hashParts)
	return number % 1000000
}

// GetSecret 获取秘钥
func (v *VerifierAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, v.un())
	return strings.ToUpper(v.base32encode(v.hmacSha1(buf.Bytes(), nil)))
}

// GetCode 获取动态码
func (v *VerifierAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := v.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	number := v.oneTimePassword(secretKey, v.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

// GetQrcode 获取动态码二维码内容
func (v *VerifierAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// GetQrcodeUrl 获取动态码二维码图片地址,这里是第三方二维码api
func (v *VerifierAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := v.GetQrcode(user, secret)
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://www.google.com/chart?chs=200x200&chld=M|0&cht=qr&chl=" + data.Encode()
}

// VerifyCode 验证动态码
func (v *VerifierAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := v.GetCode(secret)
	fmt.Println(_code, code, err)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
