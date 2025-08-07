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

	// 验证配置值
	if len(Conf.GithubTokens) != 3 {
		t.Errorf("Expected 3 GitHub tokens, got %d", len(Conf.GithubTokens))
	}

	if Conf.GithubTokens[0] != "token1" || Conf.GithubTokens[1] != "token2" || Conf.GithubTokens[2] != "token3" {
		t.Errorf("Expected tokens to be ['token1', 'token2', 'token3'], got %v", Conf.GithubTokens)
	}

	if Conf.DataPath != "/tmp/test" {
		t.Errorf("Expected DataPath to be '/tmp/test', got '%s'", Conf.DataPath)
	}

	if Conf.DateRangeDays != 365 {
		t.Errorf("Expected DateRangeDays to be 365, got %d", Conf.DateRangeDays)
	}

	if Conf.QueriesFile != "test_queries.txt" {
		t.Errorf("Expected QueriesFile to be 'test_queries.txt', got '%s'", Conf.QueriesFile)
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

	// 验证默认值
	if Conf.DataPath != "/app/data" {
		t.Errorf("Expected default DataPath to be '/app/data', got '%s'", Conf.DataPath)
	}

	if Conf.DateRangeDays != 730 {
		t.Errorf("Expected default DateRangeDays to be 730, got %d", Conf.DateRangeDays)
	}

	if Conf.QueriesFile != "queries.txt" {
		t.Errorf("Expected default QueriesFile to be 'queries.txt', got '%s'", Conf.QueriesFile)
	}
}
