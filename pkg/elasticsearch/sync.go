package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wuzhispace.com/internal/model"
	"wuzhispace.com/pkg/logger"
)

// SyncArticle upsert 单篇文档到 ES。ES 不可用时静默返回 nil
func (es *ESClient) SyncArticle(article *model.Article) error {
	if !es.Available {
		return nil
	}

	doc := map[string]interface{}{
		"id":          article.ID,
		"title":       article.Title,
		"content":     article.Content,
		"tags":        splitTags(article.Tags),
		"category_id": article.CategoryID,
		"status":      article.Status,
		"create_time": article.CreateTime.Format("2006-01-02 15:04:05"),
	}
	if article.PublishedAt != nil {
		doc["published_at"] = article.PublishedAt.Format("2006-01-02 15:04:05")
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("序列化文档失败: %w", err)
	}

	docID := fmt.Sprintf("%d", article.ID)
	res, err := es.Client.Index(
		es.IndexName,
		bytes.NewReader(body),
		es.Client.Index.WithDocumentID(docID),
		es.Client.Index.WithRefresh("false"),
	)
	if err != nil {
		logger.Error("ES同步文章失败", "articleID", article.ID, "error", err)
		return fmt.Errorf("同步文章到ES失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("ES同步文章返回错误", "articleID", article.ID, "status", res.Status())
		return fmt.Errorf("同步文章到ES返回错误: %s", res.Status())
	}

	return nil
}

// DeleteArticle 从 ES 删除文档
func (es *ESClient) DeleteArticle(articleID int64) error {
	if !es.Available {
		return nil
	}

	docID := fmt.Sprintf("%d", articleID)
	res, err := es.Client.Delete(
		es.IndexName,
		docID,
	)
	if err != nil {
		logger.Error("ES删除文档失败", "articleID", articleID, "error", err)
		return fmt.Errorf("删除ES文档失败: %w", err)
	}
	defer res.Body.Close()

	// 404 也视为成功（文档本就不存在）
	if res.IsError() && res.StatusCode != 404 {
		logger.Error("ES删除文档返回错误", "articleID", articleID, "status", res.Status())
		return fmt.Errorf("删除ES文档返回错误: %s", res.Status())
	}

	return nil
}

// RebuildIndex Bulk API 批量写入文章到 ES
func (es *ESClient) RebuildIndex(articles []model.Article) (success, failed int, err error) {
	if !es.Available {
		return 0, 0, fmt.Errorf("ES不可用")
	}

	var buf bytes.Buffer
	for _, article := range articles {
		// action line
		meta := fmt.Sprintf(`{"index":{"_index":"%s","_id":"%d"}}`, es.IndexName, article.ID)
		buf.WriteString(meta)
		buf.WriteString("\n")

		// document line
		doc := map[string]interface{}{
			"id":          article.ID,
			"title":       article.Title,
			"content":     article.Content,
			"tags":        splitTags(article.Tags),
			"category_id": article.CategoryID,
			"status":      article.Status,
			"create_time": article.CreateTime.Format("2006-01-02 15:04:05"),
		}
		if article.PublishedAt != nil {
			doc["published_at"] = article.PublishedAt.Format("2006-01-02 15:04:05")
		}
		docBytes, marshalErr := json.Marshal(doc)
		if marshalErr != nil {
			failed++
			continue
		}
		buf.Write(docBytes)
		buf.WriteString("\n")
	}

	res, err := es.Client.Bulk(
		bytes.NewReader(buf.Bytes()),
		es.Client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("bulk写入ES失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, 0, fmt.Errorf("bulk写入ES返回错误: %s", res.Status())
	}

	// 解析 bulk 响应统计成功/失败
	var bulkResp map[string]interface{}
	if decodeErr := json.NewDecoder(res.Body).Decode(&bulkResp); decodeErr != nil {
		return len(articles), 0, nil // 写入成功但解析响应失败，乐观返回
	}

	if items, ok := bulkResp["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if indexOp, ok := itemMap["index"].(map[string]interface{}); ok {
					if errObj, ok := indexOp["error"]; ok && errObj != nil {
						failed++
						logger.Error("ES bulk写入单条失败", "error", errObj)
					} else {
						success++
					}
				}
			}
		}
	} else {
		success = len(articles)
	}

	return success, failed, nil
}

// splitTags 将逗号分隔的标签字符串拆分为切片
func splitTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	parts := strings.Split(tags, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
