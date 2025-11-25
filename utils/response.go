package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

// 状态码定义
const (
	SUCCESS                  = 200  // 成功
	ERROR                    = 500  // 服务器内部错误
	INVALID_PARAMS           = 400  // 请求参数错误
	UNAUTHORIZED             = 401  // 未授权
	FORBIDDEN                = 403  // 禁止访问
	NOT_FOUND                = 404  // 资源不存在
	CONFLICT                 = 409  // 资源冲突
	TOKEN_EXPIRED            = 1001 // Token过期
	TOKEN_INVALID            = 1002 // Token无效
	PERMISSION_DENIED        = 1003 // 权限不足
	PASSWORD_ERROR           = 1004 // 密码错误
	USER_NOT_EXIST           = 1005 // 用户不存在
	USER_ALREADY_EXIST       = 1006 // 用户已存在

	// 学校相关错误码
	SCHOOL_NOT_EXIST    = 2001 // 学校不存在
	SCHOOL_ALREADY_EXIST = 2002 // 学校已存在
	SCHOOL_SUSPENDED    = 2003 // 学校已被封禁

	// 团队相关错误码
	TEAM_NOT_EXIST          = 3001 // 团队不存在
	TEAM_ALREADY_EXIST      = 3002 // 团队已存在
	TEAM_ALREADY_JOINED     = 3003 // 已经加入团队
	TEAM_NOT_JOINED         = 3004 // 未加入团队
	TEAM_FULL               = 3005 // 团队已满
	TEAM_PASSWORD_ERROR     = 3006 // 团队密码错误
	TEAM_NOT_CAPTAIN        = 3007 // 非队长无法执行此操作
	TEAM_CAPTAIN_CANNOT_LEAVE = 3008 // 队长无法离开团队
)

// 错误信息映射
var codeMsg = map[int]string{
	SUCCESS:                   "操作成功",
	ERROR:                     "服务器内部错误",
	INVALID_PARAMS:            "请求参数错误",
	UNAUTHORIZED:              "未授权，请先登录",
	FORBIDDEN:                 "禁止访问",
	NOT_FOUND:                 "资源不存在",
	CONFLICT:                  "资源冲突",
	TOKEN_EXPIRED:             "Token已过期",
	TOKEN_INVALID:             "Token无效",
	PERMISSION_DENIED:         "权限不足",
	USER_NOT_EXIST:            "用户不存在",
	USER_ALREADY_EXIST:        "用户已存在",
	PASSWORD_ERROR:            "密码错误",
	SCHOOL_NOT_EXIST:          "学校不存在",
	SCHOOL_ALREADY_EXIST:      "学校已存在",
	SCHOOL_SUSPENDED:          "学校已被封禁",
	TEAM_NOT_EXIST:            "团队不存在",
	TEAM_ALREADY_EXIST:        "团队已存在",
	TEAM_ALREADY_JOINED:       "已经加入团队",
	TEAM_NOT_JOINED:           "未加入团队",
	TEAM_FULL:                 "团队已满",
	TEAM_PASSWORD_ERROR:       "团队密码错误",
	TEAM_NOT_CAPTAIN:          "非队长无法执行此操作",
	TEAM_CAPTAIN_CANNOT_LEAVE: "队长无法离开团队",
}

// GetMsg 获取状态码对应的信息
func GetMsg(code int) string {
	msg, ok := codeMsg[code]
	if ok {
		return msg
	}
	return codeMsg[ERROR]
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  GetMsg(SUCCESS),
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  GetMsg(code),
		Data: nil,
	})
}

// ErrorWithMsg 错误响应（自定义消息）
func ErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
