package validator

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/genai"
)

// GoogleValidator 验证 Google API 密钥
type GoogleValidator struct{}

// NewGoogleValidator 创建一个新的 GoogleValidator
func NewGoogleValidator() *GoogleValidator {
	return &GoogleValidator{}
}

// Validate 验证 Google API 密钥
func (gv *GoogleValidator) Validate(ctx context.Context, secretType SecretType, value string) (*ValidationResult, error) {
	switch secretType {
	case GoogleAPIKey:
		return gv.validateAPIKey(ctx, value)
	default:
		return nil, fmt.Errorf("不支持的密钥类型: %s", secretType)
	}
}

func (gv *GoogleValidator) validateAPIKey(ctx context.Context, apiKey string) (*ValidationResult, error) {
	option := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}
	client, err := genai.NewClient(ctx, option)
	parts := []*genai.Part{
		{Text: "hi"},
	}
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %w", err)
	}
	_, err = client.Models.GenerateContent(ctx, "gemini-2.0-flash", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		if googleApiError, ok := err.(genai.APIError); ok {
			code := strconv.Itoa(googleApiError.Code)
			return NewValidationResult(false, GoogleAPIKey, apiKey, "API 密钥无效", code), nil
		} else {
			return nil, fmt.Errorf("生成内容失败: %w", err)
		}

	}
	validateResult := NewValidationResult(true, GoogleAPIKey, apiKey, "API 密钥有效", "")
	// 付费key校验
	payText := strings.Repeat("You are an expert at analyzing transcripts.", 150)
	parts = []*genai.Part{
		{Text: payText},
	}
	_, err = client.Caches.Create(ctx, "gemini-2.5-flash", &genai.CreateCachedContentConfig{
		Contents: []*genai.Content{{Parts: parts, Role: genai.RoleUser}},
		TTL:      300,
	})
	if err != nil {
		if googleApierror, ok := err.(genai.APIError); ok {
			code := strconv.Itoa(googleApierror.Code)
			validateResult.Details = "检测到免费密钥，付费密钥失败code：" + code
			return validateResult, nil
		}
	}
	validateResult.Pay = true
	return validateResult, nil
}

// SupportedTypes 返回支持的密钥类型
func (gv *GoogleValidator) SupportedTypes() []SecretType {
	return []SecretType{GoogleAPIKey}
}
