package repository

import (
	"fmt"

	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// RoleRepository 角色数据访问层
type RoleRepository struct {
	DB *gorm.DB
}

// NewRoleRepository 创建角色数据访问实例
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

// FindByID 根据ID查找角色
func (r *RoleRepository) FindByID(id int64) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindAll 查找所有角色
func (r *RoleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.Order("id ASC").Find(&roles).Error
	return roles, err
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.Role) error {
	return r.DB.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.Role) error {
	return r.DB.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id int64) error {
	return r.DB.Where("id = ?", id).Delete(&model.Role{}).Error
}

// FindByMenuID 查询 menu_ids JSONB 包含指定菜单ID的所有角色
func (r *RoleRepository) FindByMenuID(menuID int64) ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.Where("menu_ids @> ?::jsonb", fmt.Sprintf("[%d]", menuID)).Find(&roles).Error
	return roles, err
}
