package validator

import (
	"context"
	"testing"
)

var testKey = "your_test_key_here"

func TestGoogleValidator_SupportedTypes(t *testing.T) {
	validator := NewGoogleValidator()
	types := validator.SupportedTypes()

	if len(types) != 1 {
		t.Errorf("Expected 1 supported type, got %d", len(types))
	}

	if types[0] != GoogleAPIKey {
		t.Errorf("Expected GoogleAPIKey, got %s", types[0])
	}
}

func TestGoogleValidator_ValidateUnsupportedType(t *testing.T) {
	validator := NewGoogleValidator()
	ctx := context.Background()

	_, err := validator.Validate(ctx, GitHubToken, "some_token")
	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
}

func TestGoogleValidaotr(t *testing.T) {
	validator := NewGoogleValidator()
	ctx := context.Background()

	_, err := validator.Validate(ctx, GoogleAPIKey, testKey)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
