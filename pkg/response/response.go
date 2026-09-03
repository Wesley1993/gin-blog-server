package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wuzhispace.com/pkg/logger"
)

// Response 统一 API 响应结构
type Response struct {
	Code int    `json:"code"` // 业务错误码
	Msg  string `json:"msg"`  // 提示信息
	Data any    `json:"data"` // 响应数据，无数据时为 null
}

// PageResult 通用泛型分页返回结构
type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// Success 成功响应（无数据）
func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "ok",
	})
}

// SuccessData 成功响应（带数据）
func SuccessData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "ok",
		Data: data,
	})
}

// SuccessPage 成功返回分页数据
func SuccessPage[T any](c *gin.Context, list []T, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "ok",
		Data: PageResult[T]{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 业务错误响应（统一记录一条 Warn 便于排查，userID 由 JWT 中间件写入）
func Fail(c *gin.Context, httpCode int, errCode int, msg string) {
	logger.Warn("请求处理失败",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"userID", c.GetInt64("userID"),
		"httpCode", httpCode,
		"code", errCode,
		"msg", msg,
	)
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  msg,
	})
}
