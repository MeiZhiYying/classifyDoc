package service

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"file-classifier/internal/config"
	"file-classifier/internal/models"
	"file-classifier/internal/utils"
)

// ClassifyByFilename 根据文件名进行分类
func ClassifyByFilename(filename string) string {
	lowerName := strings.ToLower(filename)

	for category, keywords := range config.ClassificationKeywords {
		for _, keyword := range keywords {
			if strings.Contains(lowerName, strings.ToLower(keyword)) {
				return category
			}
		}
	}
	return "未分类"
}

// ClassifyByAI AI分析功能 - 集成WPS AI API
func ClassifyByAI(filename, content string) string {
	// 获取安全的文件标题
	title := utils.GetSafeFileName(filename)

	// 如果内容为空，使用文件名作为内容
	if strings.TrimSpace(content) == "" {
		content = filename
	}

	// 调用WPS AI分类API
	category, err := ClassifyWithAI(title, content)
	if err != nil {
		log.Printf("AI分析失败: %s, 错误: %v", filename, err)
		// AI分析失败时返回未分类
		return "未分类"
	}

	log.Printf("AI分析成功: %s -> %s", filename, category)
	return category
}

// ResetClassificationStats 重置分类统计
func ResetClassificationStats() {
	for key := range config.ClassificationStats {
		config.ClassificationStats[key] = models.CategoryStats{Count: 0, Files: []models.FileInfo{}}
	}
}

// AddFileToCategory 添加文件到分类
func AddFileToCategory(category string, fileInfo models.FileInfo) {
	// 使用全局互斥锁保护共享数据
	config.StatsMutex.Lock()
	defer config.StatsMutex.Unlock()

	if stats, exists := config.ClassificationStats[category]; exists {
		stats.Files = append(stats.Files, fileInfo)
		stats.Count = len(stats.Files)
		config.ClassificationStats[category] = stats
	}
}

func CheckFiles(c *gin.Context, files []*multipart.FileHeader) {
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "没有上传文件",
		})
		return
	}
	if len(files) > config.MaxFileCount {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "最多支持上传200个文件",
		})
		return
	}

}
func ClassificDOC(c *gin.Context, files []*multipart.FileHeader) {
	results := &models.UploadResult{
		Total:               len(files),
		Processed:           0,
		FirstStepClassified: 0,
		AIClassified:        0,
		Classifications:     config.ClassificationStats,
	}

	// 使用原子操作保证线程安全的计数
	var processed, firstStepClassified, aiClassified int64

	// 为每个文件创建一个goroutine，实现最大并发
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			filename := file.Filename

			// 第一步：关键词检查
			category := ClassifyByFilename(filename)

			// 保存文件
			savePath := filepath.Join(config.UploadDir, filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				log.Printf("保存文件失败: %s, %v", filename, err)
				atomic.AddInt64(&processed, 1)
				return
			}

			fileInfo := models.FileInfo{
				Name: filename,
				Path: savePath,
				Size: file.Size,
			}

			if category != "未分类" {
				// 关键词分类成功
				fileInfo.Type = "filename"
				fileInfo.Category = category
				AddFileToCategory(category, fileInfo)
				atomic.AddInt64(&firstStepClassified, 1)
				log.Printf("关键词分类成功: %s -> %s", filename, category)
			} else {
				// 第二步：AI分析（只有关键词分类失败的文件才进入此步骤）
				log.Printf("开始AI分析: %s", filename)

				// 读取文件内容
				content, err := utils.ReadFileContent(savePath)
				if err != nil {
					log.Printf("读取文件内容失败: %s, %v", filename, err)
					// 即使读取内容失败，仍尝试用文件名进行AI分析
					content = ""
				}

				// 调用AI分析
				aiCategory := ClassifyByAI(filename, content)
				fileInfo.Type = "AI"
				fileInfo.Category = aiCategory
				AddFileToCategory(aiCategory, fileInfo)
				atomic.AddInt64(&aiClassified, 1)
				log.Printf("AI分析完成: %s -> %s", filename, aiCategory)
			}

			atomic.AddInt64(&processed, 1)
		}(file)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 更新结果统计
	results.Processed = int(processed)
	results.FirstStepClassified = int(firstStepClassified)
	results.AIClassified = int(aiClassified)
	results.Classifications = config.GetClassificationStats()

	log.Printf("================= 分类完成 =================")
	log.Printf("总文件数: %d", results.Total)
	log.Printf("已处理: %d", results.Processed)
	log.Printf("关键词分类成功: %d", results.FirstStepClassified)
	log.Printf("AI分析成功: %d", results.AIClassified)
	log.Printf("==========================================")

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "文件分类完成",
		Results: results,
	})
}
