package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 文件信息结构
type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Type string `json:"type"` // "filename", "ai", "failed"
}

// 分类统计结构
type CategoryStats struct {
	Count int        `json:"count"`
	Files []FileInfo `json:"files"`
}

// 上传结果结构
type UploadResult struct {
	Total               int                      `json:"total"`
	Processed           int                      `json:"processed"`
	FirstStepClassified int                      `json:"firstStepClassified"`
	AIClassified        int                      `json:"aiClassified"`
	Classifications     map[string]CategoryStats `json:"classifications"`
}

// 响应结构
type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Results *UploadResult `json:"results,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// 全局变量
var classificationStats = map[string]CategoryStats{
	"合同":   {Count: 0, Files: []FileInfo{}},
	"简历":   {Count: 0, Files: []FileInfo{}},
	"发票":   {Count: 0, Files: []FileInfo{}},
	"论文":   {Count: 0, Files: []FileInfo{}},
	"未分类":  {Count: 0, Files: []FileInfo{}},
	"新增分类": {Count: 0, Files: []FileInfo{}},
}

// 分类关键词配置
var classificationKeywords = map[string][]string{
	"合同": {"合同", "协议", "契约", "contract", "agreement", "合作", "签署"},
	"简历": {"简历", "履历", "resume", "cv", "个人简历", "求职", "应聘"},
	"发票": {"发票", "票据", "invoice", "收据", "账单", "bill", "费用"},
	"论文": {"论文", "研究", "paper", "thesis", "学术", "期刊", "研究报告", "毕业论文"},
}

// 确保上传目录存在
func ensureUploadDir() {
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
}

// 根据文件名进行分类
func classifyByFilename(filename string) string {
	lowerName := strings.ToLower(filename)

	for category, keywords := range classificationKeywords {
		for _, keyword := range keywords {
			if strings.Contains(lowerName, strings.ToLower(keyword)) {
				return category
			}
		}
	}
	return "未分类"
}

// AI分析占位符函数
func classifyByAI(filePath, filename string) string {
	// 模拟AI分析延迟
	time.Sleep(1 * time.Second)

	// 模拟AI分析结果（随机返回一个分类）
	categories := []string{"合同", "简历", "发票", "论文"}
	randomCategory := categories[rand.Intn(len(categories))]

	log.Printf("AI分析结果: %s -> %s", filename, randomCategory)
	return randomCategory
}

// 重置分类统计
func resetClassificationStats() {
	for key := range classificationStats {
		classificationStats[key] = CategoryStats{Count: 0, Files: []FileInfo{}}
	}
}

// 添加文件到分类
func addFileToCategory(category string, fileInfo FileInfo) {
	if stats, exists := classificationStats[category]; exists {
		stats.Files = append(stats.Files, fileInfo)
		stats.Count = len(stats.Files)
		classificationStats[category] = stats
	}
}

// 文件上传处理
func uploadHandler(c *gin.Context) {
	// 解析多文件上传
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "文件上传解析失败: " + err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "没有上传文件",
		})
		return
	}

	if len(files) > 200 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "最多支持上传200个文件",
		})
		return
	}

	// 重置统计
	resetClassificationStats()

	results := &UploadResult{
		Total:               len(files),
		Processed:           0,
		FirstStepClassified: 0,
		AIClassified:        0,
		Classifications:     classificationStats,
	}

	log.Printf("开始处理 %d 个文件", len(files))

	// 第一步：根据文件名分类
	var unclassifiedFiles []FileInfo

	for _, file := range files {
		filename := file.Filename
		category := classifyByFilename(filename)

		// 保存文件
		savePath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("保存文件失败: %s, %v", filename, err)
			continue
		}

		fileInfo := FileInfo{
			Name: filename,
			Path: savePath,
			Size: file.Size,
		}

		if category != "未分类" {
			fileInfo.Type = "filename"
			addFileToCategory(category, fileInfo)
			results.FirstStepClassified++
		} else {
			unclassifiedFiles = append(unclassifiedFiles, fileInfo)
		}
		results.Processed++
	}

	log.Printf("第一步分类完成: %d 个文件被分类", results.FirstStepClassified)
	log.Printf("待AI分析文件: %d 个", len(unclassifiedFiles))

	// 第二步：AI分析剩余文件
	for _, file := range unclassifiedFiles {
		aiCategory := classifyByAI(file.Path, file.Name)
		file.Type = "ai"
		addFileToCategory(aiCategory, file)
		results.AIClassified++
	}

	log.Printf("AI分析完成: %d 个文件被分类", results.AIClassified)

	// 更新结果中的分类统计
	results.Classifications = classificationStats

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "文件分类完成",
		Results: results,
	})
}

// 获取分类统计
func statsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, classificationStats)
}

// 获取指定分类的文件列表
func filesHandler(c *gin.Context) {
	category := c.Param("category")

	if stats, exists := classificationStats[category]; exists {
		c.JSON(http.StatusOK, stats)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
	}
}

func main() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 确保上传目录存在
	ensureUploadDir()

	// 创建Gin引擎
	r := gin.Default()

	// 允许更大的文件上传（100MB）
	r.MaxMultipartMemory = 100 << 20

	// CORS中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// API路由 - 必须在静态文件路由之前定义
	api := r.Group("/api")
	{
		api.GET("/stats", statsHandler)
		api.GET("/files/:category", filesHandler)
	}

	// 文件上传路由
	r.POST("/upload", uploadHandler)

	// 静态文件服务 - 必须在最后定义
	r.StaticFile("/", "./public/index.html")
	r.Static("/static", "./public")

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("🚀 Go文件分类服务器启动中...\n")
	fmt.Printf("🌐 访问地址: http://localhost:%s\n", port)
	fmt.Printf("====================================\n")

	log.Fatal(r.Run(":" + port))
}
