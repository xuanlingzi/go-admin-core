package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/xuanlingzi/go-admin-core/logger"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

func Hmac(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func IsStringEmpty(str string) bool {
	return strings.Trim(str, " ") == ""
}

func GetUUID() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), "-", "")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func Base64ToImage(imageBase64 string) ([]byte, error) {
	image, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func GetDirFiles(dir string) ([]string, error) {
	dirList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesRet := make([]string, 0)

	for _, file := range dirList {
		if file.IsDir() {
			files, err := GetDirFiles(dir + string(os.PathSeparator) + file.Name())
			if err != nil {
				return nil, err
			}

			filesRet = append(filesRet, files...)
		} else {
			filesRet = append(filesRet, dir+string(os.PathSeparator)+file.Name())
		}
	}

	return filesRet, nil
}

func GetCurrentTimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// slice去重
func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func ClientIpAddr(r *http.Request) string {
	xForwardFor := r.Header.Get("X-Forwarded-For")
	ipAddr := strings.TrimSpace(strings.Split(xForwardFor, ",")[0])
	if StringIsNotEmpty(ipAddr) {
		return ipAddr
	}

	ipAddr = r.Header.Get("X-Real-IP")
	if StringIsNotEmpty(ipAddr) {
		return ipAddr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func UserAgentToOsBrowser(userAgent string) (os string, browser string) {

	lIdx := strings.Index(userAgent, "(")
	rIdx := strings.Index(userAgent, ")")

	os = strings.TrimSpace(userAgent[lIdx+1 : rIdx])
	browser = strings.TrimSpace(userAgent[rIdx+1:])

	return
}

func ResolveIPFromHostsFile() (string, error) {
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		logger.Errorf("Problem reading /etc/hosts: %v", err.Error())
		return "", fmt.Errorf("problem reading /etc/hosts: %v", err.Error())
	}

	lines := strings.Split(string(data), "\n")
	line := lines[len(lines)-1]
	if len(line) < 2 {
		line = lines[len(lines)-2]
	}

	parts := strings.Split(line, "\t")
	return parts[0], nil
}

func GetIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func GetIPWithPrefix(prefix string) string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), prefix) {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func PasswordCheck(ps string, minLength, maxLength int, level int) error {
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if len(ps) < minLength || len(ps) > maxLength {
		return fmt.Errorf("密码长度必须大于%d，小于%d个字符", minLength, maxLength)
	}
	if level > 0 {
		if b, err := regexp.MatchString(num, ps); !b || err != nil {
			return fmt.Errorf("密码必须包含数字")
		}
	}
	if level > 1 {
		if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
			return fmt.Errorf("密码必须包含小写字母")
		}
		if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
			return fmt.Errorf("密码必须包含大写字母")
		}
	}
	if level > 2 {
		if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
			return fmt.Errorf("密码必须包含特殊字符")
		}
	}
	return nil
}

func StringIsEmpty(text string) bool {
	return text == "" || strings.TrimSpace(text) == ""
}

func StringIsNotEmpty(text string) bool {
	return !StringIsEmpty(text)
}

func HashFile(filename string) (string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return filename, err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return filename, err
	}

	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

func HashString(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func Sha1String(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func FastSimpleUUID() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(uuid.String(), "-", "")
}

func GenerateRandomCode(num int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, num)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate random code: %w", err)
		}
		b[i] = digits[n.Int64()]
	}
	return string(b), nil
}

func GenerateRandomString(num int) (string, error) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, num)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("generate random string: %w", err)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func CompareSignature(signature, message, accessSecret string) bool {
	expectedMac := HmacSignature(message, accessSecret)
	return hmac.Equal([]byte(signature), expectedMac)
}

func HmacSignature(message, accessSecret string) []byte {
	mac := hmac.New(sha1.New, []byte(accessSecret))
	mac.Write([]byte(message))
	return mac.Sum(nil)
}

func UniqueId(prefix string) string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}
	now := time.Now().In(loc)
	timeString := strings.Replace(now.Format("20060102150405.999"), ".", "", 1)
	code, err := GenerateRandomCode(4)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s%s", prefix, timeString, code)
}

func Copy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HttpPost(url, content string) ([]byte, error) {

	response, err := http.Post(url, "application/json", strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return origin
	}
	var (
		searchLength  = len(search)
		replaceLength = len(replace)
		searchLower   = strings.ToLower(search)
		originLower   string
		pos           int
	)
	for {
		originLower = strings.ToLower(origin)
		if pos = Pos(originLower, searchLower, pos); pos != -1 {
			origin = origin[:pos] + replace + origin[pos+searchLength:]
			pos += replaceLength
			if n--; n == 0 {
				break
			}
		} else {
			break
		}
	}
	return origin
}

func Pos(haystack, needle string, startOffset ...int) int {
	length := len(haystack)
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	if length == 0 || offset > length || -offset > length {
		return -1
	}
	if offset < 0 {
		offset += length
	}
	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}
	return pos + offset
}
