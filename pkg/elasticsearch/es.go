package elasticsearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"wuzhispace.com/config"
	"wuzhispace.com/pkg/logger"
)

// ESClient 封装 ES 客户端
type ESClient struct {
	Client    *elasticsearch.Client
	IndexName string
	Available bool // ES 是否可用（启动检查后设置）
}

// NewESClient 创建 ES 客户端
func NewESClient(cfg *config.ESConfig) (*ESClient, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("创建ES客户端失败: %w", err)
	}

	indexName := cfg.IndexPrefix + "_article_index"

	es := &ESClient{
		Client:    client,
		IndexName: indexName,
		Available: false,
	}

	return es, nil
}

// CheckAndInit 启动时检查 ES 连通性、IK 分词器、索引存在性
func (es *ESClient) CheckAndInit(ctx context.Context) {
	// 1. Ping ES
	res, err := es.Client.Info()
	if err != nil {
		logger.Warn("ES连接失败，搜索将降级到PG", "error", err)
		es.Available = false
		return
	}
	res.Body.Close()

	// 2. 检查 IK 分词器
	if !es.checkIKAnalyzer(ctx) {
		logger.Warn("ES IK分词器未安装，搜索将降级到PG")
		es.Available = false
		return
	}

	// 3. 确保索引存在
	if err := es.EnsureIndex(ctx); err != nil {
		logger.Error("ES创建索引失败", "error", err)
		es.Available = false
		return
	}

	es.Available = true
	logger.Info("ES初始化成功")
}

// checkIKAnalyzer 验证 IK 分词器
func (es *ESClient) checkIKAnalyzer(ctx context.Context) bool {
	body := strings.NewReader(`{"analyzer":"ik_max_word","text":"测试中文分词"}`)
	res, err := es.Client.API.Indices.Analyze(
		es.Client.API.Indices.Analyze.WithBody(body),
	)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return !res.IsError()
}

// EnsureIndex 幂等创建索引
func (es *ESClient) EnsureIndex(ctx context.Context) error {
	// 检查索引是否存在
	res, err := es.Client.Indices.Exists([]string{es.IndexName})
	if err != nil {
		return fmt.Errorf("检查索引存在性失败: %w", err)
	}

	if res.StatusCode == 200 {
		return nil // 索引已存在
	}

	// 创建索引
	mapping := strings.NewReader(ArticleIndexMapping)
	createRes, err := es.Client.Indices.Create(
		es.IndexName,
		es.Client.Indices.Create.WithBody(mapping),
	)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("创建索引返回错误: %s", createRes.String())
	}

	return nil
}
