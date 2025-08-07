package config

import (
	"file-classifier/internal/models"
	"sync"
)

// ClassificationStats 全局分类统计
var ClassificationStats = map[string]models.CategoryStats{
	"合同":  {Count: 0, Files: []models.FileInfo{}},
	"简历":  {Count: 0, Files: []models.FileInfo{}},
	"发票":  {Count: 0, Files: []models.FileInfo{}},
	"论文":  {Count: 0, Files: []models.FileInfo{}},
	"未分类": {Count: 0, Files: []models.FileInfo{}},
}

// ClassificationKeywords 分类关键词配置
var ClassificationKeywords = map[string][]string{
	"合同": {"合同", "协议", "契约", "contract", "agreement", "合作", "签署"},
	"简历": {"简历", "履历", "resume", "cv", "个人简历", "求职", "应聘"},
	"发票": {"发票", "票据", "invoice", "收据", "账单", "bill", "费用"},
	"论文": {"论文", "研究", "paper", "thesis", "学术", "期刊", "研究报告", "报告", "毕业论文"},
}

// 添加互斥锁来保护动态分类
var (
	StatsMutex    sync.RWMutex
	KeywordsMutex sync.RWMutex
)

// AddCategory 动态添加分类
func AddCategory(categoryName string, keywords []string) {
	KeywordsMutex.Lock()
	defer KeywordsMutex.Unlock()

	// 添加分类关键词
	ClassificationKeywords[categoryName] = keywords

	// 添加分类统计
	StatsMutex.Lock()
	defer StatsMutex.Unlock()
	ClassificationStats[categoryName] = models.CategoryStats{Count: 0, Files: []models.FileInfo{}}
}

// GetClassificationStats 线程安全地获取分类统计
func GetClassificationStats() map[string]models.CategoryStats {
	StatsMutex.RLock()
	defer StatsMutex.RUnlock()

	// 创建副本以避免并发问题
	result := make(map[string]models.CategoryStats)
	for k, v := range ClassificationStats {
		result[k] = v
	}
	return result
}

// GetClassificationKeywords 线程安全地获取分类关键词
func GetClassificationKeywords() map[string][]string {
	KeywordsMutex.RLock()
	defer KeywordsMutex.RUnlock()

	// 创建副本以避免并发问题
	result := make(map[string][]string)
	for k, v := range ClassificationKeywords {
		result[k] = v
	}
	return result
}

// IsPredefinedCategory 判断是否为预定义分类
func IsPredefinedCategory(categoryName string) bool {
	predefinedCategories := []string{"合同", "简历", "发票", "论文", "未分类"}
	for _, cat := range predefinedCategories {
		if cat == categoryName {
			return true
		}
	}
	return false
}

// Server 服务器配置
const (
	DefaultPort  = "3000"
	UploadDir    = "uploads"
	MaxFileSize  = 100 << 20 // 100MB
	MaxFileCount = 200
	StaticDir    = "./public"
	IndexFile    = "./public/index.html"
)
