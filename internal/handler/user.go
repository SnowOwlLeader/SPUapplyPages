package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"applepages/config"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从请求头获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未提供认证信息",
		})
		return
	}

	// 解析token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "认证格式错误",
		})
		return
	}

	tokenType := parts[0]
	token := parts[1]

	// 请求Linux Do API获取用户信息
	req, err := http.NewRequest("GET", config.GlobalConfig.OAuth.UserInfoURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建请求失败",
		})
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "请求用户信息失败",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取响应失败",
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{
			"error": fmt.Sprintf("获取用户信息失败: %s", string(body)),
		})
		return
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "解析用户信息失败",
		})
		return
	}

	// 将用户信息存储到上下文中
	c.Set("user", userInfo)

	// 如果是API调用，返回用户信息
	if c.Request.URL.Path == "/api/user/info" {
		c.JSON(http.StatusOK, gin.H{
			"user": userInfo,
		})
		return
	}

	// 如果是中间件调用，继续下一个处理器
	c.Next()
}
