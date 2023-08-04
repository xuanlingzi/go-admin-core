package file_store

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var _tencentFile *cos.Client

func GetTencentFileClient() *cos.Client {
	return _tencentFile
}

type TencentFileClient struct {
	client   *cos.Client
	accessId string
}

func NewTencentFile(client *cos.Client, accessKey, secretKey, cosUrl, ciUrl string) *TencentFileClient {
	if client == nil {
		cosURL, err := url.Parse(cosUrl)
		if err != nil {
			panic(fmt.Sprintf("Tencent file store init error: %s", err.Error()))
		}
		ciURL, err := url.Parse(ciUrl)
		if err != nil {
			panic(fmt.Sprintf("Tencent file store init error: %s", err.Error()))
		}
		base := &cos.BaseURL{
			BucketURL: cosURL,
			CIURL:     ciURL,
		}
		transport := &cos.AuthorizationTransport{
			SecretID:  accessKey,
			SecretKey: secretKey,
			Transport: &http.Transport{},
		}
		client = cos.NewClient(base, &http.Client{
			Transport: transport,
		})
		_tencentFile = client
	}

	r := &TencentFileClient{
		client:   client,
		accessId: accessKey,
	}
	return r
}

func (rc *TencentFileClient) String() string {
	return rc.accessId
}

func (rc *TencentFileClient) Check() bool {
	return rc.client != nil
}

func (rc *TencentFileClient) Close() {

}

func (rc *TencentFileClient) Upload(ctx context.Context, name, fileLocation string) (string, error) {
	response, err := rc.client.Object.PutFromFile(ctx, name, fileLocation, nil)
	if err != nil {
		return fileLocation, err
	}
	defer response.Body.Close()

	return response.Response.Request.URL.String(), nil
}

// GetClient 暴露原生client
func (rc *TencentFileClient) GetClient() interface{} {
	return rc.client
}
