package elasticsearch

// ArticleIndexMapping 文章索引 mapping
const ArticleIndexMapping = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "id":          {"type": "long"},
      "title":       {"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart"},
      "content":     {"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart"},
      "tags":        {"type": "keyword"},
      "category_id": {"type": "long"},
      "status":      {"type": "integer"},
      "create_time": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss||epoch_millis"},
      "published_at": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss||epoch_millis"}
    }
  }
}`
