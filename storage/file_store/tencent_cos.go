package file_store

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
)

var _tencent_cos *cos.Client

func GetTencentCosClient() *cos.Client {
	return _tencent_cos
}

type TencentCosClient struct {
	client *cos.Client
}

func NewTencentCos(client *cos.Client, base *cos.BaseURL, transport *cos.AuthorizationTransport) *TencentCosClient {
	if client == nil {
		client = cos.NewClient(base, &http.Client{
			Transport: transport,
		})
	}

	r := &TencentCosClient{
		client: client,
	}
	return r
}

func (rc *TencentCosClient) String() string {
	return "cos"
}

func (rc *TencentCosClient) Check() bool {
	return rc.client != nil
}

func (rc *TencentCosClient) Close() {

}

func (rc *TencentCosClient) Upload(ctx context.Context, name, fileLocation string) (string, error) {
	response, err := rc.client.Object.PutFromFile(ctx, name, fileLocation, nil)
	if err != nil {
		return fileLocation, err
	}
	defer response.Body.Close()

	return response.Response.Request.URL.String(), nil
}

// GetClient 暴露原生client
func (rc *TencentCosClient) GetClient() interface{} {
	return rc.client
}
