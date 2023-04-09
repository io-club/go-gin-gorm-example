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
		image.POST("/upload", api.UploadImage)
		image.DELETE("/:id", api.DeleteImageById)
		// image.DELETE("/", api.DeleteImagesByRecordId)
	}
	fabric := base.Group("/fabric")
	{
		fabric.POST("", api.CreateFabric)
		fabric.GET("/:id", api.GetFabric)
		fabric.PUT("/:id", api.UpdateFabric)
		fabric.DELETE("/:id", api.DeleteFabric)
		fabric.GET("/list", api.GetFabrics)
	}
	brand := base.Group("/brand")
	{
		brand.DELETE("/:id", api.DeleteBrandById)
		brand.GET("/:id", api.GetBrandById)
		brand.POST("", api.CreateBrand)
		brand.PUT("/:id", api.UpdateBrand)
		brand.GET("/list", api.GetBrands)
	}
	trend := base.Group("/trend")
	{
		trend.DELETE("/:id", api.DeleteTrendById)
		trend.GET("/:id", api.GetTrendById)
		trend.POST("", api.CreateTrend)
		trend.PUT("/:id", api.UpdateTrend)
		trend.GET("/list", api.GetTrends)
	}
	cloth := base.Group("/cloth")
	{
		cloth.DELETE("/:id", api.DeleteClothById)
		cloth.GET("/:id", api.GetClothById)
		cloth.POST("", api.CreateCloth)
		cloth.PUT("/:id", api.UpdateCloth)
		cloth.GET("/list", api.GetCloths)
	}
	dress := base.Group("/dress")
	{
		dress.DELETE("/:id", api.DeleteDressById)
		dress.GET("/:id", api.GetDressById)
		dress.POST("", api.CreateDress)
		dress.PUT("/:id", api.UpdateDress)
		dress.GET("/list", api.GetDresss)
	}
	news := base.Group("/news")
	{
		news.DELETE("/:id", api.DeleteNewsById)
		news.GET("/:id", api.GetNewsById)
		news.POST("", api.CreateNews)
		news.PUT("/:id", api.UpdateNews)
		news.GET("/list", api.GetNewss)
	}
	// 启动 HTTP 服务器
	r.Run(":8080")
}
