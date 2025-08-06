package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Test creating a new logger
	logger := New()
	if logger == nil {
		t.Error("Expected a logger, got nil")
	}

	if logger.Logger == nil {
		t.Error("Expected logger.Logger to be initialized")
	}
}

func TestLoggerInfo(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Info method
	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "[INFO] test message") {
		t.Errorf("Expected log output to contain '[INFO] test message', got '%s'", output)
	}
}

func TestLoggerError(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Error method
	logger.Error("test error")
	output := buf.String()

	if !strings.Contains(output, "[ERROR] test error") {
		t.Errorf("Expected log output to contain '[ERROR] test error', got '%s'", output)
	}
}

func TestLoggerWarn(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Warn method
	logger.Warn("test warning")
	output := buf.String()

	if !strings.Contains(output, "[WARN] test warning") {
		t.Errorf("Expected log output to contain '[WARN] test warning', got '%s'", output)
	}
}

func TestLoggerDebug(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Debug method
	logger.Debug("test debug")
	output := buf.String()

	if !strings.Contains(output, "[DEBUG] test debug") {
		t.Errorf("Expected log output to contain '[DEBUG] test debug', got '%s'", output)
	}
}

func TestLoggerInfof(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Infof method
	logger.Infof("test %s", "message")
	output := buf.String()

	if !strings.Contains(output, "[INFO] test message") {
		t.Errorf("Expected log output to contain '[INFO] test message', got '%s'", output)
	}
}

func TestLoggerErrorf(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Errorf method
	logger.Errorf("test %s", "error")
	output := buf.String()

	if !strings.Contains(output, "[ERROR] test error") {
		t.Errorf("Expected log output to contain '[ERROR] test error', got '%s'", output)
	}
}

func TestLoggerWarnf(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Warnf method
	logger.Warnf("test %s", "warning")
	output := buf.String()

	if !strings.Contains(output, "[WARN] test warning") {
		t.Errorf("Expected log output to contain '[WARN] test warning', got '%s'", output)
	}
}

func TestLoggerDebugf(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", 0),
	}

	// Test Debugf method
	logger.Debugf("test %s", "debug")
	output := buf.String()

	if !strings.Contains(output, "[DEBUG] test debug") {
		t.Errorf("Expected log output to contain '[DEBUG] test debug', got '%s'", output)
	}
}
