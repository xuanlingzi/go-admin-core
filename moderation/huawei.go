package moderation

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	moderation "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/region"
	"github.com/spf13/cast"
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

func (rc *HuaweiAuditClient) AuditText(content string, result *int, label *string, score *int, detail *string) error {

	eventType := "article"
	requestData := model.TextDetectionDataReq{
		Text: content,
	}
	textReq := model.TextDetectionReq{
		Data:      &requestData,
		EventType: &eventType,
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

	if response.Result.Details != nil && len(*response.Result.Details) > 0 {
		_detail := *response.Result.Details
		if _detail[0].Confidence != nil {
			*score = cast.ToInt(*_detail[0].Confidence * 100)
		}
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

func (rc *HuaweiAuditClient) AuditImage(url string, result *int, label *string, score *int, detail *string) error {

	eventType := "article"

	var imageCategories = []string{
		"porn",
		"politics",
		"terrorism",
	}
	imageReq := model.ImageDetectionReq{
		Url:        &url,
		Categories: &imageCategories,
		EventType:  &eventType,
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

	if response.Result.Details != nil && len(*response.Result.Details) > 0 {
		_detail := *response.Result.Details
		*score = cast.ToInt(_detail[0].Confidence * 100)
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

	eventType := model.GetVideoCreateRequestEventTypeEnum().DEFAULT

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
		ImageCategories: &imageCategories,
		EventType:       &eventType,
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

func (rc *HuaweiAuditClient) AuditResult(body *[]byte, result *int, label *string, score *int, detail *string, jobId *string) error {

	var resultSuggestion, resultLabel, resultDetail, resultJobId string
	var resultBody []byte
	var resultScore int

	if len(*body) > 0 { // 从回调返回内容

		if gjson.GetBytes(*body, "result.job_id").Exists() {
			resultJobId = gjson.GetBytes(*body, "result.job_id").String()
		}

		if gjson.GetBytes(*body, "result.suggestion").Exists() {
			resultSuggestion = gjson.GetBytes(*body, "result.suggestion").String()
		}

		if gjson.GetBytes(*body, "result.label").Exists() {
			resultLabel = gjson.GetBytes(*body, "result.label").String()
		}

		if gjson.GetBytes(*body, "result.image_detail").Exists() {

			resultDetail = gjson.GetBytes(*body, "result.image_detail").String()

			if resultLabel == "" {
				imageDetails := gjson.GetBytes(*body, "result.image_detail").Array()
				if len(imageDetails) > 0 && imageDetails[0].Get("category").Exists() {
					resultLabel = imageDetails[0].Get("category").String()
				}
			}
		}
	} else if jobId != nil { // 用SDK调用查询结果

		request := model.RunQueryVideoModerationJobRequest{
			JobId: *jobId,
		}
		response, err := rc.client.RunQueryVideoModerationJob(&request)
		if err != nil {
			return fmt.Errorf("huawei Audit error: %v", err.Error())
		}
		if response.Status == nil { // unfinished
			return nil
		}
		if response.Status.Value() == "running" { // unfinished
			return nil
		}
		if response.Status.Value() == "failed" { // failed
			return fmt.Errorf("huawei Audit error: %v", response)
		}

		// 标识已经处理过
		resultBody = []byte(response.Status.Value())
		resultSuggestion = response.Result.Suggestion.Value()
		if response.Result.ImageDetail != nil && len(*response.Result.ImageDetail) > 0 {

			_imageDetail := *response.Result.ImageDetail

			if _imageDetail[0].Detail != nil && len(*_imageDetail[0].Detail) > 0 {
				_subImageDetail := *_imageDetail[0].Detail
				if _subImageDetail[0].Confidence != nil {
					resultScore = cast.ToInt(*_subImageDetail[0].Confidence * 100)
				}
			}

			resultLabel = _imageDetail[0].Category.Value()

			b, err := json.Marshal(response.Result.ImageDetail)
			if err == nil {
				resultDetail = string(b)
			}
		}
	}

	if resultSuggestion != "" {
		switch strings.ToLower(resultSuggestion) {
		case "pass":
			*result = 0
		case "block":
			*result = 1
		case "review":
			*result = 2
		default:
			*result = 1
		}
	}

	*score = resultScore

	if resultLabel != "" {
		*label = strings.ToLower(resultLabel)
	}

	if resultDetail != "" {
		*detail = resultDetail
	}

	if jobId == nil && resultJobId != "" {
		*jobId = resultJobId
	}

	if len(resultBody) > 0 {
		*body = resultBody
	}

	return nil
}

// GetClient 暴露原生client
func (rc *HuaweiAuditClient) GetClient() interface{} {
	return rc.client
}
