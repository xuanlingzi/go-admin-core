package file_store

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

var _tencentFile *cos.Client

func GetTencentFileClient() *cos.Client {
	return _tencentFile
}

type TencentFileClient struct {
	client *cos.Client
}

func NewTencentFile(client *cos.Client, accessKey, secretKey, cosUrl, ciUrl string) *TencentFileClient {
	if client == nil {
		cosURL, err := url.Parse(cosUrl)
		if err != nil {
			panic(err)
		}
		ciURL, err := url.Parse(ciUrl)
		if err != nil {
			panic(err)
		}
		base := &cos.BaseURL{
			BucketURL: cosURL,
			CIURL:     ciURL,
		}
		transport := &cos.AuthorizationTransport{
			SecretID:  accessKey,
			SecretKey: secretKey,
		}
		client = cos.NewClient(base, &http.Client{
			Transport: transport,
		})
	}

	r := &TencentFileClient{
		client: client,
	}
	return r
}

func (rc *TencentFileClient) String() string {
	return "tencent_file"
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
