package model

import "gorm.io/gorm"

type Cloth struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Detail   string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
	Type     string `gorm:"not null"`
}

func (cloth Cloth) TableName() string {
	return "cloth"
}

func GetIdsFromCloths(cloths []Cloth) []int64 {
	ids := make([]int64, len(cloths))
	for i, cloth := range cloths {
		ids[i] = int64(cloth.ID)
	}
	return ids
}

func DeleteClothById(clothId int64) error {
	var cloth Cloth
	result := DB.Delete(&cloth, clothId)
	return result.Error
}

func GetClothById(id int) (Cloth, error) {
	var cloth Cloth
	result := DB.First(&cloth, id)
	return cloth, result.Error
}

func UpdateCloth(cloth Cloth) (Cloth, error) {
	result := DB.Updates(&cloth)
	return cloth, result.Error
}