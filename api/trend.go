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
		c.JSON(404, "err")
	}

}

func GetTrendById(c *gin.Context) {
	id := c.Param("id")
	trendId, err := strconv.ParseInt(id, 10, 64)
	trend, err :=model.GetTrendById(trendId)
	if err != nil {
		c.JSON(404,"未查询到此id")
		return
	}
	model.GetTrendById(trendId)
	if err != nil {
		c.JSON(400,"err:转换失败")
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	trend := model.Trend{Name: req.Name, Type: req.Type, Detail: req.Detail, ImageURL: filename}
	if err := model.DB.Create(&trend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	images := make([]model.Image, len(req.Images))
	for _, file := range req.Images {
		filename := util.CreateFileName(file)
		if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, model.Image{TableName: model.Trend{}.TableName(), RecordID: trend.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "trend": trend})
}
