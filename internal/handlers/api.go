package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"file-classifier/internal/config"
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
