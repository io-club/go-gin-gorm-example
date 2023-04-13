package api

import (
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateNewsRequest struct {
	Main         string                  `form:"name" json:"Main" binding:"required"`
	Type         string                  `form:"type" json:"type" binding:"required"`
	Title        string                  `form:"detail" json:"title" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func DeleteNewsById(c *gin.Context) {
	id := c.Param("id")
	newsId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "parseint fail"})
		return
	}
	err = model.DeleteClothById(newsId)
	if err != nil {
		c.JSON(404, gin.H{"error": "delete fail"})
		return
	}
	c.JSON(200, gin.H{"success": "delete success"})
}
func GetNewsById(c *gin.Context) {
	id := c.Param("id")
	var news model.News
	if err := model.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}
	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", news.TableName(), news.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}
	var ret GetNewsResponse
	ret.News = news
	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetNewsResponse struct {
	model.News
	Images []SimpleImageResponse `json:"images"`
}

type GetNewssRequest struct {
	model.Pageable
}
type GetNewssResponse struct {
	model.News
	Images []SimpleImageResponse `json:"images"`
}

func GetNewss(c *gin.Context) {
	var req GetNewssRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}

	var newss []model.News

	conn := model.DB
	if err := conn.Limit(*req.Size).Offset((*req.Page - 1) * *req.Size).Order("id desc").Find(&newss).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get newss"})
		return
	}

	// 获取图片信息
	imagesMap, err := model.GetImagesByRecordIds(model.News{}.TableName(), model.GetIdsFromNewss(newss))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片信息获取失败"})
		return
	}

	ret := []GetNewssResponse{}
	// 将图片信息添加到面料信息中
	for _, news := range newss {
		images := make([]SimpleImageResponse, 0)
		for _, image := range imagesMap[int64(news.ID)] {
			images = append(images, SimpleImageResponse{int64(image.ID), image.FileName})
		}
		ret = append(ret, GetNewssResponse{news, images})
	}

	c.JSON(http.StatusOK, ret)
}
func CreateNews(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": "预览图保存失败"})
		return
	}

	news := model.News{Main: req.Main, Title: req.Title}
	if err := model.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件创建失败"})
		return
	}

	images := make([]model.Image, len(req.Images))
	for _, file := range req.Images {
		filename := util.CreateFileName(file)
		if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "图片保存失败"})
			return
		}
		images = append(images, model.Image{TableName: model.News{}.TableName(), RecordID: news.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "news": news})
}

type UpdateNewsRequest struct {
	Main  *string `form:"main" json:"main" `
	Title *string `form:"title" json:"title" `
}

func UpdateNews(c *gin.Context) {
	id := c.Param("id")
	var req UpdateNewsRequest
	var old model.News
	if err := model.DB.First(&old, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": " not found"})
		return
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "参数解析失败"})
	}

	if req.Title != nil {
		old.Title = *req.Title
	}
	if req.Main != nil {
		old.Main = *req.Main
	}

	if err := model.DB.Save(&old).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Cloth"})
		return
	}

	c.JSON(http.StatusOK, old)
}
