package handler

import (
	"applepages/config"
	"applepages/internal/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct{}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{}
}

type RegisterRequest struct {
	LastName    string `json:"lastName" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	SchoolEmail string `json:"schoolEmail" binding:"required"`
	BackupEmail string `json:"backupEmail" binding:"required"`
}

type RegData struct {
	LastName    string `json:"lastName" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	SchoolEmail string `json:"schoolEmail" binding:"required"`
	BackupEmail string `json:"backupEmail" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// 验证输入是否安全
func validateInput(input string) bool {
	// 检查是否包含可疑的 SQL 关键字
	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "UNION",
		"WHERE", "HAVING", "GROUP BY", "ORDER BY", "--", "/*", "*/",
		"EXEC", "EXECUTE", "DECLARE", "CAST", "CONVERT",
	}

	upperInput := strings.ToUpper(input)
	for _, keyword := range sqlKeywords {
		if strings.Contains(upperInput, keyword) {
			return false
		}
	}

	// 检查是否包含特殊字符
	specialChars := []string{"'", "\"", ";", "\\", "=", "&", "|", "<", ">"}
	for _, char := range specialChars {
		if strings.Contains(input, char) {
			return false
		}
	}

	return true
}

// 验证姓名
func validateName(name string) bool {
	// 允许字母、数字、空格和常见的名字字符
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-']+$`)
	return nameRegex.MatchString(name) && len(name) >= 2 && len(name) <= 50
}

// 验证邮箱
func validateEmail(email string) bool {
	// 严格的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) && len(email) <= 100
}

// HandleRegister 处理用户注册请求
func (h *RegisterHandler) HandleRegister(c *gin.Context) {
	// 获取用户信息
	userInfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
		})
		return
	}

	user, ok := userInfo.(UserInfo)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户信息类型错误",
		})
		return
	}

	// 验证信任等级
	if user.TrustLevel < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "信任等级不足",
		})
		return
	}

	// 解析请求体
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求数据格式错误",
		})
		return
	}

	// 验证输入安全性
	if !validateInput(req.LastName) || !validateInput(req.FirstName) ||
		!validateInput(req.SchoolEmail) || !validateInput(req.BackupEmail) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入包含非法字符",
		})
		return
	}

	// 验证姓名格式
	if !validateName(req.LastName) || !validateName(req.FirstName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "姓名格式不正确",
		})
		return
	}

	// 验证邮箱格式
	if !validateEmail(req.SchoolEmail) || !validateEmail(req.BackupEmail) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱格式不正确",
		})
		return
	}

	// 验证学校邮箱后缀
	if !strings.HasSuffix(strings.ToLower(req.SchoolEmail), "polyu.edu.rs") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "学校邮箱必须以 polyu.edu.rs 结尾",
		})
		return
	}

	// 检查用户是否已注册
	exists, err := database.CheckUserRegistration(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查用户注册状态失败",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "该 Linux Do 账号已注册",
		})
		return
	}

	// 创建注册记录
	registration := &RegData{
		LastName:    strings.TrimSpace(req.LastName),
		FirstName:   strings.TrimSpace(req.FirstName),
		SchoolEmail: strings.ToLower(strings.TrimSpace(req.SchoolEmail)),
		BackupEmail: strings.ToLower(strings.TrimSpace(req.BackupEmail)),
		Password:    generateRandomPassword(),
	}

	if err := CreateWorkspaceUser(registration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册失败",
		})
		return
	}
	database.CreateUserRegistration(&database.UserRegistration{
		LinuxDoUsername:   user.Username,
		LinuxDoTrustLevel: user.TrustLevel,
		LastName:          req.LastName,
		FirstName:         req.FirstName,
		SchoolEmail:       req.SchoolEmail,
		BackupEmail:       req.BackupEmail,
	})
	c.JSON(http.StatusOK, gin.H{
		"message":  "注册成功",
		"password": registration.Password,
	})
}

func getGoogleAccessToken() (string, error) {
	// 构建请求参数
	data := url.Values{}
	data.Set("client_id", config.GlobalConfig.Google.ClientID)
	data.Set("client_secret", config.GlobalConfig.Google.ClientSecret)
	data.Set("refresh_token", config.GlobalConfig.Google.RefreshToken)
	data.Set("grant_type", "refresh_token")

	// 创建请求
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("获取访问令牌失败: %s", string(body))
	}

	// 解析响应
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return result.AccessToken, nil
}

func CreateWorkspaceUser(user *RegData) error {
	log.Printf("开始创建工作区用户: %s %s <%s>", user.FirstName, user.LastName, user.SchoolEmail)

	// 获取访问令牌
	accessToken, err := getGoogleAccessToken()
	if err != nil {
		log.Printf("获取访问令牌失败: %v", err)
		return fmt.Errorf("获取访问令牌失败: %v", err)
	}
	log.Printf("成功获取访问令牌")

	// 构建请求
	reqBody := map[string]interface{}{
		"primaryEmail":  user.SchoolEmail,
		"password":      user.Password,
		"recoveryEmail": user.BackupEmail,
		"name": map[string]string{
			"givenName":  user.FirstName,
			"familyName": user.LastName,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("序列化请求体失败: %v", err)
		return fmt.Errorf("序列化请求体失败: %v", err)
	}
	log.Printf("请求体数据: %s", string(jsonData))

	// 创建请求
	req, err := http.NewRequest("POST", "https://admin.googleapis.com/admin/directory/v1/users", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	log.Printf("发送请求到 Google Admin API")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("创建用户失败: HTTP %d - %s", resp.StatusCode, string(body))
		return fmt.Errorf("创建用户失败: %s", string(body))
	}

	log.Printf("成功创建用户: %s", user.SchoolEmail)
	return nil
}

// 生成随机密码
func generateRandomPassword() string {
	// 定义密码字符集
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	specialChars := "!@#$%^&*"

	// 确保密码至少包含每种字符
	password := make([]byte, 12)
	password[0] = lowerChars[rand.Intn(len(lowerChars))]
	password[1] = upperChars[rand.Intn(len(upperChars))]
	password[2] = numbers[rand.Intn(len(numbers))]
	password[3] = specialChars[rand.Intn(len(specialChars))]

	// 随机填充剩余位置
	allChars := lowerChars + upperChars + numbers + specialChars
	for i := 4; i < 12; i++ {
		password[i] = allChars[rand.Intn(len(allChars))]
	}

	// 打乱密码顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}
