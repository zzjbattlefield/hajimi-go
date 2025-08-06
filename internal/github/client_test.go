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

func TestRotateToken(t *testing.T) {
	// 测试使用单个令牌轮换令牌
	tokens := []string{"token1"}
	client := NewClient(tokens)
	rotated := client.RotateToken()

	// With only one token, rotating should return the same client
	if rotated.tokenIdx != client.tokenIdx {
		t.Error("Expected token index to remain the same with only one token")
	}

	// Test rotating tokens with multiple tokens
	tokens = []string{"token1", "token2", "token3"}
	client = NewClient(tokens)

	// Initially should be at index 0
	if client.tokenIdx != 0 {
		t.Errorf("Expected initial token index to be 0, got %d", client.tokenIdx)
	}

	// After first rotation should be at index 1
	client = client.RotateToken()
	if client.tokenIdx != 1 {
		t.Errorf("Expected token index to be 1 after first rotation, got %d", client.tokenIdx)
	}

	// After second rotation should be at index 2
	client = client.RotateToken()
	if client.tokenIdx != 2 {
		t.Errorf("Expected token index to be 2 after second rotation, got %d", client.tokenIdx)
	}

	// After third rotation should be back at index 0
	client = client.RotateToken()
	if client.tokenIdx != 0 {
		t.Errorf("Expected token index to be 0 after third rotation, got %d", client.tokenIdx)
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
