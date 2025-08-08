package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"file-classifier/internal/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler 注册用户（内存存储 + 加盐哈希）
func RegisterHandler(c *gin.Context) {
	var req registerRequest
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数错误"})
		return
	}
	if len(req.Username) < 3 || len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "用户名/密码过短"})
		return
	}
	if err := auth.RegisterUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LoginHandler 简单登录：开发模式统一密码 admin
func LoginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "请求体错误"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "用户名密码不能为空"})
		return
	}

	// 校验用户
	if err := auth.ValidateUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}

	sid, _ := auth.CreateSession(req.Username)

	// 设置 HttpOnly 会话 Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	c.JSON(http.StatusOK, gin.H{"success": true, "username": req.Username})
}

// MeHandler 返回当前登录用户
func MeHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("sid")
	if err != nil || cookie.Value == "" {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}
	if sess, ok := auth.GetSession(cookie.Value); ok {
		c.JSON(http.StatusOK, gin.H{"authenticated": true, "user": gin.H{"username": sess.Username}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authenticated": false})
}

// LogoutHandler 退出登录
func LogoutHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("sid")
	if err == nil && cookie.Value != "" {
		auth.DeleteSession(cookie.Value)
	}

	// 清除 Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	c.JSON(http.StatusOK, gin.H{"success": true})
}
