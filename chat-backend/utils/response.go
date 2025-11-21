package utils

import (
	"time"
)

// ResponseBody 统一响应体结构
// swagger:model
type ResponseBody struct {
	// 登录状态
	// required: true
	IsLogin bool `json:"is_login"`
	// 请求是否成功
	// required: true
	Code bool `json:"code"`
	// 响应消息
	// required: true
	Message string `json:"msg"`
	// 响应数据
	// required: true
	Data interface{} `json:"data"`
	// 持续时间（秒）
	// required: true
	Time float64 `json:"time"`
	// 数字错误代码
	// required: true
	CodeNum int `json:"code_num"`
	// 错误代码消息
	// required: true
	CodeMsg string `json:"code_msg"`
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}
