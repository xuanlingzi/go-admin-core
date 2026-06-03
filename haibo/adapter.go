package haibo

// AdapterHaibo 海博开放平台出站适配器接口。
//
// 由 HaiboClient 实现，经 sdk.Runtime 全局注册，供 retail-core 业务编排层
// 统一调用。接口保持最小化：组装/签名/POST 的细节封装在实现内部。
type AdapterHaibo interface {
	// String 返回适配器标识（appId），用于注册时作为 key。
	String() string
	// Invoke 组装系统级参数、序列化 body、签名并发起出站调用。
	Invoke(path string, bizParams interface{}) (*HaiboResponse, error)
	// SetCallLogger 注入出站调用日志钩子（在 DI 装配处接线到 HaiboMessageLog）。
	SetCallLogger(l CallLogger)
}
