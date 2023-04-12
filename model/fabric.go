package model

import "gorm.io/gorm"

// Fabric 结构体表示服装面料信息
type Fabric struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Category string `gorm:"not null"`
	Detail   string `gorm:"not null"`
	ImageURL string `gorm:"not null"`
}

func (fabric Fabric) TableName() string {
	return "fabric"
}

func GetIdsFromFabrics(fabrics []Fabric) []int64 {
	ids := make([]int64, len(fabrics))
	for i, fabric := range fabrics {
		ids[i] = int64(fabric.ID)
	}
	return ids
}
