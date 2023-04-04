package api

import (
	"fibric/model"
	"fibric/util"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateFabricRequest struct {
	Name         string                  `form:"name" json:"name" binding:"required"`
	Detail       string                  `form:"detail" json:"detail" binding:"required"`
	Category     string                  `form:"category" json:"category" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func CreateFabric(c *gin.Context) {
	// 获取请求中的参数
	var req CreateFabricRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 图片不能超过 5 个
	if len(req.Images) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "最多上传 5 张图片"})
		return
	}

	// 将上传的图片保存到本地
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将面料信息保存到数据库中
	fabric := model.Fabric{Name: req.Name, Detail: req.Detail, Category: req.Category, ImageURL: filename}
	if err := model.DB.Create(&fabric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	images := make([]model.Image, len(req.Images))
	// 将上传的图片保存到本地
	for _, file := range req.Images {
		filename := util.CreateFileName(file)
		if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, model.Image{TableName: model.Fabric{}.TableName(), RecordID: fabric.ID, FileName: filename})
	}

	// 将图片信息保存到数据库中
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "fabric": fabric})
}

func GetFabric(c *gin.Context) {
	id := c.Param("id")

	var fabric model.Fabric

	if err := model.DB.First(&fabric, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fabric not found"})
		return
	}

	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", fabric.TableName(), fabric.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}

	var ret GetFabricsResponse
	ret.Fabric = fabric

	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetFabricsRequest struct {
	model.Pageable
	Category string `form:"category" json:"category"`
}
type GetFabricsResponse struct {
	model.Fabric
	Images []SimpleImageResponse `json:"images"`
}

func GetFabrics(c *gin.Context) {
	var req GetFabricsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fabrics []model.Fabric

	conn := model.DB
	fmt.Printf("category: [%s]\n", req.Category)
	if req.Category != "" {
		conn = conn.Where("category = ?", req.Category)
	}
	if err := conn.Limit(*req.Size).Offset((*req.Page - 1) * *req.Size).Order("id desc").Find(&fabrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get fabrics"})
		return
	}

	// 获取图片信息
	imagesMap, err := model.GetImagesByRecordIds(model.Fabric{}.TableName(), model.GetIdsFromFabrics(fabrics))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ret []GetFabricsResponse
	// 将图片信息添加到面料信息中
	for _, fabric := range fabrics {
		images := make([]SimpleImageResponse, 0)
		for _, image := range imagesMap[int64(fabric.ID)] {
			images = append(images, SimpleImageResponse{int64(image.ID), image.FileName})
		}
		ret = append(ret, GetFabricsResponse{fabric, images})
	}

	c.JSON(http.StatusOK, ret)
}

type UpdateFabricRequest struct {
	Name         *string               `form:"name" json:"name" `
	Category     *string               `form:"category" json:"category" `
	Detail       *string               `form:"detail" json:"detail" `
	PreviewImage *multipart.FileHeader `form:"image" json:"image"`
}

func UpdateFabric(c *gin.Context) {
	id := c.Param("id")
	var req UpdateFabricRequest
	var old model.Fabric
	if err := model.DB.First(&old, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fabric not found"})
		return
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "参数解析失败"})
	}

	if req.Name != nil {
		old.Name = *req.Name
	}
	if req.Category != nil {
		old.Category = *req.Category
	}
	if req.Detail != nil {
		old.Name = *req.Detail
	}
	if req.PreviewImage != nil {
		filename := util.CreateFileName(req.PreviewImage)
		if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		old.Name = filename
	}

	if err := model.DB.Save(&old).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update fabric"})
		return
	}

	c.JSON(http.StatusOK, req)
}

func DeleteFabric(c *gin.Context) {
	id := c.Param("id")

	var fabric model.Fabric

	if err := model.DB.First(&fabric, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Cloth not found"})
		return
	}

	if err := model.DB.Delete(&fabric).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cloth"})
		return
	}

	// 删除图片
	if err := model.DeleteImagesByRecordId(model.Fabric{}.TableName(), int64(fabric.ID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}
