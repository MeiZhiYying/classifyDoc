package extractor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type plainTextExtractor struct{}

func (p *plainTextExtractor) Extract(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() && builder.Len() < maxContentSize {
		builder.WriteString(scanner.Text())
		builder.WriteString("\n")
		line++
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	content := builder.String()
	if len(content) > maxContentSize {
		content = content[:maxContentSize]
	}
	return content, nil
}

func init() {
	exts := []string{".txt", ".md", ".log", ".csv", ".json", ".yaml", ".yml"}
	p := &plainTextExtractor{}
	for _, ext := range exts {
		Register(ext, p)
	}
}
