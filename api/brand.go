package api

import (
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateBrandRequest struct {
	Name         string                  `form:"name" json:"name" binding:"required"`
	Detail       string                  `form:"detail" json:"detail" binding:"required"`
	PreviewImage *multipart.FileHeader   `form:"image" json:"image" binding:"required"`
	Images       []*multipart.FileHeader `form:"images" json:"images"`
}

func DeleteBrandById(c *gin.Context) {
	id := c.Param("id")
	brandId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, "err")
		return
	}
	err = model.DeleteBrandById(brandId)
	if err != nil {
		c.JSON(404, "err")
	}
}
func GetBrandById(c *gin.Context) {
	id := c.Param("id")
	var brand model.Brand
	if err := model.DB.First(&brand, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "brand not found"})
		return
	}
	var images []model.Image
	if err := model.DB.Where("table_name = ? AND record_id = ?", brand.TableName(), brand.ID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get images"})
		return
	}
	var ret GetBrandResponse
	ret.Brand = brand
	simpleImages := make([]SimpleImageResponse, len(images))
	for i, image := range images {
		simpleImages[i] = SimpleImageResponse{ID: int64(image.ID), Name: image.FileName}
	}
	ret.Images = simpleImages

	c.JSON(http.StatusOK, ret)
}

type GetBrandResponse struct {
	model.Brand
	Images []SimpleImageResponse `json:"images"`
}

func CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := util.CreateFileName(req.PreviewImage)
	if err := c.SaveUploadedFile(req.PreviewImage, "images/"+filename); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	brand := model.Brand{Name: req.Name, Detail: req.Detail, ImageURL: filename}
	if err := model.DB.Create(&brand).Error; err != nil {
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
		images = append(images, model.Image{TableName: model.Brand{}.TableName(), RecordID: brand.ID, FileName: filename})
	}
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "brand": brand})
}
