package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"applepages/config"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct{}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	TrustLevel int    `json:"trust_level"`
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

// 获取用户信息
func (h *OAuthHandler) getUserInfo(token string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", config.GlobalConfig.OAuth.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取用户信息失败: %s", string(body))
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// HandleCallback 处理OAuth2回调
func (h *OAuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Println("未收到授权码")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未收到授权码",
		})
		return
	}

	log.Printf("收到授权码: %s\n", code)

	// 准备获取token的请求数据
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.GlobalConfig.OAuth.RedirectURI)

	// 发送请求获取token
	req, err := http.NewRequest("POST", config.GlobalConfig.OAuth.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("创建请求失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建请求失败",
		})
		return
	}

	// 添加Basic认证
	auth := base64.StdEncoding.EncodeToString([]byte(config.GlobalConfig.OAuth.ClientID + ":" + config.GlobalConfig.OAuth.ClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("发送token请求到: %s\n", config.GlobalConfig.OAuth.TokenURL)
	log.Printf("Basic认证信息: %s\n", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("请求token失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "请求token失败",
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取响应失败",
		})
		return
	}

	log.Printf("收到响应: HTTP %d\n", resp.StatusCode)
	log.Printf("响应内容: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{
			"error": fmt.Sprintf("OAuth服务器返回错误: %d", resp.StatusCode),
		})
		return
	}

	// 解析响应
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Printf("解析token响应失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "解析token响应失败",
		})
		return
	}

	// 获取用户信息
	userInfo, err := h.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("获取用户信息失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户信息失败",
		})
		return
	}

	log.Printf("成功获取token: %s\n", tokenResp.AccessToken[:10]+"...")
	log.Printf("成功获取用户信息: %+v\n", userInfo)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenResp.AccessToken,
		"type":  tokenResp.TokenType,
		"user":  userInfo,
	})
}
