package api

import (
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateClothRequest struct {
	Name         string                  `form:"name" json:"name" binding:"required"`
	Type         string                  `form:"type" json:"type" binding:"required"`
	Detail       string                  `form:"detail" json:"detail" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func DeleteClothById(c *gin.Context) {
	id := c.Param("id")
	clothId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "parseint fail"})
		return
	}
	err = model.DeleteClothById(clothId)
	if err != nil {
		c.JSON(404, gin.H{"error": "delete fail"})
		return
	}
	c.JSON(200, gin.H{"success": "delete success"})
}
func GetClothById(c *gin.Context) {
	id := c.Param("id")
	var cloth model.Cloth
	if err := model.DB.First(&cloth, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloth not found"})
		return
	}
	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", cloth.TableName(), cloth.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}
	var ret GetClothResponse
	ret.Cloth = cloth
	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetClothResponse struct {
	model.Cloth
	Images []SimpleImageResponse `json:"images"`
}

func CreateCloth(c *gin.Context) {
	var req CreateClothRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": "预览图保存失败"})
		return
	}

	cloth := model.Cloth{Name: req.Name, Detail: req.Detail, Type: req.Type, ImageURL: filename}
	if err := model.DB.Create(&cloth).Error; err != nil {
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
		images = append(images, model.Image{TableName: model.Cloth{}.TableName(), RecordID: cloth.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "cloth": cloth})
}

type UpdateClothRequest struct {
	Name         *string               `form:"name" json:"name" `
	Detail       *string               `form:"detail" json:"detail" `
	Type         *string               `form:"type" json:"type"`
	PreviewImage *multipart.FileHeader `form:"image" json:"image"`
}

func UpdateCloth(c *gin.Context) {
	id := c.Param("id")
	var req UpdateClothRequest
	var old model.Cloth
	if err := model.DB.First(&old, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cloth not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Cloth"})
		return
	}

	c.JSON(http.StatusOK, old)
}

type GetClothsRequest struct {
	model.Pageable
	Name *string `form:"name" json:"name"`
	Type *string `form:"type" json:"type"`
}
type GetClothsResponse struct {
	model.Cloth
	Images []SimpleImageResponse `json:"images"`
}

func GetCloths(c *gin.Context) {
	var req GetClothsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}

	var cloths []model.Cloth

	conn := model.DB
	if req.Name != nil {
		conn = conn.Where("name like ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		conn = conn.Where("type = ?", *req.Type)
	}
	if err := conn.Limit(*req.Size).Offset((*req.Page - 1) * *req.Size).Order("id desc").Find(&cloths).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cloths"})
		return
	}

	// 获取图片信息
	imagesMap, err := model.GetImagesByRecordIds(model.Cloth{}.TableName(), model.GetIdsFromCloths(cloths))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片信息获取失败"})
		return
	}

	ret := []GetClothsResponse{}
	// 将图片信息添加到面料信息中
	for _, cloth := range cloths {
		images := make([]SimpleImageResponse, 0)
		for _, image := range imagesMap[int64(cloth.ID)] {
			images = append(images, SimpleImageResponse{int64(image.ID), image.FileName})
		}
		ret = append(ret, GetClothsResponse{cloth, images})
	}

	c.JSON(http.StatusOK, ret)
}
