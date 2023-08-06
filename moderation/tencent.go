package moderation

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var _tencentAudit *cos.Client

func GetTencentAuditClient() *cos.Client {
	return _tencentAudit
}

type TencentAuditClient struct {
	client      *cos.Client
	accessId    string
	callbackUrl string
}

func NewTencentAudit(client *cos.Client, accessKey, secretKey, cosUrl, ciUrl, callbackUrl string) *TencentAuditClient {
	if client == nil {
		cosURL, err := url.Parse(cosUrl)
		if err != nil {
			panic(fmt.Sprintf("Tencent Audit init error: %s", err.Error()))
		}
		ciURL, err := url.Parse(ciUrl)
		if err != nil {
			panic(fmt.Sprintf("Tencent Audit init error: %s", err.Error()))
		}

		cosBaseUrl := &cos.BaseURL{
			BucketURL: cosURL,
			CIURL:     ciURL,
		}
		transport := &cos.AuthorizationTransport{
			SecretID:  accessKey,
			SecretKey: secretKey,
			Transport: &http.Transport{},
		}
		client = cos.NewClient(cosBaseUrl, &http.Client{
			Transport: transport,
		})
		_tencentAudit = client
	}

	r := &TencentAuditClient{
		client:      client,
		accessId:    accessKey,
		callbackUrl: callbackUrl,
	}
	return r
}

func (rc *TencentAuditClient) String() string {
	return rc.accessId
}

func (rc *TencentAuditClient) Check() bool {
	return rc.client != nil
}

func (rc *TencentAuditClient) Close() {

}

func (rc *TencentAuditClient) AuditText(content string, result *int, label *string, score *int, detail *string) error {

	bString := base64.StdEncoding.EncodeToString([]byte(content))
	opt := cos.PutTextAuditingJobOptions{
		InputContent: bString,
		Conf: &cos.TextAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Teenager,Illegal,Abuse",
		},
	}
	response, _, err := rc.client.CI.PutTextAuditingJob(context.TODO(), &opt)
	if err != nil {
		return err
	}

	*result = response.JobsDetail.Result
	*label = strings.ToLower(response.JobsDetail.Label)
	if response.JobsDetail != nil {
		b, err := json.Marshal(response.JobsDetail)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *TencentAuditClient) AuditImage(url string, result *int, label *string, score *int, detail *string) error {

	_, filename := filepath.Split(url)

	opt := &cos.ImageRecognitionOptions{
		CIProcess:  "sensitive-content-recognition",
		DetectType: "porn,terrorist,politics,terrorism,teenager",
		Interval:   5,
		MaxFrames:  100,
	}
	response, _, err := rc.client.CI.ImageAuditing(context.TODO(), filename, opt)
	if err != nil {
		return err
	}

	*result = response.Result
	*label = strings.ToLower(response.Label)
	if response != nil {
		b, err := json.Marshal(response)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *TencentAuditClient) AuditVideo(url string, frame int32, jobId *string) error {

	_, filename := filepath.Split(url)

	opt := &cos.PutVideoAuditingJobOptions{
		InputObject: filename,
		Conf: &cos.VideoAuditingJobConf{
			DetectType: "Porn,Terrorism,Politics,Teenager", //Terrorist",
			Snapshot: &cos.PutVideoAuditingJobSnapshot{
				/*
					截帧模式，默认值为Interval。
					Interval 表示间隔模式；
					Average 表示平均模式；
					Fps 表示固定帧率模式。
					Interval 模式：TimeInterval，Count 参数生效。当设置 Count，未设置 TimeInterval 时，表示截取所有帧，共 Count 张图片。
					Average 模式：Count 参数生效。表示整个视频，按平均间隔截取共 Count 张图片。
					Fps 模式：TimeInterval 表示每秒截取多少帧，未设置 TimeInterval 时，表示截取所有帧，Count 表示共截取多少帧。
				*/
				Mode:         "Fps",
				TimeInterval: cast.ToFloat32(frame),
				Count:        10000,
			},
			Callback:        rc.callbackUrl,
			CallbackVersion: "Detail", // Simple（回调内容包含基本信息）、Detail（回调内容包含详细信息）。默认为 Simple
			DetectContent:   1,        // 当值为0时：表示只审核视频画面截图；值为1时：表示同时审核视频画面截图和视频声音。默认值为0
		},
	}
	response, _, err := rc.client.CI.PutVideoAuditingJob(context.TODO(), opt)
	if err != nil {
		return err
	}
	*jobId = response.JobsDetail.JobId

	return nil
}

func (rc *TencentAuditClient) AuditResult(body *[]byte, result *int, label *string, score *int, detail *string, jobId *string) error {

	*jobId = gjson.GetBytes(*body, "JobsDetail.JobId").String()
	*label = gjson.GetBytes(*body, "JobsDetail.Label").String()
	*result = cast.ToInt(gjson.GetBytes(*body, "JobsDetail.Result").Int())

	*detail = gjson.GetBytes(*body, "JobsDetail").String()

	return nil
}

// GetClient 暴露原生client
func (rc *TencentAuditClient) GetClient() interface{} {
	return rc.client
}

func genInfo(info *cos.RecognitionInfo) string {

	result := "{}"
	if info != nil {
		result, _ = sjson.Set(result, "code", info.Code)
		result, _ = sjson.Set(result, "msg", info.Msg)
		result, _ = sjson.Set(result, "hit_flag", info.HitFlag)
		result, _ = sjson.Set(result, "score", info.Score)
		result, _ = sjson.Set(result, "label", info.Label)
		result, _ = sjson.Set(result, "sub_label", info.SubLabel)
		result, _ = sjson.Set(result, "count", info.Count)
		result, _ = sjson.Set(result, "keywords", strings.Join(info.Keywords, ","))
	}

	return result
}
