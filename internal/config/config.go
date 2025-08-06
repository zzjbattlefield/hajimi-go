package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 存储应用程序的配置信息
type Config struct {
	GithubTokens  []string
	Proxy         string
	DataPath      string
	DateRangeDays int
	QueriesFile   string

	// 同步配置
	SyncEnabled       bool
	SyncBatchSize     int
	SyncFlushInterval int // 以秒为单位
	SyncMaxRetries    int
}

// Load 从环境变量中加载配置
func Load() (*Config, error) {
	// 如果存在 .env 文件则加载它
	_ = godotenv.Load()

	cfg := &Config{
		GithubTokens:      getEnvAsSlice("GITHUB_TOKENS", []string{}),
		Proxy:             getEnv("PROXY", ""),
		DataPath:          getEnv("DATA_PATH", "/app/data"),
		DateRangeDays:     getEnvAsInt("DATE_RANGE_DAYS", 730),
		QueriesFile:       getEnv("QUERIES_FILE", "queries.txt"),
		SyncEnabled:       getEnvAsBool("SYNC_ENABLED", false),
		SyncBatchSize:     getEnvAsInt("SYNC_BATCH_SIZE", 10),
		SyncFlushInterval: getEnvAsInt("SYNC_FLUSH_INTERVAL", 30), // 30秒
		SyncMaxRetries:    getEnvAsInt("SYNC_MAX_RETRIES", 3),
	}

	return cfg, nil
}

// 辅助函数
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string) []string {
	if value, exists := os.LookupEnv(name); exists {
		return strings.Split(value, ",")
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
