package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// CategoryRepository 分类数据访问层
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository 创建分类数据访问实例
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// FindAll 查询所有分类，按 sort 排序
func (r *CategoryRepository) FindAll() ([]model.Category, error) {
	var cats []model.Category
	err := r.DB.Order("sort ASC").Find(&cats).Error
	if err != nil {
		return nil, err
	}
	return cats, nil
}

// FindByID 根据 ID 查询分类
func (r *CategoryRepository) FindByID(id int64) (*model.Category, error) {
	var cat model.Category
	err := r.DB.Where("id = ?", id).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create 创建分类
func (r *CategoryRepository) Create(cat *model.Category) error {
	return r.DB.Create(cat).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.DB.Save(cat).Error
}

// Delete 物理删除分类（service 层先检查是否被文章引用）
func (r *CategoryRepository) Delete(id int64) error {
	return r.DB.Where("id = ?", id).Delete(&model.Category{}).Error
}

// BuildTree 构建分类树（parent_id=0 为根）
func (r *CategoryRepository) BuildTree(cats []model.Category) []model.Category {
	if len(cats) == 0 {
		return []model.Category{}
	}

	// 按 parent_id 分组
	childrenMap := make(map[int64][]model.Category)
	for _, cat := range cats {
		childrenMap[cat.ParentID] = append(childrenMap[cat.ParentID], cat)
	}

	// 递归构建根节点（parent_id=0）
	tree := buildChildren(childrenMap, 0)
	if tree == nil {
		tree = []model.Category{}
	}
	return tree
}

// buildChildren 递归构建指定父节点下的子树
func buildChildren(childrenMap map[int64][]model.Category, parentID int64) []model.Category {
	children := childrenMap[parentID]
	if len(children) == 0 {
		return nil
	}
	result := make([]model.Category, len(children))
	for i, child := range children {
		result[i] = child
		result[i].Children = buildChildren(childrenMap, child.ID)
	}
	return result
}

// CountArticlesByCategoryID 统计引用该分类的文章数
func (r *CategoryRepository) CountArticlesByCategoryID(categoryID int64) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Article{}).Where("category_id = ? AND is_deleted = 0", categoryID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
