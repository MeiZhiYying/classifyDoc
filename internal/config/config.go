package config

import "file-classifier/internal/models"

// ClassificationStats 全局分类统计
var ClassificationStats = map[string]models.CategoryStats{
	"合同":   {Count: 0, Files: []models.FileInfo{}},
	"简历":   {Count: 0, Files: []models.FileInfo{}},
	"发票":   {Count: 0, Files: []models.FileInfo{}},
	"论文":   {Count: 0, Files: []models.FileInfo{}},
	"未分类":  {Count: 0, Files: []models.FileInfo{}},
	"新增分类": {Count: 0, Files: []models.FileInfo{}},
}

// ClassificationKeywords 分类关键词配置
var ClassificationKeywords = map[string][]string{
	"合同": {"合同", "协议", "契约", "contract", "agreement", "合作", "签署"},
	"简历": {"简历", "履历", "resume", "cv", "个人简历", "求职", "应聘"},
	"发票": {"发票", "票据", "invoice", "收据", "账单", "bill", "费用"},
	"论文": {"论文", "研究", "paper", "thesis", "学术", "期刊", "研究报告", "报告", "毕业论文"},
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
