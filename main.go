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
		image.POST("/upload", api.UploadImage).Use(mw.CheckSuperAdmin())
		image.DELETE("/:id", api.DeleteImageById).Use(mw.CheckSuperAdmin())
		// image.DELETE("/", api.DeleteImagesByRecordId)
	}
	fabric := base.Group("/fabric").Use(mw.CheckLoginMiddleware())
	{
		fabric.POST("", api.CreateFabric).Use(mw.CheckSuperAdmin())
		fabric.GET("/:id", api.GetFabric)
		fabric.PUT("/:id", api.UpdateFabric).Use(mw.CheckSuperAdmin())
		fabric.DELETE("/:id", api.DeleteFabric).Use(mw.CheckSuperAdmin())
		fabric.GET("/list", api.GetFabrics)
	}
	brand := base.Group("/brand").Use(mw.CheckLoginMiddleware())
	{
		brand.DELETE("/:id", api.DeleteBrandById).Use(mw.CheckSuperAdmin())
		brand.GET("/:id", api.GetBrandById)
		brand.POST("", api.CreateBrand).Use(mw.CheckSuperAdmin())
		brand.PUT("/:id", api.UpdateBrand).Use(mw.CheckSuperAdmin())
		brand.GET("/list", api.GetBrands)
	}
	trend := base.Group("/trend").Use(mw.CheckLoginMiddleware())
	{
		trend.DELETE("/:id", api.DeleteTrendById).Use(mw.CheckSuperAdmin())
		trend.GET("/:id", api.GetTrendById)
		trend.POST("", api.CreateTrend).Use(mw.CheckSuperAdmin())
		trend.PUT("/:id", api.UpdateTrend).Use(mw.CheckSuperAdmin())
		trend.GET("/list", api.GetTrends)
	}
	cloth := base.Group("/cloth").Use(mw.CheckLoginMiddleware())
	{
		cloth.DELETE("/:id", api.DeleteClothById).Use(mw.CheckSuperAdmin())
		cloth.GET("/:id", api.GetClothById)
		cloth.POST("", api.CreateCloth).Use(mw.CheckSuperAdmin())
		cloth.PUT("/:id", api.UpdateCloth).Use(mw.CheckSuperAdmin())
		cloth.GET("/list", api.GetCloths)
	}
	dress := base.Group("/dress").Use(mw.CheckLoginMiddleware())
	{
		dress.DELETE("/:id", api.DeleteDressById).Use(mw.CheckSuperAdmin())
		dress.GET("/:id", api.GetDressById)
		dress.POST("", api.CreateDress).Use(mw.CheckSuperAdmin())
		dress.PUT("/:id", api.UpdateDress).Use(mw.CheckSuperAdmin())
		dress.GET("/list", api.GetDresss)
	}
	news := base.Group("/news").Use(mw.CheckLoginMiddleware())
	{
		news.DELETE("/:id", api.DeleteNewsById).Use(mw.CheckSuperAdmin())
		news.GET("/:id", api.GetNewsById)
		news.POST("", api.CreateNews).Use(mw.CheckSuperAdmin())
		news.PUT("/:id", api.UpdateNews).Use(mw.CheckSuperAdmin())
		news.GET("/list", api.GetNewss)
	}
	// 启动 HTTP 服务器
	r.Run(":8080")
}
