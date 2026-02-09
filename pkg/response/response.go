package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	// 业务错误码 (可根据需要扩展)
	ERROR_USER_NOT_EXIST  = 10001
	ERROR_PASSWORD_WRONG  = 10002
	ERROR_TOKEN_GEN_FAIL  = 10003
	ERROR_USER_EXIST      = 10004
	ERROR_ENCRYPT_FAIL    = 10005
	ERROR_CREATE_FAIL     = 10006
	ERROR_TOKEN_INVALID   = 10007
	ERROR_TOKEN_MALFORMED = 10008
)

var MsgFlags = map[int]string{
	SUCCESS:               "操作成功",
	ERROR:                 "操作失败",
	INVALID_PARAMS:        "请求参数错误",
	ERROR_USER_NOT_EXIST:  "用户不存在或角色不匹配",
	ERROR_PASSWORD_WRONG:  "密码错误",
	ERROR_TOKEN_GEN_FAIL:  "生成令牌失败",
	ERROR_USER_EXIST:      "用户名已存在",
	ERROR_ENCRYPT_FAIL:    "密码加密失败",
	ERROR_CREATE_FAIL:     "创建账号失败",
	ERROR_TOKEN_INVALID:   "Token无效",
	ERROR_TOKEN_MALFORMED: "Token格式错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"` // 使用 data 统一包裹数据
}

func Result(c *gin.Context, httpCode int, errCode int, data interface{}) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, SUCCESS, data)
}

// Fail 失败响应
func Fail(c *gin.Context, errCode int) {
	Result(c, http.StatusOK, errCode, nil)
}

// FailWithMessage 自定义消息失败响应
func FailWithMessage(c *gin.Context, errCode int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  msg,
		Data: nil,
	})
}
