package util

import "gorm.io/gorm"

// PageQuery 通用分页请求参数
type PageQuery struct {
	Page      int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	SortField string `form:"sort_field" json:"sort_field"`
	SortOrder string `form:"sort_order" json:"sort_order"` // asc / desc
}

// GetOffset 计算偏移量
func (p *PageQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取每页条数
func (p *PageQuery) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// Paginate 泛型分页查询方法
// T: 数据模型类型
// query: 基础查询条件 (如 db.Where("status = ?", 1))
// pageQuery: 分页参数
// list: 接收结果的切片指针 (如 &[]model.Product{})
func Paginate[T any](db *gorm.DB, query *gorm.DB, pageQuery *PageQuery, list *[]T) (int64, error) {
	// 1. 设置默认分页参数
	if pageQuery.Page <= 0 {
		pageQuery.Page = 1
	}
	if pageQuery.PageSize <= 0 {
		pageQuery.PageSize = 10
	}

	// 2. 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	// 3. 查询列表数据
	if err := query.
		Offset(pageQuery.GetOffset()).
		Limit(pageQuery.PageSize).
		Order("id DESC"). // 默认按 ID 倒序，可在外部覆盖
		Find(list).Error; err != nil {
		return 0, err
	}

	return total, nil
}
