package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess   = 200
	CodeError     = 500
	InvalidParams = 400

	ErrUserNotExist   = 10001
	ErrPasswordWrong  = 10002
	ErrTokenGenFail   = 10003
	ErrUserExist      = 10004
	ErrEncryptFail    = 10005
	ErrCreateFail     = 10006
	ErrTokenInvalid   = 10007
	ErrTokenMalformed = 10008
	

	ErrItemNotFound   = 20001
	ErrItemCreateFail = 20002
	ErrItemUpdateFali = 20003
	ErrItemDeleteFail = 20004
	ErrNoPerMission   = 20005
	ErrDBQueryFail    = 20006

	ErrClaimCreateFail = 30001
	ErrClaimNotFound   = 30002
	ErrClaimUpdateFail = 30003
)

var MsgFlags = map[int]string{
	CodeSuccess:       "操作成功",
	CodeError:         "操作失败",
	InvalidParams:     "请求参数错误",
	ErrUserNotExist:   "用户不存在或角色不匹配",
	ErrPasswordWrong:  "密码错误",
	ErrTokenGenFail:   "生成令牌失败",
	ErrUserExist:      "用户名已存在",
	ErrEncryptFail:    "密码加密失败",
	ErrCreateFail:     "创建账号失败",
	ErrTokenInvalid:   "Token无效",
	ErrTokenMalformed: "Token格式错误",

	ErrItemNotFound:   "物品不存在",
	ErrItemCreateFail: "发布物品失败",
	ErrItemUpdateFali: "更新物品失败",
	ErrItemDeleteFail: "删除物品失败",
	ErrNoPerMission:   "无权限执行此操作",
	ErrDBQueryFail:    "数据库查询失败",

	ErrClaimCreateFail: "申请认领失败",
	ErrClaimNotFound:   "认领记录不存在",
	ErrClaimUpdateFail: "更新认领状态失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[CodeError]
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
	Result(c, http.StatusOK, CodeSuccess, data)
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
