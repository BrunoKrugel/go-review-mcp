package prompts

import "testing"

func TestGetAllPrompts(t *testing.T) {
	all := GetAllPrompts()
	if len(all) != 6 {
		t.Fatalf("expected 6 prompts, got %d", len(all))
	}

	names := map[string]bool{
		"review_code":          false,
		"check_naming":         false,
		"check_error_handling": false,
		"check_concurrency":    false,
		"check_testing":        false,
		"check_interfaces":     false,
	}

	for _, p := range all {
		if _, ok := names[p.Name]; !ok {
			t.Errorf("unexpected prompt name: %q", p.Name)
		}
		names[p.Name] = true
		if p.Template == "" {
			t.Errorf("prompt %q has empty template", p.Name)
		}
		if p.Description == "" {
			t.Errorf("prompt %q has empty description", p.Name)
		}
	}

	for name, found := range names {
		if !found {
			t.Errorf("missing prompt: %q", name)
		}
	}
}

func TestPromptsContainGoLanguage(t *testing.T) {
	all := GetAllPrompts()
	for _, p := range all {
		if !containsGoRef(p.Description) {
			t.Errorf("prompt %q description should mention Go language: %q", p.Name, p.Description)
		}
	}
}

func TestPromptTemplates_HavePlaceholders(t *testing.T) {
	all := GetAllPrompts()
	for _, p := range all {
		hasCode := len(p.Arguments) > 0
		if !hasCode {
			continue
		}
		found := false
		for _, arg := range p.Arguments {
			if arg.Name == "code" && arg.Required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("prompt %q should have a required 'code' argument", p.Name)
		}
	}
}

func containsGoRef(s string) bool {
	return len(s) >= 2 && (containsStr(s, "Go") || containsStr(s, "Golang"))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
