package models

// FileInfo 文件信息结构
type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Type string `json:"type"` // "filename", "ai", "failed"
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
