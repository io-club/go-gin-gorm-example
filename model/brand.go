package model

import "gorm.io/gorm"

// Brand 结构体表示品牌信息
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
func UpdateBrand(brand Brand) (Brand, error) {
	result := DB.Updates(&brand)
	return brand, result.Error
}
