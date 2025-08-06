package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"file-classifier/internal/router"
	"file-classifier/internal/utils"
)

func main() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 确保上传目录存在
	utils.EnsureUploadDir()

	// 设置路由
	r := router.SetupRouter()

	// 获取端口
	port := utils.GetPort()

	fmt.Printf("🚀 Go文件分类服务器启动中...\n")
	fmt.Printf("🌐 访问地址: http://localhost:%s\n", port)
	fmt.Printf("====================================\n")

	// 启动服务器
	log.Fatal(r.Run(":" + port))
}
