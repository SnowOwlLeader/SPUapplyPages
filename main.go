package main

import (
	"fmt"
	"log"

	"applepages/config"
	"applepages/internal/database"
	"applepages/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("配置初始化失败:", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 预检请求结果缓存12小时
	}))

	// 创建处理器
	oauthHandler := handler.NewOAuthHandler()
	userHandler := handler.NewUserHandler()
	registerHandler := handler.NewRegisterHandler()

	// 设置静态文件服务
	r.Static("/assets", "./www/school/SPUapplyPages/assets")
	r.StaticFile("/", "./www/school/SPUapplyPages/web/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./www/school/SPUapplyPages/web/index.html")
	})

	// 设置API路由
	api := r.Group("/api")
	{
		oauth := api.Group("/oauth")
		{
			oauth.GET("/callback", oauthHandler.HandleCallback)
		}

		// 需要认证的路由组
		authenticated := api.Group("")
		authenticated.Use(userHandler.GetUserInfo)
		{
			authenticated.GET("/user/info", func(c *gin.Context) {}) // 实际的处理在中间件中
			authenticated.POST("/register", registerHandler.HandleRegister)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%s", config.GlobalConfig.Server.Port)
	log.Printf("服务器启动在 %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
