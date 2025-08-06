package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"file-classifier/internal/config"
	"file-classifier/internal/models"
	"file-classifier/internal/service"
)

// AddCategoryHandler 添加新分类
func AddCategoryHandler(c *gin.Context) {
	var request struct {
		CategoryName string `json:"categoryName" binding:"required"`
		Username     string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查新增分类数量限制（最多3个）
	keywords := config.GetClassificationKeywords()
	newCategoryCount := 0
	for categoryName := range keywords {
		if !config.IsPredefinedCategory(categoryName) {
			newCategoryCount++
		}
	}

	if newCategoryCount >= 3 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "新增分类数量已达上限（最多3个），请删除一些分类后再添加",
		})
		return
	}

	// 使用用户名作为关键词
	keywordList := []string{request.Username}

	// 添加新分类
	config.AddCategory(request.CategoryName, keywordList)

	// 重新扫描文件以应用新分类
	go func() {
		// 在后台重新扫描文件
		service.ResetClassificationStats()

		// 检查uploads目录是否存在
		if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
			return
		}

		// 遍历uploads目录并重新分类
		filepath.Walk(config.UploadDir, func(path string, info os.FileInfo, err error) error {
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
				relPath, err := filepath.Rel(config.UploadDir, path)
				if err != nil {
					return err
				}

				// 分类文件
				category := service.ClassifyByFilename(info.Name())

				// 创建文件信息
				fileInfo := models.FileInfo{
					Name:     info.Name(),
					Path:     relPath,
					Size:     info.Size(),
					Type:     "filename",
					Category: category,
					ModTime:  info.ModTime(),
				}

				// 添加到分类
				service.AddFileToCategory(category, fileInfo)
			}
			return nil
		})
	}()

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: fmt.Sprintf("分类 '%s' 添加成功，关键词: %s", request.CategoryName, request.Username),
	})
}

// GetCategoriesHandler 获取所有分类
func GetCategoriesHandler(c *gin.Context) {
	categories := config.GetClassificationKeywords()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"categories": categories,
	})
}

// StatsHandler 获取分类统计
func StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetClassificationStats())
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

// AllFilesHandler 获取所有文件列表
func AllFilesHandler(c *gin.Context) {
	// 获取查询参数
	sortBy := c.DefaultQuery("sort", "time")         // 排序方式：time(默认), size
	order := c.DefaultQuery("order", "desc")         // 排序顺序：desc(默认), asc
	filterCategory := c.DefaultQuery("category", "") // 分类筛选：空表示全部

	var allFiles []models.FileInfo

	// 遍历所有分类，收集文件
	for categoryName, categoryStats := range config.ClassificationStats {
		// 如果指定了分类筛选，只返回对应分类的文件
		if filterCategory != "" && categoryName != filterCategory {
			continue
		}

		for _, file := range categoryStats.Files {
			// 为每个文件添加分类信息
			fileWithCategory := file
			fileWithCategory.Category = categoryName

			// 获取文件的修改时间
			fullPath := filepath.Join(config.UploadDir, file.Path)
			if fileInfo, err := os.Stat(fullPath); err == nil {
				fileWithCategory.ModTime = fileInfo.ModTime()
			} else {
				// 如果无法获取文件信息，使用当前时间
				fileWithCategory.ModTime = time.Now()
			}

			allFiles = append(allFiles, fileWithCategory)
		}
	}

	// 根据指定方式排序
	switch sortBy {
	case "size":
		if order == "asc" {
			sort.Slice(allFiles, func(i, j int) bool {
				return allFiles[i].Size < allFiles[j].Size
			})
		} else {
			sort.Slice(allFiles, func(i, j int) bool {
				return allFiles[i].Size > allFiles[j].Size
			})
		}
	case "time":
		fallthrough
	default:
		if order == "asc" {
			sort.Slice(allFiles, func(i, j int) bool {
				return allFiles[i].ModTime.Before(allFiles[j].ModTime)
			})
		} else {
			sort.Slice(allFiles, func(i, j int) bool {
				return allFiles[i].ModTime.After(allFiles[j].ModTime)
			})
		}
	}

	// 构建响应
	response := gin.H{
		"success": true,
		"total":   len(allFiles),
		"files":   allFiles,
		"sort":    sortBy,
		"order":   order,
		"filter":  filterCategory,
	}

	c.JSON(http.StatusOK, response)
}
