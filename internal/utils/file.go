package utils

import (
	"os"

	"file-classifier/internal/config"
)

// EnsureUploadDir 确保上传目录存在
func EnsureUploadDir() {
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		os.Mkdir(config.UploadDir, 0755)
	}
}

// GetPort 获取服务器端口
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = config.DefaultPort
	}
	return port
}
