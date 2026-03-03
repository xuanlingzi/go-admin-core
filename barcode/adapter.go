package barcode

// AdapterBarcode 条码查询适配器接口
// QueryByBarcode 返回通用 map，key 由各实现自行约定，调用方按需断言。
type AdapterBarcode interface {
	String() string
	Setup() error
	Close()
	// QueryByBarcode 根据条形码查询商品信息
	QueryByBarcode(barcode string) (map[string]interface{}, error)
}
