package models

import "time"

// FileInfo 文件信息结构
type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	Type     string    `json:"type"`     // "filename", "ai", "failed"
	Category string    `json:"category"` // 文件分类
	ModTime  time.Time `json:"modTime"`  // 修改时间
}

// CategoryStats 分类统计结构
type CategoryStats struct {
	Count int        `json:"count"`
	Files []FileInfo `json:"files"`
}

// UploadResult 上传结果结构
type UploadResult struct {
	Total               int                      `json:"total"`
	Processed           int                      `json:"processed"`
	FirstStepClassified int                      `json:"firstStepClassified"`
	AIClassified        int                      `json:"aiClassified"`
	Classifications     map[string]CategoryStats `json:"classifications"`
}

// Response 响应结构
type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Results *UploadResult `json:"results,omitempty"`
	Error   string        `json:"error,omitempty"`
}
