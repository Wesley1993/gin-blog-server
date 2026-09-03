package repository

import (
	"gorm.io/gorm"
	"wuzhispace.com/internal/model"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByRoleID 查找某角色下所有用户
func (r *UserRepository) FindByRoleID(roleID int64) ([]model.User, error) {
	var users []model.User
	err := r.DB.Where("role_id = ?", roleID).Find(&users).Error
	return users, err
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// Delete 物理删除用户
func (r *UserRepository) Delete(id int64) error {
	return r.DB.Where("id = ?", id).Delete(&model.User{}).Error
}

// Page 分页查询用户
func (r *UserRepository) Page(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.DB.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
