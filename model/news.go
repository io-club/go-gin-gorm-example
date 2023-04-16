package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Main     string `gorm:"not null"`
	Type     string `gorm:"not null"`
	ImageURL string ``
}

func (news News) TableName() string {
	return "news"
}

type NewsType string

const (
	// 行业资讯
	NewsTypeIndustry NewsType = "industry"
	// 校企合作
	NewsTypeSchoolCompany NewsType = "school_company"
)

func GetNewsTypeList() []NewsType {
	return []NewsType{NewsTypeIndustry, NewsTypeSchoolCompany}
}
func NewsTypeExists(newsType string) bool {
	for _, t := range GetNewsTypeList() {
		if string(t) == newsType {
			return true
		}
	}
	return false
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
