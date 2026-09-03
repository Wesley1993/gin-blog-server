package service

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"gorm.io/gorm"

	"wuzhispace.com/config"
	"wuzhispace.com/pkg/elasticsearch"
	"wuzhispace.com/pkg/logger"
	pkgredis "wuzhispace.com/pkg/redis"
)

// healthCheckTimeout 依赖健康探测的超时时间
const healthCheckTimeout = 2 * time.Second

// SystemService 系统运行信息服务
type SystemService struct {
	startTime   time.Time
	db          *gorm.DB
	redisClient *pkgredis.RedisClient
	esClient    *elasticsearch.ESClient

	// healthMu/lastHealth 缓存各依赖上次健康状态，仅在状态翻转时告警，避免仪表盘轮询刷屏
	healthMu   sync.Mutex
	lastHealth map[string]string
}

// NewSystemService 创建系统服务实例
func NewSystemService(
	startTime time.Time,
	db *gorm.DB,
	redisClient *pkgredis.RedisClient,
	esClient *elasticsearch.ESClient,
) *SystemService {
	return &SystemService{
		startTime:   startTime,
		db:          db,
		redisClient: redisClient,
		esClient:    esClient,
		lastHealth:  make(map[string]string),
	}
}

// reportHealth 记录依赖健康状态，仅在健康→异常翻转时打 Warn，恢复时打 Info
func (s *SystemService) reportHealth(name, status string) {
	if status != "healthy" && status != "unhealthy" {
		return
	}
	s.healthMu.Lock()
	defer s.healthMu.Unlock()
	last := s.lastHealth[name]
	if last == status {
		return
	}
	s.lastHealth[name] = status
	if status == "unhealthy" {
		logger.Warn("依赖健康探测异常", "dependency", name, "lastStatus", last)
	} else if last == "unhealthy" {
		logger.Info("依赖健康探测恢复", "dependency", name)
	}
}

// SystemOverview 系统运行概览
type SystemOverview struct {
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	AppMode     string `json:"app_mode"`
	GoVersion   string `json:"go_version"`
	Goroutines  int    `json:"goroutines"`
	MemAlloc    uint64 `json:"mem_alloc_mb"`
	Uptime      string `json:"uptime"`
	DBStatus    string `json:"db_status"`
	RedisStatus string `json:"redis_status"`
	ESStatus    string `json:"es_status"`
}

// GetOverview 采集当前系统运行概览（含各依赖健康状态）
func (s *SystemService) GetOverview() *SystemOverview {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	overview := &SystemOverview{
		AppName:     config.GlobalConfig.App.Name,
		AppVersion:  config.GlobalConfig.App.Version,
		AppMode:     config.GlobalConfig.App.Mode,
		GoVersion:   runtime.Version(),
		Goroutines:  runtime.NumGoroutine(),
		MemAlloc:    memStats.Alloc / 1024 / 1024,
		Uptime:      formatDuration(time.Since(s.startTime)),
		DBStatus:    s.checkDB(),
		RedisStatus: s.checkRedis(),
		ESStatus:    s.checkES(),
	}

	return overview
}

// checkDB 通过底层连接 Ping 检测 PostgreSQL 健康状态
func (s *SystemService) checkDB() string {
	status := s.probeDB()
	s.reportHealth("db", status)
	return status
}

func (s *SystemService) probeDB() string {
	if s.db == nil {
		return "not_configured"
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return "unhealthy"
	}
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// checkRedis 通过 Ping 检测 Redis 健康状态
func (s *SystemService) checkRedis() string {
	status := s.probeRedis()
	s.reportHealth("redis", status)
	return status
}

func (s *SystemService) probeRedis() string {
	if s.redisClient == nil {
		return "not_configured"
	}
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()
	if err := s.redisClient.Ping(ctx); err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// checkES 通过 Info 接口检测 Elasticsearch 健康状态
func (s *SystemService) checkES() string {
	status := s.probeES()
	s.reportHealth("es", status)
	return status
}

func (s *SystemService) probeES() string {
	if s.esClient == nil {
		return "not_configured"
	}
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()
	res, err := s.esClient.Client.Info(s.esClient.Client.Info.WithContext(ctx))
	if err != nil {
		return "unhealthy"
	}
	defer res.Body.Close()
	if res.IsError() {
		return "unhealthy"
	}
	return "healthy"
}

// formatDuration 将时长格式化为 "1d 2h 3m" 形式
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
