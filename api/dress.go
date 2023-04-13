package api

import (
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateDressRequest struct {
	Name         string                  `form:"name" json:"name" binding:"required"`
	Type         string                  `form:"type" json:"type" binding:"required"`
	Detail       string                  `form:"detail" json:"detail" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func DeleteDressById(c *gin.Context) {
	id := c.Param("id")
	dressId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "parseint fail"})
		return
	}
	err = model.DeleteDressById(dressId)
	if err != nil {
		c.JSON(404, gin.H{"error": "delete fail"})
		return
	}
	c.JSON(200, gin.H{"success": "delete success"})
}
func GetDressById(c *gin.Context) {
	id := c.Param("id")
	var dress model.Dress
	if err := model.DB.First(&dress, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dress not found"})
		return
	}
	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", dress.TableName(), dress.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}
	var ret GetDressResponse
	ret.Dress = dress
	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetDressResponse struct {
	model.Dress
	Images []SimpleImageResponse `json:"images"`
}

func CreateDress(c *gin.Context) {
	var req CreateDressRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": "预览图保存失败"})
		return
	}

	dress := model.Dress{Name: req.Name, Detail: req.Detail, ImageURL: filename}
	if err := model.DB.Create(&dress).Error; err != nil {
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
		images = append(images, model.Image{TableName: model.Dress{}.TableName(), RecordID: dress.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "dress": dress})
}

type UpdateDressRequest struct {
	Name         *string               `form:"name" json:"name" `
	Detail       *string               `form:"detail" json:"detail" `
	PreviewImage *multipart.FileHeader `form:"image" json:"image"`
}

func UpdateDress(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDressRequest
	var old model.Dress
	if err := model.DB.First(&old, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dress not found"})
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
	if req.PreviewImage != nil {
		filename := util.CreateFileName(req.PreviewImage)
		if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "预览图片保存失败"})
			return
		}
		old.ImageURL = filename
	}

	if err := model.DB.Save(&old).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Cloth"})
		return
	}

	c.JSON(http.StatusOK, old)
}

type GetDresssRequest struct {
	model.Pageable
}
type GetDresssResponse struct {
	model.Dress
	Images []SimpleImageResponse `json:"images"`
}

func GetDresss(c *gin.Context) {
	var req GetDresssRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}

	var dresss []model.Dress

	conn := model.DB
	if err := conn.Limit(*req.Size).Offset((*req.Page - 1) * *req.Size).Order("id desc").Find(&dresss).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cloths"})
		return
	}

	// 获取图片信息
	imagesMap, err := model.GetImagesByRecordIds(model.Dress{}.TableName(), model.GetIdsFromDresss(dresss))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片信息获取失败"})
		return
	}

	ret := []GetDresssResponse{}
	// 将图片信息添加到面料信息中
	for _, dress := range dresss {
		images := make([]SimpleImageResponse, 0)
		for _, image := range imagesMap[int64(dress.ID)] {
			images = append(images, SimpleImageResponse{int64(image.ID), image.FileName})
		}
		ret = append(ret, GetDresssResponse{dress, images})
	}

	c.JSON(http.StatusOK, ret)
}
