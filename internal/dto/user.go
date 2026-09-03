package dto

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	RoleID   int64  `json:"role_id" binding:"required"`
	Status   int16  `json:"status"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	ID       int64  `json:"id" binding:"required"`
	Nickname string `json:"nickname"`
	RoleID   int64  `json:"role_id"`
	Status   int16  `json:"status"`
}

// ResetPwdReq 重置密码请求
type ResetPwdReq struct {
	Password string `json:"password" binding:"required"`
}

// UpdateUserStatusReq 用户启停用请求（独立于编辑，便于按 user:status 权限控制）
type UpdateUserStatusReq struct {
	Status int16 `json:"status"`
}
