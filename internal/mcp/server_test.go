package mcp

import (
	"testing"
)

func TestTruncateContent_Short(t *testing.T) {
	result := truncateContent("hello", 10)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTruncateContent_ExactLength(t *testing.T) {
	result := truncateContent("hello", 5)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTruncateContent_Truncated(t *testing.T) {
	result := truncateContent("hello world", 5)
	if result != "hello..." {
		t.Errorf("expected 'hello...', got %q", result)
	}
}

func TestTruncateContent_UTF8(t *testing.T) {
	content := "café résumé naïve"
	result := truncateContent(content, 5)
	if result != "café ..." {
		t.Errorf("expected 'café ...', got %q", result)
	}
}

func TestTruncateContent_Emojis(t *testing.T) {
	content := "🎉🎊🎈🎁🎂"
	result := truncateContent(content, 3)
	if result != "🎉🎊🎈..." {
		t.Errorf("expected '🎉🎊🎈...', got %q", result)
	}
}

func TestTruncateContent_Empty(t *testing.T) {
	result := truncateContent("", 10)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
