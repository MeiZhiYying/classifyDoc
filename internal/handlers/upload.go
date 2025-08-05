package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

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
	service.CheckFiles(c, files)
	// 重置统计
	service.ResetClassificationStats()
	log.Printf("开始处理 %d 个文件", len(files))
	service.ClassificDOC(c, files)
	// 更新结果中的分类统计

}
