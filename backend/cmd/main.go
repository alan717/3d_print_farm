package main

import (
	"log"
	"os"

	"github.com/3d-printing-farm/backend/internal/config"
	"github.com/3d-printing-farm/backend/internal/handlers"
	"github.com/3d-printing-farm/backend/internal/middleware"
	"github.com/3d-printing-farm/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 初始化配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 设置Gin模式为release模式以减少debug信息
	gin.SetMode(gin.ReleaseMode)

	// 初始化Gin路由
	r := gin.Default()

	// 配置信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	// 添加CORS中间件
	r.Use(middleware.CORS())

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	fileHandler := handlers.NewFileHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	materialHandler := handlers.NewMaterialHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// API路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 用户路由
		user := api.Group("/user")
		user.Use(middleware.AuthRequired())
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		// 文件路由
		files := api.Group("/files")
		{
			files.POST("/upload", middleware.AuthRequired(), fileHandler.Upload)
			files.GET("/:id", fileHandler.GetFile) // 移除认证要求以允许3D预览
			files.DELETE("/:id", middleware.AuthRequired(), fileHandler.DeleteFile)
		}

		// 材料路由 - 移除末尾斜杠避免重定向
		api.GET("/materials", materialHandler.GetMaterials)

		// 订单路由
		orders := api.Group("/orders")
		orders.Use(middleware.AuthRequired())
		{
			orders.POST("/", orderHandler.CreateOrder)
			orders.GET("/", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(), handlers.AdminMiddleware(db))
		{
			// 仪表板
			admin.GET("/dashboard", adminHandler.GetDashboardStats)

			// 用户管理
			admin.GET("/users", adminHandler.GetUsers)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)

			// 订单管理
			admin.GET("/orders", adminHandler.GetOrders)
			admin.PUT("/orders/:id/status", adminHandler.UpdateOrderStatus)

			// 材料管理
			admin.POST("/materials", adminHandler.CreateMaterial)
			admin.PUT("/materials/:id", adminHandler.UpdateMaterial)
			admin.DELETE("/materials/:id", adminHandler.DeleteMaterial)

			// 系统设置
			admin.GET("/settings", adminHandler.GetSystemSettings)
			admin.PUT("/settings/:key", adminHandler.UpdateSystemSetting)

			// 操作日志
			admin.GET("/logs", adminHandler.GetAdminLogs)
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}