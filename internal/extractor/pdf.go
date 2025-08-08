package extractor

import (
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type pdfExtractor struct{}

func (p *pdfExtractor) Extract(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开 PDF 失败: %v", err)
	}
	defer f.Close()

	var builder strings.Builder
	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("解析 PDF 失败: %v", err)
	}
	_, err = io.Copy(&builder, b)
	if err != nil {
		return "", fmt.Errorf("读取 PDF 文本失败: %v", err)
	}

	content := builder.String()
	if len(content) == 0 {
		return "", fmt.Errorf("PDF 无可提取文本")
	}
	if len(content) > maxContentSize {
		content = content[:maxContentSize]
	}
	return content, nil
}

func init() {
	Register(".pdf", &pdfExtractor{})
}
