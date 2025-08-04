package service

import (
	"log"
	"math/rand"
	"strings"
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
func ClassifyByAI(filePath, filename string) string {
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
