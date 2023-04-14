package api

import (
	"fibric/model"
	"fibric/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 定义登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
func Login(c *gin.Context) {
	var form LoginRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询用户
	var user model.User
	result := model.DB.Where("username = ?", form.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}
	// 生成 jwt
	token, err := util.GenerateToken(user.ID)
	// 保存登录状态到 session
	session := sessions.Default(c)
	session.Set("token", token)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

// 定义注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	var form struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断用户名是否已经存在
	var user model.User
	result := model.DB.Where("username = ?", form.Username).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	// 创建用户
	user = model.User{Username: form.Username, Password: string(hashedPassword)}
	model.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}
