package extractor

import (
	"regexp"
	"testing"
)

func TestExtract(t *testing.T) {
	// 创建一个新的提取器
	NewExtractor()

	// 测试包含各种密钥的内容
	content := `
		# This is a test file with various secrets
		
		# Google API Key
		api_key = "AIzaSyC123456789012345678901234567890123"
		some_other_string = "text..."
	`

	// 提取密钥
	secrets := SecretExtractor.Extract(content, "test.txt", "test/repo", "abc123")

	// 打印找到的密钥用于调试
	for _, secret := range secrets {
		t.Logf("Found secret: Type=%s, Value=%s", secret.Type, secret.Value)
	}

	// 验证我们找到了预期数量的密钥
	// 注意：我们期望找到超过5个密钥，因为某些模式可能会匹配其他密钥的部分内容
	// 例如，google_oauth_secret模式可能会匹配google_oauth_id的部分内容
	if len(secrets) < 5 {
		t.Errorf("Expected at least 5 secrets, got %d", len(secrets))
	}

	// 验证每种密钥类型是否被正确识别
	secretTypes := make(map[string]int)
	for _, secret := range secrets {
		secretTypes[secret.Type]++
	}

	// 检查我们是否至少找到了每种预期类型的密钥
	if secretTypes["google_api_key"] < 1 {
		t.Errorf("Expected to find at least 1 'google_api_key', but found %d", secretTypes["google_api_key"])
	}

	if secretTypes["google_oauth_id"] < 1 {
		t.Errorf("Expected to find at least 1 'google_oauth_id', but found %d", secretTypes["google_oauth_id"])
	}

	if secretTypes["aws_access_key"] < 1 {
		t.Errorf("Expected to find at least 1 'aws_access_key', but found %d", secretTypes["aws_access_key"])
	}

	if secretTypes["github_token"] < 1 {
		t.Errorf("Expected to find at least 1 'github_token', but found %d", secretTypes["github_token"])
	}

	if secretTypes["generic_api_key"] < 1 {
		t.Errorf("Expected to find at least 1 'generic_api_key', but found %d", secretTypes["generic_api_key"])
	}
}

func TestExtractNoSecrets(t *testing.T) {
	// 创建一个新的提取器
	NewExtractor()

	// 测试不包含密钥的内容
	content := `
		# This is a test file with no secrets
		
		regular_text = "this is just regular text"
		more_text = "another line of regular text"
	`

	// 提取密钥
	secrets := SecretExtractor.Extract(content, "test.txt", "test/repo", "abc123")

	// 验证我们没有找到任何密钥
	if len(secrets) != 0 {
		t.Errorf("Expected 0 secrets, got %d", len(secrets))
	}
}

func TestAddPattern(t *testing.T) {
	// 创建一个新的提取器
	NewExtractor()

	// 计算初始模式数量
	initialCount := len(SecretExtractor.GetPatterns())

	// 添加一个新模式
	pattern := regexp.MustCompile(`test_pattern_[a-z]+`)
	SecretExtractor.AddPattern("test_pattern", pattern)

	// 验证模式已添加
	patterns := SecretExtractor.GetPatterns()
	if len(patterns) != initialCount+1 {
		t.Errorf("Expected %d patterns, got %d", initialCount+1, len(patterns))
	}

	// 验证模式存在
	if patterns["test_pattern"] != pattern {
		t.Error("Expected to find the added pattern")
	}
}
