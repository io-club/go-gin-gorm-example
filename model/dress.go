package model

import "gorm.io/gorm"

type Dress struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Detail   string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
	Type     string `gorm:"not null"`
}

func (dress Dress) TableName() string {
	return "dress"
}

func GetIdsFromDresss(dresss []Dress) []int64 {
	ids := make([]int64, len(dresss))
	for i, dress := range dresss {
		ids[i] = int64(dress.ID)
	}
	return ids
}

func DeleteDressById(dressId int64) error {
	var dress Dress
	result := DB.Delete(&dress, dressId)
	return result.Error
}

func GetDressById(id int) (Dress, error) {
	var dress Dress
	result := DB.First(&dress, id)
	return dress, result.Error
}

func UpdateDress(dress Dress) (Dress, error) {
	result := DB.Updates(&dress)
	return dress, result.Error
}
