package utils

import (
	"archive/zip"
	"bufio"
	"file-classifier/internal/extractor"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadFileContent 读取文件内容（统一调用 extractor）
func ReadFileContent(filePath string) (string, error) {
	// 直接调用统一的 extractor 处理各种文件类型
	return extractor.ExtractText(filePath)
}

// ReadUploadedFileContent 从上传的文件中读取内容
func ReadUploadedFileContent(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer file.Close()

	// 根据文件扩展名决定读取方式
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	switch ext {
	case ".txt", ".md":
		return readTextFromMultipart(file)
	case ".docx":
		// multipart 读取 docx 复杂，直接返回错误提示，后续可考虑优化
		return "", fmt.Errorf("无法直接读取 docx 内容")
	default:
		return readAsTextFromMultipart(file)
	}
}

// readTextFile 读取纯文本文件
func readTextFile(file *os.File) (string, error) {
	var content strings.Builder
	scanner := bufio.NewScanner(file)

	// 限制读取的行数，避免文件过大
	maxLines := 1000
	lineCount := 0

	for scanner.Scan() && lineCount < maxLines {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	result := content.String()

	// 限制内容长度，避免过长的内容
	if len(result) > 10000 {
		result = result[:10000] + "...[内容过长，已截断]"
	}

	return result, nil
}

// extractDocxText 提取 docx 中的纯文本
func extractDocxText(filePath string) (string, error) {
	// 打开 zip 文件
	zf, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("打开 docx 失败: %v", err)
	}
	defer zf.Close()

	var docXML *zip.File
	for _, f := range zf.File {
		if f.Name == "word/document.xml" {
			docXML = f
			break
		}
	}
	if docXML == nil {
		return "", fmt.Errorf("未找到 document.xml")
	}

	rc, err := docXML.Open()
	if err != nil {
		return "", fmt.Errorf("读取 document.xml 失败: %v", err)
	}
	defer rc.Close()

	contentBytes, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("读取内容失败: %v", err)
	}

	// 移除 XML 标签
	re := regexp.MustCompile(`<[^>]+>`)
	plain := re.ReplaceAllString(string(contentBytes), " ")
	plain = strings.TrimSpace(plain)

	// 限制长度
	if len(plain) > 10000 {
		plain = plain[:10000] + "...[内容过长，已截断]"
	}
	return plain, nil
}

// readAsText 尝试将文件读取为文本
func readAsText(file *os.File) (string, error) {
	// 读取前1024字节来判断是否为文本文件
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 检查是否包含过多的二进制字符
	if isBinaryData(buffer[:n]) {
		return "", fmt.Errorf("文件为二进制格式，无法读取文本内容")
	}

	// 重置文件指针到开头
	file.Seek(0, 0)

	// 读取完整内容
	var content strings.Builder
	scanner := bufio.NewScanner(file)

	maxLines := 1000
	lineCount := 0

	for scanner.Scan() && lineCount < maxLines {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	result := content.String()

	// 限制内容长度
	if len(result) > 10000 {
		result = result[:10000] + "...[内容过长，已截断]"
	}

	return result, nil
}

// readTextFromMultipart 从multipart文件读取文本
func readTextFromMultipart(file multipart.File) (string, error) {
	var content strings.Builder
	scanner := bufio.NewScanner(file)

	maxLines := 1000
	lineCount := 0

	for scanner.Scan() && lineCount < maxLines {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	result := content.String()

	if len(result) > 10000 {
		result = result[:10000] + "...[内容过长，已截断]"
	}

	return result, nil
}

// readAsTextFromMultipart 尝试从multipart文件读取文本
func readAsTextFromMultipart(file multipart.File) (string, error) {
	// 读取前1024字节来判断是否为文本文件
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 检查是否包含过多的二进制字符
	if isBinaryData(buffer[:n]) {
		return "", fmt.Errorf("文件为二进制格式，无法读取文本内容")
	}

	// 重置文件指针到开头
	file.Seek(0, 0)

	// 读取完整内容
	var content strings.Builder
	scanner := bufio.NewScanner(file)

	maxLines := 1000
	lineCount := 0

	for scanner.Scan() && lineCount < maxLines {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	result := content.String()

	if len(result) > 10000 {
		result = result[:10000] + "...[内容过长，已截断]"
	}

	return result, nil
}

// isBinaryData 检查数据是否为二进制数据
func isBinaryData(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// 计算非打印字符的比例
	nonPrintableCount := 0
	for _, b := range data {
		// ASCII控制字符（除了\t, \n, \r）
		if b < 32 && b != 9 && b != 10 && b != 13 {
			nonPrintableCount++
		}
		// 高位字符
		if b > 126 {
			nonPrintableCount++
		}
	}

	// 如果超过30%的字符是非打印字符，认为是二进制文件
	ratio := float64(nonPrintableCount) / float64(len(data))
	return ratio > 0.3
}

// GetSafeFileName 获取安全的文件名（用于标题）
func GetSafeFileName(filename string) string {
	// 移除文件扩展名
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 限制长度
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}
