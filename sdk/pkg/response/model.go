package response

type Response struct {
	// 数据集
	RequestId string `protobuf:"bytes,1,opt,name=requestId,proto3" json:"request_id,omitempty"`
	Code      int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code"`
	Message   string `protobuf:"bytes,3,opt,name=msg,proto3" json:"message,omitempty"`
	Status    string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`

	// 兼容性字段
	ErrCode int32 `protobuf:"varint,2,opt,name=errCode,proto3" json:"err_code"`
}

type response struct {
	Response
	Data interface{} `json:"data,omitempty"`
}

type Page struct {
	Count     int `json:"count"`
	PageIndex int `json:"page"`
	PageSize  int `json:"size"`
}

type page struct {
	Page
	List interface{} `json:"list"`
}

func (e *response) SetData(data interface{}) {
	e.Data = data
}

func (e response) Clone() Responses {
	return &e
}

func (e *response) SetTraceID(id string) {
	e.RequestId = id
}

func (e *response) SetMessage(s string) {
	e.Message = s
}

func (e *response) SetCode(code int32) {
	e.Code = code
	e.ErrCode = code
}

func (e *response) SetSuccess(success bool) {
	if !success {
		e.Status = "error"
	}
}
