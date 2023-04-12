package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Title string `gorm:"not null"`
	Main  string `gorm:"not null"`
	Type  string `gorm:"not null"`
}

func (news News) TableName() string {
	return "news"
}

func GetIdsFromNewss(newss []News) []int64 {
	ids := make([]int64, len(newss))
	for i, news := range newss {
		ids[i] = int64(news.ID)
	}
	return ids
}

func DeleteNewsById(newsId int64) error {
	var news News
	result := DB.Delete(&news, newsId)
	return result.Error
}

func GetNewsById(id int) (News, error) {
	var news News
	result := DB.First(&news, id)
	return news, result.Error
}
func UpdateNews(news News) (News, error) {
	result := DB.Updates(&news)
	return news, result.Error
}
