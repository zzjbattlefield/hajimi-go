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

func TestWithToken(t *testing.T) {
	// Test setting a specific token
	tokens := []string{"token1", "token2"}
	client := NewClient(tokens)

	newClient := client.WithToken("new_token")
	if newClient == nil {
		t.Error("Expected a new client, got nil")
	}

	// The new client should have the same tokens but use the new token
	if len(newClient.tokens) != len(client.tokens) {
		t.Errorf("Expected %d tokens, got %d", len(client.tokens), len(newClient.tokens))
	}
}
