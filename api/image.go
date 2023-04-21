package api

import (
	"fibric/config"
	"fibric/model"
	"fibric/util"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteImageById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = model.DeleteImageById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func DeleteImagesByRecordId(c *gin.Context) {
	tableName := c.Param("tableName")
	recordId, err := strconv.ParseInt(c.Param("recordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = model.DeleteImagesByRecordId(tableName, recordId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

type UploadImageRequest struct {
	TableName string                  `form:"tableName" json:"tableName" binding:"required"`
	RecordId  int64                   `form:"recordId" json:"recordId" binding:"required"`
	Images    []*multipart.FileHeader `form:"images" json:"images"`
}

func UploadImage(c *gin.Context) {
	// 获取请求中的参数
	var req UploadImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check table exists
	if !model.TableExists(req.TableName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "表不存在"})
		return
	}

	// count images in db
	count, err := model.CountImagesByTableNameAndRecordId(req.TableName, req.RecordId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count+int64(len(req.Images)) > config.MaxImageNum {
		c.JSON(http.StatusBadRequest, gin.H{"error": "最多上传 20 张图片", "your": count, "has": len(req.Images)})
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
		images = append(images, model.Image{TableName: req.TableName, RecordID: uint(req.RecordId), FileName: filename})
	}

	// 将图片信息保存到数据库中
	if err := model.CreateImages(images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type SimpleImageResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
