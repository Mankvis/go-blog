package category

import (
	"go-blog/pkg/logger"
	"go-blog/pkg/model"
	"go-blog/pkg/types"
)

// Create 创建文章分类，通过 category.ID 来判断是否创建成功
func (category *Category) Create() (err error) {
	if err := model.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Get 通过 ID 获取分类
func Get(idstr string) (Category, error) {
	var category Category
	id := types.StringToUint64(idstr)
	if err := model.DB.Find(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}

// All 获取分类数据
func All() ([]Category, error) {
	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}
