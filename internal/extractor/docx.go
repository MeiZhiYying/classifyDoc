package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type docxExtractor struct{}

func (d *docxExtractor) Extract(path string) (string, error) {
	zf, err := zip.OpenReader(path)
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

	bytes, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("读取内容失败: %v", err)
	}

	re := regexp.MustCompile(`<[^>]+>`)
	plain := re.ReplaceAllString(string(bytes), " ")
	plain = strings.TrimSpace(plain)
	if len(plain) > maxContentSize {
		plain = plain[:maxContentSize]
	}
	return plain, nil
}

func init() {
	Register(".docx", &docxExtractor{})
}
