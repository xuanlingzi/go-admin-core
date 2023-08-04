package moderation

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	moderation "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/region"
	"github.com/tidwall/gjson"
	"strings"
)

var _huaweiAudit *moderation.ModerationClient

func GetHuaweiAuditClient() *moderation.ModerationClient {
	return _huaweiAudit
}

type HuaweiAuditClient struct {
	client      *moderation.ModerationClient
	accessId    string
	callbackUrl string
}

func NewHuaweiAudit(client *moderation.ModerationClient, accessKey, secretKey, _region, callbackUrl string) *HuaweiAuditClient {
	var err error
	if client == nil {
		auth := basic.NewCredentialsBuilder().
			WithAk(accessKey).
			WithSk(secretKey).
			Build()

		client = moderation.NewModerationClient(moderation.ModerationClientBuilder().WithRegion(region.ValueOf(_region)).WithCredential(auth).Build())
		if err != nil {
			panic(fmt.Sprintf("Huawei audit init error: %s", err.Error()))
		}
		_huaweiAudit = client
	}

	r := &HuaweiAuditClient{
		client:      client,
		accessId:    accessKey,
		callbackUrl: callbackUrl,
	}

	return r
}

func (rc *HuaweiAuditClient) String() string {
	return rc.accessId
}

func (rc *HuaweiAuditClient) Check() bool {
	return rc.client != nil
}

func (rc *HuaweiAuditClient) Close() {
}

func (rc *HuaweiAuditClient) AuditText(content string, result *int, label *string, detail *string) error {

	requestData := model.TextDetectionDataReq{
		Text: content,
	}
	textReq := model.TextDetectionReq{
		Data:      &requestData,
		EventType: "comment",
	}
	request := model.RunTextModerationRequest{}
	request.Body = &textReq
	response, err := rc.client.RunTextModeration(&request)
	if err != nil {
		return fmt.Errorf("huawei Audit error: %v", err.Error())
	}

	switch *response.Result.Suggestion {
	case "pass":
		*result = 0
	case "block":
		*result = 1
	case "review":
		*result = 2
	default:
		*result = 1
	}

	if response.Result.Label != nil {
		*label = *response.Result.Label
	}
	if response.Result.Details != nil {
		b, err := json.Marshal(response.Result.Details)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *HuaweiAuditClient) AuditImage(url string, result *int, label *string, detail *string) error {

	var imageCategories = []string{
		"porn",
		"politics",
		"terrorism",
	}
	imageReq := model.ImageDetectionReq{
		Url:        &url,
		Categories: imageCategories,
		EventType:  "album",
	}
	request := model.CheckImageModerationRequest{}
	request.Body = &imageReq
	response, err := rc.client.CheckImageModeration(&request)
	if err != nil {
		return fmt.Errorf("huawei Audit error: %v", err.Error())
	}

	switch *response.Result.Suggestion {
	case "pass":
		*result = 0
	case "block":
		*result = 1
	case "review":
		*result = 2
	default:
		*result = 1
	}

	if response.Result.Category != nil {
		*label = *response.Result.Category
	}
	if response.Result.Details != nil {
		b, err := json.Marshal(response.Result.Details)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *HuaweiAuditClient) AuditVideo(url string, frame int32, jobId *string) error {

	var audioCategories = []model.VideoCreateRequestAudioCategories{
		model.GetVideoCreateRequestAudioCategoriesEnum().PORN,
		model.GetVideoCreateRequestAudioCategoriesEnum().POLITICS,
		model.GetVideoCreateRequestAudioCategoriesEnum().ABUSE,
		model.GetVideoCreateRequestAudioCategoriesEnum().MOAN,
	}
	var imageCategories = []model.VideoCreateRequestImageCategories{
		model.GetVideoCreateRequestImageCategoriesEnum().PORN,
		model.GetVideoCreateRequestImageCategoriesEnum().POLITICS,
		model.GetVideoCreateRequestImageCategoriesEnum().TERRORISM,
		model.GetVideoCreateRequestImageCategoriesEnum().IMAGE_TEXT,
	}
	requestData := model.VideoCreateRequestData{
		Url:           url,
		FrameInterval: &frame,
	}
	videoReq := model.VideoCreateRequest{
		Callback:        &rc.callbackUrl,
		AudioCategories: &audioCategories,
		ImageCategories: imageCategories,
		EventType:       model.GetVideoCreateRequestEventTypeEnum().DEFAULT,
		Data:            &requestData,
	}

	request := model.RunCreateVideoModerationJobRequest{}
	request.Body = &videoReq
	response, err := rc.client.RunCreateVideoModerationJob(&request)
	if err != nil {
		return fmt.Errorf("huawei Audit error: %v", err.Error())
	}
	*jobId = *response.JobId
	return nil
}

func (rc *HuaweiAuditClient) AuditResult(body []byte, result *int, label *string, detail *string, jobId *string) error {

	if gjson.GetBytes(body, "result.job_id").Exists() {
		*jobId = gjson.GetBytes(body, "result.job_id").String()
	}

	switch strings.ToLower(gjson.GetBytes(body, "result.suggestion").String()) {
	case "pass":
		*result = 0
	case "block":
		*result = 1
	case "review":
		*result = 2
	default:
		*result = 1
	}

	if *result > 0 {
		if gjson.GetBytes(body, "result.label").Exists() {
			*label = gjson.GetBytes(body, "result.label").String()
		}

		if gjson.GetBytes(body, "result.image_detail").Exists() {
			if *label == "" {
				imageDetails := gjson.GetBytes(body, "result.image_detail").Array()
				if len(imageDetails) > 0 {
					*label = imageDetails[0].Get("category").String()
				}
			}

			*detail = gjson.GetBytes(body, "result.image_detail").String()
		}
	}

	return nil
}

// GetClient 暴露原生client
func (rc *HuaweiAuditClient) GetClient() interface{} {
	return rc.client
}
