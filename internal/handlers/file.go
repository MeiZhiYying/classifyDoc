package handlers

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"file-classifier/internal/service"
	"github.com/gin-gonic/gin"
)

// FileHandler 处理文件访问
func FileHandler(c *gin.Context) {
	filePath := c.Param("filepath")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径不能为空"})
		return
	}

	// 解码文件路径
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径格式错误"})
		return
	}

	// 构建完整的文件路径
	var fullPath string

	// 如果路径已经是绝对路径，直接使用
	if filepath.IsAbs(decodedPath) {
		fullPath = decodedPath
	} else {
		// 如果是相对路径，与uploads目录拼接
		fullPath = filepath.Join("uploads", decodedPath)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析文件路径"})
		return
	}

	// 安全检查：确保文件路径在uploads目录内
	absUploadDir, err := filepath.Abs("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析uploads目录"})
		return
	}

	if !strings.HasPrefix(absPath, absUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置适当的Content-Type
	ext := strings.ToLower(filepath.Ext(absPath))
	contentType := getContentType(ext)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}

	// 对于图片、PDF等文件，直接在浏览器中显示
	if isDisplayableFile(ext) {
		c.File(absPath)
	} else {
		// 对于其他文件，提供下载
		c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absPath))
		c.File(absPath)
	}
}

// DownloadHandler 处理文件下载
func DownloadHandler(c *gin.Context) {
	filePath := c.Param("filepath")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径不能为空"})
		return
	}

	// 解码文件路径
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件路径格式错误"})
		return
	}

	// 构建完整的文件路径
	var fullPath string

	// 如果路径已经是绝对路径，直接使用
	if filepath.IsAbs(decodedPath) {
		fullPath = decodedPath
	} else {
		// 如果是相对路径，与uploads目录拼接
		fullPath = filepath.Join("uploads", decodedPath)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析文件路径"})
		return
	}

	// 安全检查：确保文件路径在uploads目录内
	absUploadDir, err := filepath.Abs("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法解析uploads目录"})
		return
	}

	if !strings.HasPrefix(absPath, absUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 强制下载
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(absPath))
	c.File(absPath)
}

// DeleteFileHandler 删除文件并从所有分类中剔除
func DeleteFileHandler(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数错误: " + err.Error()})
		return
	}
	ok, msg := service.DeleteFileFromAllCategories(req.Path)
	if ok {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": msg})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": msg})
	}
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".log":  "text/plain",
		".csv":  "text/csv",
		".json": "application/json",
		".xml":  "application/xml",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".ogg":  "audio/ogg",
		".aac":  "audio/aac",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	return "application/octet-stream"
}

// isDisplayableFile 判断文件是否可以在浏览器中直接显示
func isDisplayableFile(ext string) bool {
	displayableExts := []string{
		".pdf", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg",
		".txt", ".html", ".css", ".js", ".json", ".xml",
		".mp4", ".avi", ".mov", ".wmv", ".flv",
		".mp3", ".wav", ".ogg", ".aac",
		".csv", ".md", ".log",
		// 移除Office文档，因为浏览器通常不支持直接显示
		// ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	}

	for _, displayableExt := range displayableExts {
		if ext == displayableExt {
			return true
		}
	}
	return false
}
