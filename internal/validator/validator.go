package validator

import (
	"context"
	"fmt"
)

// SecretType 表示密钥的类型
type SecretType string

const (
	GoogleAPIKey      SecretType = "google_api_key"
	GoogleOAuthID     SecretType = "google_oauth_id"
	GoogleOAuthSecret SecretType = "google_oauth_secret"
	AWSAccessKey      SecretType = "aws_access_key"
	AWSSecretKey      SecretType = "aws_secret_key"
	GitHubToken       SecretType = "github_token"
	GitHubOAuthToken  SecretType = "github_oauth_token"
	GitHubPAT         SecretType = "github_pat"
	GenericAPIKey     SecretType = "generic_api_key"
)

// ValidationResult 表示验证结果
type ValidationResult struct {
	Valid     bool
	Type      SecretType
	Value     string
	Details   string
	ErrorCode string
}

// Validator 验证密钥的接口
type Validator interface {
	Validate(ctx context.Context, secretType SecretType, value string) (*ValidationResult, error)
	SupportedTypes() []SecretType
}

// NewValidationResult 创建一个新的 ValidationResult
func NewValidationResult(valid bool, secretType SecretType, value, details, errorCode string) *ValidationResult {
	return &ValidationResult{
		Valid:     valid,
		Type:      secretType,
		Value:     value,
		Details:   details,
		ErrorCode: errorCode,
	}
}

// MultiValidator 组合多个验证器
type MultiValidator struct {
	validators []Validator
}

// NewMultiValidator 创建一个新的 MultiValidator
func NewMultiValidator(validators ...Validator) *MultiValidator {
	return &MultiValidator{
		validators: validators,
	}
}

// Validate 使用所有可用的验证器验证密钥
func (mv *MultiValidator) Validate(ctx context.Context, secretType SecretType, value string) (*ValidationResult, error) {
	for _, validator := range mv.validators {
		for _, supportedType := range validator.SupportedTypes() {
			if supportedType == secretType {
				return validator.Validate(ctx, secretType, value)
			}
		}
	}

	return nil, fmt.Errorf("未找到密钥类型的验证器: %s", secretType)
}

// SupportedTypes 返回所有支持的密钥类型
func (mv *MultiValidator) SupportedTypes() []SecretType {
	var types []SecretType
	for _, validator := range mv.validators {
		types = append(types, validator.SupportedTypes()...)
	}
	return types
}
