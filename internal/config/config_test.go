package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 保存原始环境变量
	originalTokens := os.Getenv("GITHUB_TOKENS")
	originalDataPath := os.Getenv("DATA_PATH")
	originalDateRange := os.Getenv("DATE_RANGE_DAYS")
	originalQueriesFile := os.Getenv("QUERIES_FILE")

	// 确保我们恢复原始环境变量
	defer func() {
		os.Setenv("GITHUB_TOKENS", originalTokens)
		os.Setenv("DATA_PATH", originalDataPath)
		os.Setenv("DATE_RANGE_DAYS", originalDateRange)
		os.Setenv("QUERIES_FILE", originalQueriesFile)
	}()

	// 设置测试环境变量
	os.Setenv("GITHUB_TOKENS", "token1,token2,token3")
	os.Setenv("DATA_PATH", "/tmp/test")
	os.Setenv("DATE_RANGE_DAYS", "365")
	os.Setenv("QUERIES_FILE", "test_queries.txt")

	// 加载配置
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// 验证配置值
	if len(cfg.GithubTokens) != 3 {
		t.Errorf("Expected 3 GitHub tokens, got %d", len(cfg.GithubTokens))
	}

	if cfg.GithubTokens[0] != "token1" || cfg.GithubTokens[1] != "token2" || cfg.GithubTokens[2] != "token3" {
		t.Errorf("Expected tokens to be ['token1', 'token2', 'token3'], got %v", cfg.GithubTokens)
	}

	if cfg.DataPath != "/tmp/test" {
		t.Errorf("Expected DataPath to be '/tmp/test', got '%s'", cfg.DataPath)
	}

	if cfg.DateRangeDays != 365 {
		t.Errorf("Expected DateRangeDays to be 365, got %d", cfg.DateRangeDays)
	}

	if cfg.QueriesFile != "test_queries.txt" {
		t.Errorf("Expected QueriesFile to be 'test_queries.txt', got '%s'", cfg.QueriesFile)
	}
}

func TestLoadDefaultValues(t *testing.T) {
	// 保存原始环境变量
	originalTokens := os.Getenv("GITHUB_TOKENS")
	originalDataPath := os.Getenv("DATA_PATH")
	originalDateRange := os.Getenv("DATE_RANGE_DAYS")
	originalQueriesFile := os.Getenv("QUERIES_FILE")

	// 确保我们恢复原始环境变量
	defer func() {
		os.Setenv("GITHUB_TOKENS", originalTokens)
		os.Setenv("DATA_PATH", originalDataPath)
		os.Setenv("DATE_RANGE_DAYS", originalDateRange)
		os.Setenv("QUERIES_FILE", originalQueriesFile)
	}()

	// Set only required environment variable
	os.Setenv("GITHUB_TOKENS", "token1")
	os.Unsetenv("DATA_PATH")
	os.Unsetenv("DATE_RANGE_DAYS")
	os.Unsetenv("QUERIES_FILE")

	// 加载配置
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// 验证默认值
	if cfg.DataPath != "/app/data" {
		t.Errorf("Expected default DataPath to be '/app/data', got '%s'", cfg.DataPath)
	}

	if cfg.DateRangeDays != 730 {
		t.Errorf("Expected default DateRangeDays to be 730, got %d", cfg.DateRangeDays)
	}

	if cfg.QueriesFile != "queries.txt" {
		t.Errorf("Expected default QueriesFile to be 'queries.txt', got '%s'", cfg.QueriesFile)
	}
}
