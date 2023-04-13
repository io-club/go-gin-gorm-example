package api

import (
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTrendRequest struct {
	Name         string                  `form:"name" json:"name" binding:"required"`
	Type         string                  `form:"type" json:"type" binding:"required"`
	Detail       string                  `form:"detail" json:"detail" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func DeleteTrendById(c *gin.Context) {
	id := c.Param("id")
	trendId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, "err")
		return
	}
	err = model.DeleteTrendById(trendId)
	if err != nil {
		c.JSON(404, gin.H{"err": "id 解析失败"})
	}

}

func GetTrendById(c *gin.Context) {
	id := c.Param("id")
	trendId, err := strconv.ParseInt(id, 10, 64)
	trend, err := model.GetTrendById(trendId)
	if err != nil {
		c.JSON(404, gin.H{"err": "id 解析失败"})
		return
	}
	model.GetTrendById(trendId)
	if err != nil {
		c.JSON(404, gin.H{"err": "trend 不存在"})
		return
	}
	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", trend.TableName(), trend.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}
	var ret GetTrendResponse
	ret.Trend = trend
	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetTrendResponse struct {
	model.Trend
	Images []SimpleImageResponse `json:"images"`
}

func CreateTrend(c *gin.Context) {
	var req CreateTrendRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": "图片保存失败"})
		return
	}

	trend := model.Trend{Name: req.Name, Type: req.Type, Detail: req.Detail, ImageURL: filename}
	if err := model.DB.Create(&trend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}

	images := make([]model.Image, len(req.Images))
	for _, file := range req.Images {
		filename := util.CreateFileName(file)
		if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "图片保存失败"})
			return
		}
		images = append(images, model.Image{TableName: model.Trend{}.TableName(), RecordID: trend.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片保存失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "trend": trend})
}

type UpdateTrendRequest struct {
	Name         *string               `form:"name" json:"name" `
	Detail       *string               `form:"detail" json:"detail" `
	Type         *string               `form:"type" json:"type" `
	PreviewImage *multipart.FileHeader `form:"image" json:"image"`
}

func UpdateTrend(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrendRequest
	var old model.Trend
	if err := model.DB.First(&old, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trend not found"})
		return
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "参数解析失败"})
	}

	if req.Name != nil {
		old.Name = *req.Name
	}
	if req.Detail != nil {
		old.Detail = *req.Detail
	}
	if req.Type != nil {
		old.Type = *req.Type
	}
	if req.PreviewImage != nil {
		filename := util.CreateFileName(req.PreviewImage)
		if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "预览图片保存失败"})
			return
		}

		old.ImageURL = filename
	}

	if err := model.DB.Save(&old).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Trend"})
		return
	}

	c.JSON(http.StatusOK, old)
}

type GetTrendsRequest struct {
	model.Pageable
	Name *string `form:"name" json:"name" `
	Type *string `form:"type" json:"type" `
}
type GetTrendsResponse struct {
	model.Trend
	Images []SimpleImageResponse `json:"images"`
}

func GetTrends(c *gin.Context) {
	var req GetTrendsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}

	var trends []model.Trend

	conn := model.DB
	if req.Name != nil {
		conn = conn.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		conn = conn.Where("type = ?", *req.Type)
	}
	if err := conn.Limit(*req.Size).Offset((*req.Page - 1) * *req.Size).Order("id desc").Find(&trends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trends"})
		return
	}

	// 获取图片信息
	imagesMap, err := model.GetImagesByRecordIds(model.Cloth{}.TableName(), model.GetIdsFromTrends(trends))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片信息获取失败"})
		return
	}

	ret := []GetTrendsResponse{}
	// 将图片信息添加到面料信息中
	for _, trend := range trends {
		images := make([]SimpleImageResponse, 0)
		for _, image := range imagesMap[int64(trend.ID)] {
			images = append(images, SimpleImageResponse{int64(image.ID), image.FileName})
		}
		ret = append(ret, GetTrendsResponse{trend, images})
	}

	c.JSON(http.StatusOK, ret)
}
