package apiErr

import (
	"easyfind/pkg/response"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

// HandleSysError 处理系统级错误
func HandleSysError(c *gin.Context, errCode int, originErr error) {
	//  记录详细的现场日志
	logger.CtxErrorf(c, "[SysError-%d] Path:%s | Method:%s | IP:%s | Cause: %v",
		errCode,
		c.Request.URL.Path,
		c.Request.Method,
		c.ClientIP(),
		originErr,
	)
	//  向前端返回标准 JSON
	response.Fail(c, errCode)
	//  强制中断
	c.Abort()
}

// HandleBizError 处理业务逻辑错误
func HandleBizError(c *gin.Context, errCode int, msg string) {
	logger.CtxWarnf(c, "[BizError-%d] Path:%s | Msg:%s", errCode, c.Request.URL.Path, msg)

	// 如果传入了自定义 msg 就用自定义的，否则用 response 包里默认的
	if msg != "" {
		response.FailWithMessage(c, errCode, msg)
	} else {
		response.Fail(c, errCode)
	}

	c.Abort()
}

// HandleValidatorError 处理参数校验错误
func HandleValidatorError(c *gin.Context, err error) {
	logger.CtxInfof(c, "[ValidError] Path:%s | Cause: %v", c.Request.URL.Path, err)
	response.FailWithMessage(c, response.InvalidParams, "请求参数格式错误")
	c.Abort()
}
