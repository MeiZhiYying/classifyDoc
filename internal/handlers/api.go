package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"file-classifier/internal/config"
	"file-classifier/internal/models"
	"file-classifier/internal/service"
)

// StatsHandler 获取分类统计
func StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.ClassificationStats)
}

// FilesHandler 获取指定分类的文件列表
func FilesHandler(c *gin.Context) {
	category := c.Param("category")

	if stats, exists := config.ClassificationStats[category]; exists {
		c.JSON(http.StatusOK, stats)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
	}
}

// ScanUploadsHandler 扫描uploads文件夹并分类
func ScanUploadsHandler(c *gin.Context) {
	// 重置统计
	service.ResetClassificationStats()

	// 检查uploads目录是否存在
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "uploads目录不存在，无需扫描",
			Results: &models.UploadResult{
				Total:               0,
				Processed:           0,
				FirstStepClassified: 0,
				AIClassified:        0,
				Classifications:     config.ClassificationStats,
			},
		})
		return
	}

	// 遍历uploads目录
	var files []string
	var fileSizes map[string]int64 = make(map[string]int64)

	err := filepath.Walk(config.UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if path == config.UploadDir {
			return nil
		}

		// 只处理文件，跳过目录
		if !info.IsDir() {
			// 获取相对路径
			relativePath, err := filepath.Rel(config.UploadDir, path)
			if err != nil {
				return err
			}
			files = append(files, relativePath)
			fileSizes[relativePath] = info.Size()
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "扫描uploads目录失败: " + err.Error(),
		})
		return
	}

	// 并发处理找到的文件
	results := &models.UploadResult{
		Total:               len(files),
		Processed:           0,
		FirstStepClassified: 0,
		AIClassified:        0,
		Classifications:     config.ClassificationStats,
	}

	// 创建并发控制
	maxConcurrent := 30 // 最大并发数
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex // 用于保护共享数据的互斥锁

	for _, filePath := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			filename := filepath.Base(filePath)
			category := service.ClassifyByFilename(filename)

			fileInfo := models.FileInfo{
				Name: filename,
				Path: filePath,            // 使用相对路径
				Size: fileSizes[filePath], // 使用获取到的文件大小
			}

			// 使用互斥锁保护共享数据
			mu.Lock()
			if category != "未分类" {
				fileInfo.Type = "filename"
				service.AddFileToCategory(category, fileInfo)
				results.FirstStepClassified++
			} else {
				// Second step: AI classification
				aiCategory := service.ClassifyByAI(filename)
				fileInfo.Type = "AI"
				service.AddFileToCategory(aiCategory, fileInfo)
				results.AIClassified++
			}
			results.Processed++
			mu.Unlock()

		}(filePath)
	}

	// 等待所有goroutine完成
	wg.Wait()

	results.Classifications = config.ClassificationStats

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "uploads目录扫描完成",
		Results: results,
	})
}
