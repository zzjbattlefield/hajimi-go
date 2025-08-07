package extractor

import (
	"regexp"
)

var SecretExtractor *Extractor

// Secret 表示提取出的密钥
type Secret struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	File   string `json:"file"`
	Repo   string `json:"repo"`
	Line   int    `json:"line"`
	Commit string `json:"commit"`
}

// Extractor 处理密钥提取
type Extractor struct {
	patterns map[string]*regexp.Regexp
}

// NewExtractor 创建一个新的提取器，包含预定义的模式
func NewExtractor() {
	patterns := map[string]*regexp.Regexp{
		// Google API Key
		"google_api_key": regexp.MustCompile(`AIzaSy[A-Za-z0-9\-_]{33}`),
	}

	SecretExtractor = &Extractor{
		patterns: patterns,
	}
}

// Extract 从内容中提取密钥
func (e *Extractor) Extract(content, filename, repo, commit string) []Secret {
	var secrets []Secret

	for secretType, pattern := range e.patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 0 {
				// 查找行号
				line := 1
				value := match[0]

				// 简单的行计数
				for i := 0; i < len(content); i++ {
					if content[i] == '\n' {
						line++
					}
					if i+len(value) <= len(content) && content[i:i+len(value)] == value {
						break
					}
				}

				secrets = append(secrets, Secret{
					Type:   secretType,
					Value:  value,
					File:   filename,
					Repo:   repo,
					Line:   line,
					Commit: commit,
				})
			}
		}
	}

	return secrets
}

// AddPattern 向提取器添加自定义模式
func (e *Extractor) AddPattern(name string, pattern *regexp.Regexp) {
	e.patterns[name] = pattern
}

// GetPatterns 返回所有模式
func (e *Extractor) GetPatterns() map[string]*regexp.Regexp {
	return e.patterns
}
