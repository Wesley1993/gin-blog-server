package errors

// 通用错误码
const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeTooManyReq   = 429
	CodeServerError  = 500
)

// 业务错误码
const (
	// 用户/认证 10001-10005
	ErrUsernameExists  = 10001
	ErrRoleNotFound    = 10002
	ErrSuperAdminDel   = 10003
	ErrAccountDisabled = 10004
	ErrLoginFailed     = 10005

	// 分类 20001-20002
	ErrCategoryNotFound = 20001
	ErrCategoryInUse    = 20002

	// 文章 30001
	ErrArticleNotFound = 30001

	// OSS 40001-40004
	ErrOSSNotConfigured   = 40001
	ErrOSSUploadFailed    = 40002
	ErrFileTypeNotSupport = 40003
	ErrFileSizeExceed     = 40004

	// ES 50001
	ErrESSearchFallback = 50001

	// 常用网站 60001
	ErrLinkNotFound = 60001

	// 菜单 70001-70003
	ErrMenuNotFound    = 70001
	ErrMenuHasChildren = 70002
	ErrMenuPathExists  = 70003
)

// codeMessages 错误码对应的消息
var codeMessages = map[int]string{
	CodeSuccess:           "ok",
	CodeBadRequest:        "参数错误",
	CodeUnauthorized:      "未登录或token已过期",
	CodeForbidden:         "无权限",
	CodeTooManyReq:        "请求过于频繁",
	CodeServerError:       "服务器内部错误",
	ErrUsernameExists:     "用户名已存在",
	ErrRoleNotFound:       "角色不存在",
	ErrSuperAdminDel:      "超级管理员不可删除",
	ErrAccountDisabled:    "账号已禁用",
	ErrLoginFailed:        "用户名或密码错误",
	ErrCategoryNotFound:   "分类不存在",
	ErrCategoryInUse:      "分类被文章引用，禁止删除",
	ErrArticleNotFound:    "文章不存在",
	ErrOSSNotConfigured:   "OSS配置未设置",
	ErrOSSUploadFailed:    "OSS上传失败",
	ErrFileTypeNotSupport: "文件类型不支持",
	ErrFileSizeExceed:     "文件大小超限",
	ErrESSearchFallback:   "ES搜索降级",
	ErrLinkNotFound:       "常用网站不存在",
	ErrMenuNotFound:       "菜单不存在",
	ErrMenuHasChildren:    "存在子菜单，禁止删除",
	ErrMenuPathExists:     "菜单路由路径已存在",
}

// GetMsg 根据错误码获取对应的消息
func GetMsg(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
