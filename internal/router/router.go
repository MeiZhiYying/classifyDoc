package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"file-classifier/internal/config"
	"file-classifier/internal/handlers"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 允许更大的文件上传
	r.MaxMultipartMemory = config.MaxFileSize

	// CORS中间件
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// API路由 - 必须在静态文件路由之前定义
	api := r.Group("/api")
	{
		api.GET("/stats", handlers.StatsHandler)
		api.GET("/files/:category", handlers.FilesHandler)
		api.GET("/all-files", handlers.AllFilesHandler)
		api.POST("/scan-uploads", handlers.ScanUploadsHandler)
		api.POST("/add-category", handlers.AddCategoryHandler)
		api.GET("/categories", handlers.GetCategoriesHandler)
		api.POST("/delete-category", handlers.DeleteCategoryHandler)
		api.POST("/delete-file", handlers.DeleteFileHandler)
	}

	// 文件上传路由
	r.POST("/upload", handlers.UploadHandler)

	// 文件访问和下载路由
	r.GET("/files/*filepath", handlers.FileHandler)
	r.GET("/download/*filepath", handlers.DownloadHandler)

	// 静态文件服务 - 必须在最后定义
	r.StaticFile("/", config.IndexFile)
	r.Static("/static", config.StaticDir)
	// 新增：上传文件静态服务
	r.Static("/uploads", config.UploadDir)

	return r
}
