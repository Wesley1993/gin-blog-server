package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"wuzhispace.com/pkg/logger"
)

// SearchResult 搜索结果单条
type SearchResult struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Tags       string `json:"tags"`
	CategoryID int64  `json:"category_id"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

// SearchResponse 搜索结果
type SearchResponse struct {
	List  []SearchResult `json:"list"`
	Total int64          `json:"total"`
}

// Search 全文搜索 — multi_match best_fields，按 create_time 降序
func (es *ESClient) Search(keyword string, page, pageSize int) (*SearchResponse, error) {
	return es.search(keyword, page, pageSize, false)
}

// SearchPublished 全文搜索仅已发布文章（博客前台公开接口使用）
func (es *ESClient) SearchPublished(keyword string, page, pageSize int) (*SearchResponse, error) {
	return es.search(keyword, page, pageSize, true)
}

// search 内部搜索实现，publishedOnly 为 true 时追加 status=1 过滤
func (es *ESClient) search(keyword string, page, pageSize int, publishedOnly bool) (*SearchResponse, error) {
	if !es.Available {
		return nil, fmt.Errorf("ES不可用")
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	from := (page - 1) * pageSize

	// 构造查询体：仅已发布时用 bool 包裹并追加 status 过滤
	var queryObj map[string]interface{}
	multiMatch := map[string]interface{}{
		"query":  keyword,
		"type":   "best_fields",
		"fields": []string{"title^3", "content", "tags^2"},
	}
	if publishedOnly {
		queryObj = map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{map[string]interface{}{"multi_match": multiMatch}},
				"filter": []interface{}{map[string]interface{}{"term": map[string]interface{}{"status": 1}}},
			},
		}
	} else {
		queryObj = map[string]interface{}{"multi_match": multiMatch}
	}

	query := map[string]interface{}{
		"query": queryObj,
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title":   map[string]interface{}{},
				"content": map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
		"sort": []map[string]interface{}{
			{"create_time": map[string]interface{}{"order": "desc"}},
		},
		"from": from,
		"size": pageSize,
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("序列化搜索请求失败: %w", err)
	}

	res, err := es.Client.Search(
		es.Client.Search.WithIndex(es.IndexName),
		es.Client.Search.WithBody(bytes.NewReader(body)),
		es.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		logger.Warn("ES搜索降级到PG", "keyword", keyword, "error", err)
		return nil, fmt.Errorf("ES搜索请求失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("ES搜索返回错误", "response", res.String())
		return nil, fmt.Errorf("ES搜索返回错误: %s", res.Status())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析ES搜索响应失败: %w", err)
	}

	resp := &SearchResponse{
		List: make([]SearchResult, 0),
	}

	// 解析 total
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if totalObj, ok := hits["total"]; ok {
			switch t := totalObj.(type) {
			case map[string]interface{}:
				if val, ok := t["value"].(float64); ok {
					resp.Total = int64(val)
				}
			case float64:
				resp.Total = int64(t)
			}
		}

		if hitsArr, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArr {
				hitMap, ok := hit.(map[string]interface{})
				if !ok {
					continue
				}

				sr := SearchResult{}

				// 解析 _source
				if source, ok := hitMap["_source"].(map[string]interface{}); ok {
					if id, ok := source["id"].(float64); ok {
						sr.ID = int64(id)
					}
					if title, ok := source["title"].(string); ok {
						sr.Title = title
					}
					if tags, ok := source["tags"].([]interface{}); ok {
						tagStrs := make([]string, 0, len(tags))
						for _, tag := range tags {
							if s, ok := tag.(string); ok {
								tagStrs = append(tagStrs, s)
							}
						}
						sr.Tags = joinTags(tagStrs)
					}
					if catID, ok := source["category_id"].(float64); ok {
						sr.CategoryID = int64(catID)
					}
					if status, ok := source["status"].(float64); ok {
						sr.Status = int(status)
					}
					if ct, ok := source["create_time"].(string); ok {
						sr.CreateTime = ct
					}
				}

				// 高亮覆盖标题
				if highlight, ok := hitMap["highlight"].(map[string]interface{}); ok {
					if titleHL, ok := highlight["title"].([]interface{}); ok && len(titleHL) > 0 {
						if s, ok := titleHL[0].(string); ok {
							sr.Title = s
						}
					}
				}

				resp.List = append(resp.List, sr)
			}
		}
	}

	return resp, nil
}

// joinTags 将标签切片拼接为逗号分隔字符串
func joinTags(tags []string) string {
	result := ""
	for i, tag := range tags {
		if i > 0 {
			result += ","
		}
		result += tag
	}
	return result
}
