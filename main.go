package main

import (
	"fibric/api"
	"fibric/model"
	"fibric/mw"
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
		auth.GET("/logout", api.Logout)
	}
	image := base.Group("/image").Use(mw.CheckLoginMiddleware())
	{
		image.POST("/upload", mw.CheckSuperAdmin(), api.UploadImage)
		image.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteImageById)
		// image.DELETE("/", api.DeleteImagesByRecordId)
	}
	fabric := base.Group("/fabric").Use(mw.CheckLoginMiddleware())
	{
		fabric.POST("", mw.CheckSuperAdmin(), api.CreateFabric)
		fabric.GET("/:id", api.GetFabric)
		fabric.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateFabric)
		fabric.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteFabric)
		fabric.GET("/list", api.GetFabrics)
	}
	brand := base.Group("/brand").Use(mw.CheckLoginMiddleware())
	{
		brand.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteBrandById)
		brand.GET("/:id", api.GetBrandById)
		brand.POST("", mw.CheckSuperAdmin(), api.CreateBrand)
		brand.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateBrand)
		brand.GET("/list", api.GetBrands)
	}
	trend := base.Group("/trend").Use(mw.CheckLoginMiddleware())
	{
		trend.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteTrendById)
		trend.GET("/:id", api.GetTrendById)
		trend.POST("", mw.CheckSuperAdmin(), api.CreateTrend)
		trend.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateTrend)
		trend.GET("/list", api.GetTrends)
	}
	cloth := base.Group("/cloth").Use(mw.CheckLoginMiddleware())
	{
		cloth.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteClothById)
		cloth.GET("/:id", api.GetClothById)
		cloth.POST("", mw.CheckSuperAdmin(), api.CreateCloth)
		cloth.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateCloth)
		cloth.GET("/list", api.GetCloths)
	}
	dress := base.Group("/dress").Use(mw.CheckLoginMiddleware())
	{
		dress.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteDressById)
		dress.GET("/:id", api.GetDressById)
		dress.POST("", mw.CheckSuperAdmin(), api.CreateDress)
		dress.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateDress)
		dress.GET("/list", api.GetDresss)
	}
	news := base.Group("/news").Use(mw.CheckLoginMiddleware())
	{
		news.DELETE("/:id", mw.CheckSuperAdmin(), api.DeleteNewsById)
		news.GET("/:id", api.GetNewsById)
		news.POST("", mw.CheckSuperAdmin(), api.CreateNews)
		news.PUT("/:id", mw.CheckSuperAdmin(), api.UpdateNews)
		news.GET("/list", api.GetNewss)
	}
	// 启动 HTTP 服务器
	r.Run(":8080")
}
