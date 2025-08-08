package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AIClassificationRequest WPS AI API请求结构
type AIClassificationRequest struct {
	UID                string         `json:"uid"`
	Stream             bool           `json:"stream"`
	FunctionCode       string         `json:"function_code"`
	FunctionParameters FunctionParams `json:"function_parameters"`
	SecText            SecurityText   `json:"sec_text"`
}

// FunctionParams 功能参数
type FunctionParams struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	CandidateTagList string `json:"candidate_tag_list"`
}

// SecurityText 安全文本
type SecurityText struct {
	From  string `json:"from"`
	Scene string `json:"scene"`
}

// AIClassificationResponse WPS AI API响应结构
type AIClassificationResponse struct {
	Event        string `json:"event"`
	FunctionCode string `json:"function_code"`
	SessionID    string `json:"session_id"`
	ReplyID      string `json:"reply_id"`
	Reply        string `json:"reply"`
	Usage        Usage  `json:"usage"`
	Platform     string `json:"platform"`
	Model        string `json:"model"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// AIClassificationResult AI分类结果
type AIClassificationResult struct {
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

// WPS AI分类API配置
const (
	WPSAIURL      = "http://kpp.wps.cn/api/v2/aigc/completions?=null"
	ClientReqID   = "123456"
	BearerToken   = ""
	IntentionCode = "kdocs_public_autolabel_new"
	ProductName   = "kdocs-public-pc"
	UID           = "282987730"
)

// ClassifyWithAI 使用WPS AI进行文件分类
func ClassifyWithAI(title, content string) (string, error) {
	// 构建请求数据
	requestData := AIClassificationRequest{
		UID:          UID,
		Stream:       false,
		FunctionCode: "doc_classify",
		FunctionParameters: FunctionParams{
			Title:            title,
			Content:          content,
			CandidateTagList: "'合同', '简历', '发票', '论文', '其它分类'",
		},
		SecText: SecurityText{
			From:  "AI_WPS_VIP",
			Scene: "ai_autolabel",
		},
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", WPSAIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ")
	req.Header.Set("Client-Request-Id", uuid.NewString())
	req.Header.Set("Ai-Gateway-Intention-Code", IntentionCode)
	req.Header.Set("Ai-Gateway-Product-Name", ProductName)
	req.Header.Set("User-Agent", "PostmanRuntime/7.32.3")

	// ======================================================
	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second, // 30秒超时
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送HTTP请求失败: %v", err)
		// Fallback to local classification
		return simpleContentClassifier(content), nil
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return simpleContentClassifier(content), nil
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
		return simpleContentClassifier(content), nil
	}

	// 解析响应
	var aiResponse AIClassificationResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		log.Printf("解析API响应失败: %v", err)
		return simpleContentClassifier(content), nil
	}

	// 解析AI分类结果
	var result AIClassificationResult
	if err := json.Unmarshal([]byte(aiResponse.Reply), &result); err != nil {
		// 如果解析失败，尝试从reply中提取分类信息
		log.Printf("解析AI分类结果失败，使用原始回复: %s", aiResponse.Reply)
		reply := strings.ToLower(aiResponse.Reply)
		if strings.Contains(reply, "合同") {
			return "合同", nil
		} else if strings.Contains(reply, "简历") {
			return "简历", nil
		} else if strings.Contains(reply, "发票") {
			return "发票", nil
		} else if strings.Contains(reply, "论文") {
			return "论文", nil
		} else {
			return "未分类", nil
		}
	}

	// 记录AI分析结果
	log.Printf("AI分析结果: %s -> %s (置信度: %.2f, 原因: %s)",
		title, result.Category, result.Confidence, result.Reason)

	// 映射分类结果到我们的分类系统
	switch result.Category {
	case "合同", "简历", "发票", "论文":
		return result.Category, nil
	case "其它分类":
		return "未分类", nil
	default:
		return simpleContentClassifier(content), nil
	}
}

// simpleContentClassifier 本地简单关键词分类器（网络失败时使用）
func simpleContentClassifier(content string) string {
	lower := strings.ToLower(content)
	if strings.Contains(lower, "合同") || strings.Contains(lower, "agreement") {
		return "合同"
	}
	if strings.Contains(lower, "简历") || strings.Contains(lower, "resume") {
		return "简历"
	}
	if strings.Contains(lower, "发票") || strings.Contains(lower, "invoice") {
		return "发票"
	}
	if strings.Contains(lower, "论文") || strings.Contains(lower, "thesis") || strings.Contains(lower, "paper") {
		return "论文"
	}
	return "未分类"
}
