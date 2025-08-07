package extractor

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TextExtractor 定义统一的文本提取接口
// Extract 返回文件中的纯文本内容，若失败返回 error
// 实现需自行裁剪过长内容
//
// 注意：所有实现应保证线程安全
//
// 注册时请使用小写扩展名（包含点），如 ".txt"
type TextExtractor interface {
	Extract(path string) (string, error)
}

var (
	registry                     = make(map[string]TextExtractor)
	fallback       TextExtractor = &defaultExtractor{}
	maxContentSize               = 10000 // 最大截取字符数
)

// Register 在 init() 中调用，注册对应扩展名的提取器
func Register(ext string, e TextExtractor) {
	registry[strings.ToLower(ext)] = e
}

// ExtractText 根据扩展名分发到具体实现
// 未注册的类型使用 fallback
func ExtractText(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if extractor, ok := registry[ext]; ok {
		return extractor.Extract(path)
	}
	return fallback.Extract(path)
}

// ------- 默认提取器 -------

type defaultExtractor struct{}

// 读取前若干行文本，若检测到二进制/不可读则返回提示
func (d *defaultExtractor) Extract(path string) (string, error) {
	return "", fmt.Errorf("不支持的文件类型: %s", filepath.Ext(path))
}
