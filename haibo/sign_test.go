package haibo

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
	"testing"
)

// 黄金用例说明（Requirement 2.7）：
//
// 海博对接文档「7. 请求示例」给出了 appId / body / timestamp 三个入参，但：
//   - 文档展示的 sign 值 "14E14141234567892C54704EB" 仅 25 个十六进制字符，
//     并非合法的 32 位 MD5 摘要，显然是脱敏占位值；
//   - 文档未公开任何 secret。
//
// 因此文档不存在「可直接断言的外部 sign 黄金向量」。为避免 fabricate 一个
// 海博 java-sdk 永远不会产出的期望值，这里采用「独立实现交叉校验」的方式固定
// 黄金向量：
//  1. 入参 appId/body/timestamp 直接取自文档「请求示例」；
//  2. secret 取测试固定值（文档未公开真实 secret）；
//  3. 期望 sign 由 Go 之外的独立 MD5 实现（Python hashlib）按文档拼接规则
//     secret+"appId"+appId+"body"+body+"signMethod"md5"timestamp"+timestamp+secret
//     预先算出，硬编码于此，确保不是用被测函数自身回算得到的（非循环自证）。
//
// 拼接串：
//
//	testSecret123456appId4de67a45194f12345678968195fbcbody{"stationNo":"1234",
//	"skuInfos":[{"skuId":100671234},{"skuId":100671235}]}signMethodmd5
//	timestamp7888899531021testSecret123456
const (
	goldenAppID     = "4de67a45194f12345678968195fbc"
	goldenBody      = `{"stationNo":"1234","skuInfos":[{"skuId":100671234},{"skuId":100671235}]}`
	goldenTimestamp = "7888899531021"
	goldenSecret    = "testSecret123456"
	// goldenSign 由独立的 Python hashlib.md5 实现按文档拼接规则预先算出。
	goldenSign = "6A18034DF192261C1706D1D9E1F82025"
)

// TestCalcSign_GoldenVector 断言 CalcSign 对文档「请求示例」参数的输出与独立实现
// 预先算出的期望 sign 逐字符完全相同（Requirement 2.7）。
func TestCalcSign_GoldenVector(t *testing.T) {
	got := CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret)
	if got != goldenSign {
		t.Fatalf("CalcSign 与黄金向量不一致:\n  got  = %q\n  want = %q", got, goldenSign)
	}
}

// TestCalcSign_MatchesDocumentedConcatRule 独立按文档拼接规则在测试内重算一次
// （不调用被测函数的内部拼接），与 CalcSign 比对，确认拼接顺序与首尾 secret 正确。
func TestCalcSign_MatchesDocumentedConcatRule(t *testing.T) {
	// secret + "appId"+appId + "body"+body + "signMethod"+md5 + "timestamp"+ts + secret
	concat := goldenSecret +
		"appId" + goldenAppID +
		"body" + goldenBody +
		"signMethod" + signMethodMD5 +
		"timestamp" + goldenTimestamp +
		goldenSecret
	sum := md5.Sum([]byte(concat))
	want := strings.ToUpper(hex.EncodeToString(sum[:]))

	got := CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret)
	if got != want {
		t.Fatalf("CalcSign 与按文档规则独立重算结果不一致:\n  got  = %q\n  want = %q", got, want)
	}
}

// TestCalcSign_Deterministic 相同入参必须恒定产出相同结果（Requirement 2.7 一致性）。
func TestCalcSign_Deterministic(t *testing.T) {
	first := CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret)
	for i := 0; i < 100; i++ {
		again := CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret)
		if again != first {
			t.Fatalf("CalcSign 非确定性：第 %d 次结果 %q != 首次 %q", i, again, first)
		}
	}
}

// upper32Hex 校验 sign 为 32 位大写十六进制（Requirement 2.3）。
var upper32Hex = regexp.MustCompile(`^[0-9A-F]{32}$`)

// TestCalcSign_Is32UppercaseHex 断言输出恒为 32 位、全大写十六进制字符串。
func TestCalcSign_Is32UppercaseHex(t *testing.T) {
	cases := []struct {
		name                                string
		appId, body, signMethod, ts, secret string
	}{
		{"golden", goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret},
		{"empty-all", "", "", "", "", ""},
		{"unicode-body", "app", `{"name":"测试商品"}`, signMethodMD5, "1700000000000", "秘钥"},
		{"long", strings.Repeat("a", 256), strings.Repeat("b", 1024), signMethodMD5, "1548645174043", strings.Repeat("c", 64)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalcSign(tc.appId, tc.body, tc.signMethod, tc.ts, tc.secret)
			if !upper32Hex.MatchString(got) {
				t.Fatalf("sign 非 32 位大写十六进制: %q", got)
			}
		})
	}
}

// TestCalcSign_SensitiveToEachInput 任一入参变化都应改变签名，确认每个字段都参与拼接。
func TestCalcSign_SensitiveToEachInput(t *testing.T) {
	base := CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret)

	variants := map[string]string{
		"appId":     CalcSign(goldenAppID+"x", goldenBody, signMethodMD5, goldenTimestamp, goldenSecret),
		"body":      CalcSign(goldenAppID, goldenBody+"x", signMethodMD5, goldenTimestamp, goldenSecret),
		"timestamp": CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp+"1", goldenSecret),
		"secret":    CalcSign(goldenAppID, goldenBody, signMethodMD5, goldenTimestamp, goldenSecret+"x"),
	}
	for field, got := range variants {
		if got == base {
			t.Errorf("修改 %s 后签名未变化，该字段可能未参与拼接", field)
		}
	}
}
