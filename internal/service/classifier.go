package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"file-classifier/internal/config"
	"file-classifier/internal/models"
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

// ClassifyByAI AI分析占位符函数
func ClassifyByAI(filename string) string {
	// 模拟AI分析延迟
	time.Sleep(1 * time.Second)

	// 模拟AI分析结果（随机返回一个分类）
	categories := []string{"合同", "简历", "发票", "论文"}
	randomCategory := categories[rand.Intn(len(categories))]

	log.Printf("AI分析结果: %s -> %s", filename, randomCategory)
	return randomCategory
}

// ResetClassificationStats 重置分类统计
func ResetClassificationStats() {
	for key := range config.ClassificationStats {
		config.ClassificationStats[key] = models.CategoryStats{Count: 0, Files: []models.FileInfo{}}
	}
}

// AddFileToCategory 添加文件到分类
func AddFileToCategory(category string, fileInfo models.FileInfo) {
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

	// Create a channel to limit concurrent goroutines
	maxConcurrent := 30
	semaphore := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire a slot
			defer func() { <-semaphore }() // Release the slot

			filename := file.Filename
			category := ClassifyByFilename(filename)

			// Save file
			savePath := filepath.Join(config.UploadDir, filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				log.Printf("保存文件失败: %s, %v", filename, err)
				return
			}

			fileInfo := models.FileInfo{
				Name: filename,
				Path: savePath,
				Size: file.Size,
			}

			if category != "未分类" {
				fileInfo.Type = "filename"
				AddFileToCategory(category, fileInfo)
				results.FirstStepClassified++
			} else {
				// Second step: AI classification
				aiCategory := ClassifyByAI(filename)
				fileInfo.Type = "AI"
				AddFileToCategory(aiCategory, fileInfo)
				results.AIClassified++
			}
			results.Processed++

		}(file)
	}

	wg.Wait()

	results.Classifications = config.ClassificationStats
	log.Printf("关键词分类完成: %d 个文件被分类", results.FirstStepClassified)
	log.Printf("AI分析完成: %d 个文件被分类", results.AIClassified)
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "文件分类完成",
		Results: results,
	})
}
