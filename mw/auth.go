package mw

import (
	"fibric/config"
	"fibric/model"
	"fibric/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if config.MODE == "debug" {
			c.Next()
			return
		}

		session := sessions.Default(c)

		token := session.Get("token")
		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		userId, err := util.CheckToken(token.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 不合法"})
			c.Abort()
			return
		}

		user, err := model.GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func CheckSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		if config.MODE == "debug" {
			c.Next()
			return
		}

		user, _ := c.Get("user")
		if user.(model.User).Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
