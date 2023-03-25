package model

import "gorm.io/gorm"

// Fabric 结构体表示服装面料信息
type Brand struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Detail   string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
}

func (brand Brand) TableName() string {
	return "brand"
}

func GetIdsFromBrands(brands []Brand) []int64 {
	ids := make([]int64, len(brands))
	for i, brand := range brands {
		ids[i] = int64(brand.ID)
	}
	return ids
}
func GetBrandById(id int) (Brand, error) {
	var brand Brand
	result := DB.First(&brand, id)
	return brand, result.Error
}
func DeleteBrandById(brandId int64) error {
	var brand Brand
	result := DB.Delete(&brand, brandId)
	return result.Error
}
