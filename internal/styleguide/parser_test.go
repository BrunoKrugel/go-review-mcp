package styleguide

import (
	"testing"
)

func TestParseContent_Headings(t *testing.T) {
	content := `# Main Title

Some intro text.

## Section One

Content for section one.

### Subsection

More details here.

## Section Two

Content for section two.`

	idx := ParseContent(content)

	if len(idx.Sections) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(idx.Sections))
	}

	tests := []struct {
		title   string
		content string
		index   int
		level   int
	}{
		{"Main Title", "Some intro text.", 0, 1},
		{"Section One", "Content for section one.", 1, 2},
		{"Subsection", "More details here.", 2, 3},
		{"Section Two", "Content for section two.", 3, 2},
	}

	for _, tt := range tests {
		sec := idx.Sections[tt.index]
		if sec.Title != tt.title {
			t.Errorf("section %d: expected title %q, got %q", tt.index, tt.title, sec.Title)
		}
		if sec.Level != tt.level {
			t.Errorf("section %d: expected level %d, got %d", tt.index, tt.level, sec.Level)
		}
	}
}

func TestParseContent_EmptyInput(t *testing.T) {
	idx := ParseContent("")

	if len(idx.Sections) != 0 {
		t.Errorf("expected 0 sections for empty input, got %d", len(idx.Sections))
	}

	if len(idx.Topics) != 0 {
		t.Errorf("expected 0 topics for empty input, got %d", len(idx.Topics))
	}
}

func TestParseContent_HTMLStripped(t *testing.T) {
	content := `# Title

<p>Paragraph with <strong>bold</strong> text.</p>
<div class="foo">Div content</div>
<a href="https://example.com">Link text</a>`

	idx := ParseContent(content)

	if len(idx.Sections) == 0 {
		t.Fatal("expected at least 1 section")
	}

	for _, forbidden := range []string{"<p>", "</p>", "<strong>", "</strong>", "<div", "</div>", "<a ", "</a>"} {
		if contains(idx.Sections[0].Content, forbidden) {
			t.Errorf("content should not contain %q, but does: %s", forbidden, idx.Sections[0].Content)
		}
	}
}

func TestSearchContent(t *testing.T) {
	content := `# Error Handling

Always check errors explicitly.

# Naming Conventions

Use MixedCaps for names.`

	idx := ParseContent(content)

	results := idx.SearchContent("error")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'error', got %d", len(results))
	}
	if results[0].Title != "Error Handling" {
		t.Errorf("expected 'Error Handling', got %q", results[0].Title)
	}
}

func TestSearchContent_NoResults(t *testing.T) {
	content := `# Title

Some content.`

	idx := ParseContent(content)

	results := idx.SearchContent("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestGetTopic(t *testing.T) {
	content := `# Error Handling

Always check errors and use error wrapping.

# Naming Conventions

Use MixedCaps for names and identifiers.`

	idx := ParseContent(content)

	errors := idx.GetTopic("errors")
	if len(errors) == 0 {
		t.Error("expected to find sections for 'errors' topic")
	}
}

func TestGetTopic_UnknownTopic(t *testing.T) {
	content := `# Title

Some content.`

	idx := ParseContent(content)

	results := idx.GetTopic("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for unknown topic, got %d", len(results))
	}
}

func TestGetSection(t *testing.T) {
	content := `# My Section

Section content here.

# Another Section

Other content.`

	idx := ParseContent(content)

	sec, found := idx.GetSection("my section")
	if !found {
		t.Fatal("expected to find 'my section'")
	}
	if sec.Title != "My Section" {
		t.Errorf("expected 'My Section', got %q", sec.Title)
	}

	_, found = idx.GetSection("nonexistent")
	if found {
		t.Error("expected not to find 'nonexistent'")
	}
}

func TestGenerateAnchor(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"Hello World", "hello-world"},
		{"Test, With: Punctuation.", "test-with-punctuation"},
		{"It's \"quoted\"", "its-quoted"},
		{"Func (Args)", "func-args"},
	}

	for _, tt := range tests {
		result := generateAnchor(tt.input)
		if result != tt.expect {
			t.Errorf("generateAnchor(%q) = %q, want %q", tt.input, result, tt.expect)
		}
	}
}

func TestRemoveRemainingTags(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"<p>Hello</p>", "Hello"},
		{"<div class=\"foo\">bar</div>", "bar"},
		{"<a href=\"url\">link</a>", "link"},
		{"no tags here", "no tags here"},
		{"<br>", ""},
		{"<input type=\"text\" />", ""}, // self-closing tag fully stripped
		{"text <b>bold</b> more", "text bold more"},
	}

	for _, tt := range tests {
		result := removeRemainingTags(tt.input)
		if result != tt.expect {
			t.Errorf("removeRemainingTags(%q) = %q, want %q", tt.input, result, tt.expect)
		}
	}
}

func TestCleanContent(t *testing.T) {
	input := "  <p>Hello <strong>World</strong></p>  "
	result := cleanContent(input)
	if result != "Hello World" {
		t.Errorf("cleanContent(%q) = %q, want %q", input, result, "Hello World")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
