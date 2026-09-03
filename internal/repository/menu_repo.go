package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// MenuRepository 菜单数据访问层
type MenuRepository struct {
	DB *gorm.DB
}

// NewMenuRepository 创建菜单数据访问实例
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{DB: db}
}

// FindAll 查询所有启用的菜单
func (r *MenuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.DB.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// FindByIDs 按ID列表查询菜单
func (r *MenuRepository) FindByIDs(ids []int64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.DB.Where("id IN ? AND status = ?", ids, 1).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.DB.Create(menu).Error
}

// FindByID 根据ID查找菜单（不限状态）
func (r *MenuRepository) FindByID(id int64) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.DB.Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id int64) error {
	return r.DB.Where("id = ?", id).Delete(&model.Menu{}).Error
}

// FindByParentID 按父级ID查询子菜单（用于删除前校验子节点）
func (r *MenuRepository) FindByParentID(parentID int64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.DB.Where("parent_id = ?", parentID).Find(&menus).Error
	return menus, err
}

// FindByPath 按路由路径查询菜单（用于path唯一性校验）
func (r *MenuRepository) FindByPath(path string) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.Where("path = ?", path).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// BuildTree 构建菜单树（parent_id=0 为根节点）
// 两段式构树：先建立指针映射并收集父子关系（全程指针操作，
// 避免值拷贝导致子树丢失），所有挂载完成后再递归物化为值切片，
// 与输入顺序无关。
func (r *MenuRepository) BuildTree(menus []model.Menu) []model.Menu {
	menuMap := make(map[int64]*model.Menu, len(menus))
	children := make(map[int64][]*model.Menu, len(menus))

	// 第一遍：建立 id → *Menu 指针映射，重置 Children 避免脏数据累积
	for i := range menus {
		menus[i].Children = nil
		menuMap[menus[i].ID] = &menus[i]
	}

	// 第二遍：通过指针收集父子关系；父不存在的节点提升为根（保持原始顺序）
	var roots []*model.Menu
	for i := range menus {
		node := &menus[i]
		switch {
		case node.ParentID == 0:
			roots = append(roots, node)
		default:
			if _, ok := menuMap[node.ParentID]; ok {
				children[node.ParentID] = append(children[node.ParentID], node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	// 第三遍：所有挂载关系确定后，递归物化为值类型，保证子树完整
	var build func(node *model.Menu) model.Menu
	build = func(node *model.Menu) model.Menu {
		value := *node
		if kids, ok := children[node.ID]; ok {
			value.Children = make([]model.Menu, 0, len(kids))
			for _, kid := range kids {
				value.Children = append(value.Children, build(kid))
			}
		} else {
			value.Children = nil
		}
		return value
	}

	tree := make([]model.Menu, 0, len(roots))
	for _, root := range roots {
		tree = append(tree, build(root))
	}
	return tree
}
