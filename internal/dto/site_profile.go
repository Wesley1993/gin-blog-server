package dto

// SiteProfileContact 站长联系方式
type SiteProfileContact struct {
	Github  string `json:"github"`  // GitHub 地址
	Email   string `json:"email"`   // 邮箱
	Wechat  string `json:"wechat"`  // 微信号
	QQ      string `json:"qq"`      // QQ 号
	Address string `json:"address"` // 地址
}

// SiteProfileSkill 技能项
type SiteProfileSkill struct {
	Name  string `json:"name"`  // 技能名称
	Level int    `json:"level"` // 熟练度 0-100
	Group string `json:"group"` // 分组（如 后端/前端/工程化）
}

// SiteProfileProject 项目经历项
type SiteProfileProject struct {
	Name  string   `json:"name"`  // 项目名称
	Desc  string   `json:"desc"`  // 项目描述
	Tech  []string `json:"tech"`  // 技术栈
	Start string   `json:"start"` // 开始时间 如 2026-01
	End   string   `json:"end"`   // 结束时间，空表示至今
	Link  string   `json:"link"`  // 项目链接
}

// SiteProfileResp 站长个人资料响应（公开接口与管理接口共用）
type SiteProfileResp struct {
	Name     string                `json:"name"`     // 站长姓名
	Avatar   string                `json:"avatar"`   // 头像
	Title    string                `json:"title"`    // 头衔/一句话介绍
	Bio      string                `json:"bio"`      // 个人简历/自序
	Contacts SiteProfileContact    `json:"contacts"` // 联系方式
	Skills   []*SiteProfileSkill   `json:"skills"`   // 技能列表
	Projects []*SiteProfileProject `json:"projects"` // 项目经历列表
}

// SaveSiteProfileReq 保存站长个人资料请求（整体覆盖保存）
type SaveSiteProfileReq struct {
	Name     string                `json:"name"`
	Avatar   string                `json:"avatar"`
	Title    string                `json:"title"`
	Bio      string                `json:"bio"`
	Contacts SiteProfileContact    `json:"contacts"`
	Skills   []*SiteProfileSkill   `json:"skills"`
	Projects []*SiteProfileProject `json:"projects"`
}
