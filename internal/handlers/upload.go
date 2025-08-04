package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"file-classifier/internal/config"
	"file-classifier/internal/models"
	"file-classifier/internal/service"
)

// UploadHandler 文件上传处理
func UploadHandler(c *gin.Context) {
	// 解析多文件上传
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "文件上传解析失败: " + err.Error(),
		})
		return
	}

	files := form.File["files"]
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

	// 重置统计
	service.ResetClassificationStats()

	results := &models.UploadResult{
		Total:               len(files),
		Processed:           0,
		FirstStepClassified: 0,
		AIClassified:        0,
		Classifications:     config.ClassificationStats,
	}

	log.Printf("开始处理 %d 个文件", len(files))

	// 第一步：根据文件名分类
	var unclassifiedFiles []models.FileInfo

	for _, file := range files {
		filename := file.Filename
		category := service.ClassifyByFilename(filename)

		// 保存文件
		savePath := filepath.Join(config.UploadDir, filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("保存文件失败: %s, %v", filename, err)
			continue
		}

		fileInfo := models.FileInfo{
			Name: filename,
			Path: savePath,
			Size: file.Size,
		}

		if category != "未分类" {
			fileInfo.Type = "filename"
			service.AddFileToCategory(category, fileInfo)
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
		aiCategory := service.ClassifyByAI(file.Path, file.Name)
		file.Type = "ai"
		service.AddFileToCategory(aiCategory, file)
		results.AIClassified++
	}

	log.Printf("AI分析完成: %d 个文件被分类", results.AIClassified)

	// 更新结果中的分类统计
	results.Classifications = config.ClassificationStats

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "文件分类完成",
		Results: results,
	})
}
