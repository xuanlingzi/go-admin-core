package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xuanlingzi/go-admin-core/sdk/pkg"
)

var Default = &response{}

// Error 失败数据处理
func Error(c *gin.Context, code int, err error, message ...string) {
	res := Default.Clone()
	if err != nil {
		res.SetMessage(err.Error())
	}
	if len(message) > 0 {
		res.SetMessage(strings.Join(message, ","))
	}
	res.SetTraceID(pkg.GenerateMsgIDFromContext(c))
	res.SetCode(int32(code))
	res.SetSuccess(false)
	c.Set("result", res)
	c.Set("status", code)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, message ...string) {
	res := Default.Clone()
	res.SetData(data)
	res.SetSuccess(true)
	if len(message) > 0 {
		res.SetMessage(strings.Join(message, ","))
	}
	res.SetTraceID(pkg.GenerateMsgIDFromContext(c))
	res.SetCode(http.StatusOK)
	c.Set("result", res)
	c.Set("status", http.StatusOK)
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// PageOK 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, message ...string) {
	var res page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, message...)
}

// Custom 兼容函数
func Custom(c *gin.Context, data gin.H) {
	data["requestId"] = pkg.GenerateMsgIDFromContext(c)
	c.Set("result", data)
	c.AbortWithStatusJSON(http.StatusOK, data)
}
