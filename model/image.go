package model

import "gorm.io/gorm"

// Image 结构体表示图片信息
type Image struct {
	gorm.Model
	TableName string `gorm:"not null"`
	RecordID  uint   `gorm:"not null"`
	FileName  string `gorm:"not null"`
}

func GetImagesById(tableName string, recordID int64) ([]Image, error) {
	var images []Image
	result := DB.Where("table_name = ? AND record_id = ?", tableName, recordID).Find(&images)
	return images, result.Error
}

func GetImagesByRecordIds(tableName string, recordIDs []int64) (map[int64][]Image, error) {
	var images []Image
	result := DB.Where("table_name = ? AND record_id IN (?)", tableName, recordIDs).Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	imagesMap := make(map[int64][]Image)
	for _, image := range images {
		imagesMap[int64(image.RecordID)] = append(imagesMap[int64(image.RecordID)], image)
	}

	return imagesMap, nil
}

func GetImageById(id int64) (Image, error) {
	var image Image
	result := DB.First(&image, id)
	return image, result.Error
}

func CreateImage(image Image) (Image, error) {
	result := DB.Create(&image)
	return image, result.Error
}

func CreateImages(images []Image) error {
	if len(images) == 0 {
		return nil
	}
	err := DB.Create(&images).Error
	return err
}

func DeleteImageById(id int64) error {
	var image Image
	result := DB.Delete(&image, id)
	return result.Error
}

func DeleteImagesByRecordId(tableName string, recordID int64) error {
	var images []Image
	result := DB.Where("table_name = ? AND record_id = ?", tableName, recordID).Delete(&images)
	return result.Error
}

func CountImagesByTableNameAndRecordId(tableName string, recordID int64) (int64, error) {
	var count int64
	result := DB.Model(&Image{}).Where("table_name = ? AND record_id = ?", tableName, recordID).Count(&count)
	return count, result.Error
}
