package haibo

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// signMethodMD5 海博签名算法固定取值。
const signMethodMD5 = "md5"

// CalcSign 按海博开放平台规范计算出站请求签名。
//
// 规则（见 Requirement 2.1-2.3）：
//  1. 仅取 appId、body、signMethod、timestamp 四个参数（不含 sign 自身），
//     按该固定顺序拼接为 k1v1k2v2k3v3k4v4 形式：
//     "appId"+appId+"body"+body+"signMethod"+signMethod+"timestamp"+timestamp。
//  2. 在拼接结果首部与尾部各拼接一次海博分配的 secret。
//  3. 对最终字符串计算 MD5 摘要，表示为 32 位十六进制字符串后转为大写，作为 sign。
//
// 该实现复用乐刷 calcTradeSign 的「拼接 → md5 → ToUpper」范式，但拼接规则按海博
// 文档的 secret + k1v1k2v2k3v3k4v4 + secret（固定顺序，非字典序），需与海博
// java-sdk 对相同参数计算所得结果逐字符一致（Requirement 2.7）。
func CalcSign(appId, body, signMethod, timestamp, secret string) string {
	var buf strings.Builder
	buf.WriteString(secret)
	buf.WriteString("appId")
	buf.WriteString(appId)
	buf.WriteString("body")
	buf.WriteString(body)
	buf.WriteString("signMethod")
	buf.WriteString(signMethod)
	buf.WriteString("timestamp")
	buf.WriteString(timestamp)
	buf.WriteString(secret)

	sum := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}
