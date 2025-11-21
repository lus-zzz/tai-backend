package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HandlerFunc 定义新的handler函数类型
// 所有handler方法应该返回 (result, nil) 表示成功，或者 (nil, err) 表示失败
type HandlerFunc func(c *gin.Context) (interface{}, error)

// ResponseHandlerMiddleware 响应处理中间件
// 统一处理所有handler的返回值，自动格式化为ResponseBody格式
func ResponseHandlerMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		
		// 获取当前handler
		handler, exists := c.Get("handler")
		if !exists {
			// 如果没有设置handler，直接继续
			c.Next()
			return
		}

		// 类型断言
		handlerFunc, ok := handler.(HandlerFunc)
		if !ok {
			// 如果不是正确的类型，直接继续
			c.Next()
			return
		}

		// 执行handler
		result, err := handlerFunc(c)
		duration := time.Since(startTime).Seconds()
		
		if err != nil {
			// 处理错误
			handleError(c, err, duration)
			return
		}

		// 处理成功响应
		handleSuccess(c, result, duration)
	})
}

// handleSuccess 处理成功响应
func handleSuccess(c *gin.Context, data interface{}, duration float64) {
	response := ResponseBody{
		IsLogin: true,  // 简化处理，默认为true，实际应根据登录状态判断
		Code:    true,
		Message: "",  // 成功时消息可以为空
		Data:    data,
		Time:    duration,
		CodeNum: 200,
		CodeMsg: "success",
	}

	// 统一返回200状态码
	c.JSON(200, response)
}

// handleError 处理错误响应
func handleError(c *gin.Context, err error, duration float64) {
	var apiErr *APIError
	var codeNum int
	var codeMsg, message string

	// 如果已经是APIError，使用其信息
	if err, ok := err.(*APIError); ok {
		apiErr = err
	} else {
		// 否则包装为APIError
		apiErr = WrapError(err, ErrInternalServer)
	}

	// 从映射表获取信息
	mapping := GetErrorMapping(apiErr.ErrorCode)
	codeNum = mapping.CodeNum
	codeMsg = string(apiErr.ErrorCode)
	
	// 获取错误消息，优先使用原始错误信息
	if apiErr.Err != nil {
		message = apiErr.Err.Error()
	} else {
		message = mapping.Message
	}

	response := ResponseBody{
		IsLogin: false,
		Code:    false,
		Message: message,
		Data:    nil,
		Time:    duration,
		CodeNum: codeNum,
		CodeMsg: codeMsg,
	}

	// 统一返回200状态码
	c.JSON(200, response)
}

// WrapHandler 包装handler函数，使其可以被中间件处理
func WrapHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将handler存储到context中
		c.Set("handler", handler)
		c.Next()
	}
}
