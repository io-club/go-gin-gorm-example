package api

import (
	"fibric/model"
	"fibric/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateFabric(c *gin.Context) {
	// 获取请求中的参数
	name := c.PostForm("name")
	details := c.PostForm("detail")

	// 获取上传的图片
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将面料信息保存到数据库中
	fabric := model.Fabric{Name: name, Detail: details}
	if err := model.DB.Create(&fabric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将上传的图片保存到本地
	filename := util.CreateFileName(file)
	if err := c.SaveUploadedFile(file, "images/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将图片信息保存到数据库中
	image := model.Image{TableName: model.Fabric{}.TableName(), RecordID: fabric.ID, FileName: filename}
	if err := model.DB.Create(&image).Error; err != nil {
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

	c.JSON(http.StatusOK, fabric)
}

type GetFabricsResponse struct {
	model.Fabric
	Images []model.Image `json:"images"`
}

func GetFabrics(c *gin.Context) {
	var fabrics []model.Fabric

	if err := model.DB.Find(&fabrics).Error; err != nil {
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
		ret = append(ret, GetFabricsResponse{fabric, imagesMap[int64(fabric.ID)]})
	}

	c.JSON(http.StatusOK, ret)
}

type UpdateFabricRequest struct {
	model.Fabric
	ImageURL string `json:"image_url"`
}

func UpdateFabric(c *gin.Context) {
	id := c.Param("id")
	var fabric UpdateFabricRequest

	if err := model.DB.First(&fabric, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fabric not found"})
		return
	}

	// 获取请求中的参数
	newName := fabric.Name
	if name := c.PostForm("name"); name != "" {
		newName = name
	}
	newDetail := fabric.Detail
	if detail := c.PostForm("detail"); detail != "" {
		newDetail = detail
	}

	// 获取上传的图片
	newFilename := fabric.ImageURL
	file, err := c.FormFile("image")
	if err == nil {
		// 将上传的图片保存到本地
		filename := util.CreateFileName(file)
		if err := c.SaveUploadedFile(file, "images/"+filename); err == nil {
			newFilename = filename
		}
	}

	// 将面料信息保存到数据库中
	fabric.Name = newName
	fabric.Detail = newDetail
	fabric.ImageURL = newFilename
	if err := model.DB.Save(&fabric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update fabric"})
		return
	}

	c.JSON(http.StatusOK, fabric)
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
