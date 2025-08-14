package github

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// 测试创建一个没有令牌的客户端
	client := NewClient([]string{})
	if client == nil {
		t.Error("Expected a client, got nil")
	}

	// 测试创建一个带有令牌的客户端
	tokens := []string{"token1", "token2", "token3"}
	client = NewClient(tokens)
	if client == nil {
		t.Error("Expected a client, got nil")
	}

	if len(client.tokens) != len(tokens) {
		t.Errorf("Expected %d tokens, got %d", len(tokens), len(client.tokens))
	}
}
