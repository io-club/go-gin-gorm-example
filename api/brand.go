package api

import (
	"fibric/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
