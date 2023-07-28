package moderation

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	moderation "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/region"
)

var _huaweiAudit *moderation.ModerationClient

func GetHuaweiAuditClient() *moderation.ModerationClient {
	return _huaweiAudit
}

type HuaweiAuditClient struct {
	client      *moderation.ModerationClient
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
			panic(err)
		}
	}

	r := &HuaweiAuditClient{
		client:      client,
		callbackUrl: callbackUrl,
	}

	return r
}

func (rc *HuaweiAuditClient) String() string {
	return "huawei_audit"
}

func (rc *HuaweiAuditClient) Check() bool {
	return rc.client != nil
}

func (rc *HuaweiAuditClient) Close() {
}

func (rc *HuaweiAuditClient) AuditText(content string, suggestion *string, label *string, detail *string) error {

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

	suggestion = response.Result.Suggestion
	label = response.Result.Label
	if response.Result.Details != nil {
		b, err := json.Marshal(response.Result.Details)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *HuaweiAuditClient) AuditImage(url string, suggestion *string, label *string, detail *string) error {

	var imageCategories = []string{
		"porn",
		"politics",
		"abuse",
	}
	imageReq := model.ImageDetectionReq{
		Url:        &url,
		Categories: imageCategories,
		EventType:  "head_image",
	}
	request := model.CheckImageModerationRequest{}
	request.Body = &imageReq
	response, err := rc.client.CheckImageModeration(&request)
	if err != nil {
		return fmt.Errorf("huawei Audit error: %v", err.Error())
	}

	suggestion = response.Result.Suggestion
	label = response.Result.Category
	if response.Result.Details != nil {
		b, err := json.Marshal(response.Result.Details)
		if err == nil {
			*detail = string(b)
		}
	}

	return nil
}

func (rc *HuaweiAuditClient) AuditVideo(url string, frame int32) error {

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
	_, err := rc.client.RunCreateVideoModerationJob(&request)
	if err != nil {
		return fmt.Errorf("huawei Audit error: %v", err.Error())
	}
	return nil
}

// GetClient 暴露原生client
func (rc *HuaweiAuditClient) GetClient() interface{} {
	return rc.client
}
