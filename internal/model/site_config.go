package model

// SiteConfig 站点配置模型
type SiteConfig struct {
	BaseModel
	SiteName     string `gorm:"column:site_name;type:varchar(100)" json:"site_name"`
	SiteDesc     string `gorm:"column:site_desc;type:varchar(255)" json:"site_desc"`
	Copyright    string `gorm:"column:copyright;type:varchar(255)" json:"copyright"`
	FoundedAt    string `gorm:"column:founded_at;type:varchar(20)" json:"founded_at"`  // 建站日期 如 2024-01-01
	Icp          string `gorm:"column:icp;type:varchar(100)" json:"icp"`               // ICP备案号 如 京ICP备XXXXXXXX号
	PoliceIcp    string `gorm:"column:police_icp;type:varchar(100)" json:"police_icp"` // 公安备案号 如 京公网安备 11010502000000号
	OssAccessKey string `gorm:"column:oss_access_key;type:varchar(200)" json:"oss_access_key"`
	OssSecretKey string `gorm:"column:oss_secret_key;type:varchar(200)" json:"oss_secret_key"`
	OssBucket    string `gorm:"column:oss_bucket;type:varchar(100)" json:"oss_bucket"`
	OssEndpoint  string `gorm:"column:oss_endpoint;type:varchar(200)" json:"oss_endpoint"`
	OssDomain    string `gorm:"column:oss_domain;type:varchar(200)" json:"oss_domain"`
	OssProvider  string `json:"oss_provider" gorm:"column:oss_provider;default:aliyun"`
	OssRegion    string `json:"oss_region" gorm:"column:oss_region"`
	OssInsecure  int    `json:"oss_insecure" gorm:"column:oss_insecure;type:smallint;default:0"` // 跳过HTTPS证书校验（自签名证书/私有端点）0否 1是
}

// TableName 指定表名
func (SiteConfig) TableName() string {
	return "blog_site_config"
}
