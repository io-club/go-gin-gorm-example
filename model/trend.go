package model

import "gorm.io/gorm"

type Trend struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Type     string `gorm:"not null"`
	Detail   string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
}

func (trend Trend) TableName() string {
	return "trend"
}

func GetTrendById(id int64) (Trend, error) {
	var trend Trend
	result := DB.First(&trend, id)
	return trend, result.Error
}

func DeleteTrendById(trendId int64) error {
	var trend Trend
	result := DB.Delete(&trend, trendId)
	return result.Error
}
