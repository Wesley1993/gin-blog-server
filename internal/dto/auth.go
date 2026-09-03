package dto

import "wuzhispace.com/internal/model"

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token string `json:"token"`
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
	User        UserDTO      `json:"user"`
	Menus       []model.Menu `json:"menus"`
	Permissions []string     `json:"permissions"`
}

// UserDTO 用户信息DTO
type UserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	RoleID   int64  `json:"role_id"`
}
