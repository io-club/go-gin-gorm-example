package main

import (
	"fibric/api"
	"fibric/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 HTTP 服务器
	r := gin.Default()

	// 使用 cookie 作为 session 存储方式
	store := cookie.NewStore([]byte(model.SessionSecret))
	r.Use(sessions.Sessions("session", store))

	// 注册路由
	base := r.Group("/api")
	base.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})
	auth := base.Group("/auth")
	{
		auth.POST("/login", api.Login)
		auth.POST("/register", api.Register)
	}
	image := base.Group("/image")
	{
		// TODO: curl add
		image.DELETE("/:id", api.DeleteImageById)
		// TODO: curl add
		image.DELETE("/", api.DeleteImagesByRecordId)
	}
	fabric := base.Group("/fabric")
	{
		// TODO: curl modify
		fabric.POST("/", api.CreateFabric)
		// TODO: curl modify
		fabric.GET("/:id", api.GetFabric)
		// TODO: curl modify
		fabric.PUT("/:id", api.UpdateFabric)
		// TODO: curl modify
		fabric.DELETE("/:id", api.DeleteFabric)
		// TODO: curl modify
		fabric.GET("/list", api.GetFabrics)
	}

	// 启动 HTTP 服务器
	r.Run(":8080")
}
