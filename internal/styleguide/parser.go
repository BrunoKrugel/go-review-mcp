package styleguide

import (
	"strings"
)

// ContentIndex represents indexed style guide content
type ContentIndex struct {
	Content  string
	Topics   map[string][]Section
	Sections []Section
}

// Section represents a section in the style guide
type Section struct {
	Title   string
	Content string
	Anchor  string
	Level   int
}

// ParseContent parses style guide content into structured sections
func ParseContent(content string) *ContentIndex {
	index := &ContentIndex{
		Content: content,
		Topics:  make(map[string][]Section),
	}

	lines := strings.Split(content, "\n")
	var currentSection *Section
	var sectionContent strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			if currentSection != nil {
				currentSection.Content = strings.TrimSpace(sectionContent.String())
				index.Sections = append(index.Sections, *currentSection)
			}

			level := 0
			for i := 0; i < len(line) && line[i] == '#'; i++ {
				level++
			}

			title := strings.TrimSpace(line[level:])
			anchor := generateAnchor(title)

			currentSection = &Section{
				Title:  title,
				Level:  level,
				Anchor: anchor,
			}
			sectionContent.Reset()
		} else if currentSection != nil {
			sectionContent.WriteString(line)
			sectionContent.WriteString("\n")
		}
	}

	if currentSection != nil {
		currentSection.Content = cleanContent(sectionContent.String())
		index.Sections = append(index.Sections, *currentSection)
	}

	// Clean content for all sections
	for i := range index.Sections {
		index.Sections[i].Content = cleanContent(index.Sections[i].Content)
	}

	index.indexTopics()
	return index
}

func cleanContent(content string) string {
	content = strings.TrimSpace(content)

	// Remove HTML tags
	content = removeHTMLTags(content)

	return content
}

func removeHTMLTags(content string) string {
	// Common HTML tags to remove
	htmlTags := []string{
		"<table>", "</table>",
		"<thead>", "</thead>",
		"<tbody>", "</tbody>",
		"<tr>", "</tr>",
		"<td>", "</td>",
		"<th>", "</th>",
		"<br>", "<br/>",
		"<hr>", "<hr/>",
		"<div>", "</div>",
		"<span>", "</span>",
		"<p>", "</p>",
		"<ul>", "</ul>",
		"<ol>", "</ol>",
		"<li>", "</li>",
		"<a>", "</a>",
		"<strong>", "</strong>",
		"<em>", "</em>",
		"<b>", "</b>",
		"<i>", "</i>",
	}

	result := content
	for _, tag := range htmlTags {
		result = strings.ReplaceAll(result, tag, "")
	}

	// Remove HTML attributes from remaining tags using regex-like approach
	// Simple approach: remove anything that looks like <tagname ...> or </tagname>
	result = removeRemainingTags(result)

	return result
}

func removeRemainingTags(content string) string {
	var result strings.Builder
	inTag := false

	//nolint:intrange // Need index for byte access
	for i := 0; i < len(content); i++ {
		if content[i] == '<' {
			inTag = true
			continue
		}
		if content[i] == '>' && inTag {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteByte(content[i])
		}
	}

	return result.String()
}

func (idx *ContentIndex) indexTopics() {
	keywords := getTopicKeywords()
	assigned := make(map[int]bool)

	// First pass: assign sections based on title keywords (most specific)
	idx.assignByTitle(keywords, assigned)

	// Second pass: assign remaining sections based on content keywords
	idx.assignByContent(keywords, assigned)
}

func getTopicKeywords() map[string][]string {
	return map[string][]string{
		"naming":      {"naming", "identifier", "variable name", "function name", "type name", "package name", "receiver name", "name convention"},
		"errors":      {"error", "error handling", "error wrapping", "error return", "panic", "recover"},
		"concurrency": {"concurrency", "concurrent", "goroutine", "channel", "mutex", "sync.", "parallel", "race condition", "waitgroup"},
		"testing":     {"testing", "test", "test function", "benchmark", "test case", "table-driven"},
		"interfaces":  {"interface", "interface design", "interface pollution", "accept interfaces"},
		"formatting":  {"formatting", "gofmt", "goimports", "code style", "indentation", "format"},
		"comments":    {"comment", "documentation", "godoc", "doc comment"},
		"imports":     {"import", "import statement", "import grouping", "import path"},
		"context":     {"context", "context.context", "context package", "ctx", "cancellation", "deadline", "context value"},
	}
}

func (idx *ContentIndex) assignByTitle(keywords map[string][]string, assigned map[int]bool) {
	for i, section := range idx.Sections {
		titleLower := strings.ToLower(section.Title)
		bestMatch, bestScore := findBestMatch(titleLower, keywords, true)

		if bestMatch != "" && bestScore > 0 {
			idx.Topics[bestMatch] = append(idx.Topics[bestMatch], section)
			assigned[i] = true
		}
	}
}

func (idx *ContentIndex) assignByContent(keywords map[string][]string, assigned map[int]bool) {
	for i, section := range idx.Sections {
		if assigned[i] {
			continue
		}

		contentLower := strings.ToLower(section.Content)
		bestMatch, bestScore := findBestMatch(contentLower, keywords, false)

		// Only assign if there's a reasonable match (at least 2 keyword hits)
		if bestMatch != "" && bestScore >= 2 {
			idx.Topics[bestMatch] = append(idx.Topics[bestMatch], section)
		}
	}
}

func findBestMatch(text string, keywords map[string][]string, scoreByLength bool) (string, int) {
	bestMatch := ""
	bestScore := 0

	for topic, words := range keywords {
		score := 0
		for _, word := range words {
			if strings.Contains(text, word) {
				if scoreByLength {
					score += len(word) // Longer matches = better score
				} else {
					score++
				}
			}
		}
		if score > bestScore {
			bestScore = score
			bestMatch = topic
		}
	}

	return bestMatch, bestScore
}

// SearchContent searches for specific content in the style guide
func (idx *ContentIndex) SearchContent(query string) []Section {
	queryLower := strings.ToLower(query)
	var results []Section

	for _, section := range idx.Sections {
		if strings.Contains(strings.ToLower(section.Title), queryLower) ||
			strings.Contains(strings.ToLower(section.Content), queryLower) {
			results = append(results, section)
		}
	}

	return results
}

// GetTopic returns all sections related to a topic
func (idx *ContentIndex) GetTopic(topic string) []Section {
	return idx.Topics[strings.ToLower(topic)]
}

// GetSection returns a section by title
func (idx *ContentIndex) GetSection(title string) *Section {
	for _, section := range idx.Sections {
		if strings.EqualFold(section.Title, title) {
			return &section
		}
	}
	return nil
}

func generateAnchor(title string) string {
	anchor := strings.ToLower(title)
	anchor = strings.ReplaceAll(anchor, " ", "-")
	anchor = strings.ReplaceAll(anchor, ".", "")
	anchor = strings.ReplaceAll(anchor, ",", "")
	anchor = strings.ReplaceAll(anchor, ":", "")
	anchor = strings.ReplaceAll(anchor, "'", "")
	anchor = strings.ReplaceAll(anchor, "\"", "")
	anchor = strings.ReplaceAll(anchor, "(", "")
	anchor = strings.ReplaceAll(anchor, ")", "")
	return anchor
}

// ExtractCodeExamples extracts code examples from content
func ExtractCodeExamples(content string) []string {
	var examples []string
	lines := strings.Split(content, "\n")
	var inCodeBlock bool
	var currentExample strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				examples = append(examples, currentExample.String())
				currentExample.Reset()
				inCodeBlock = false
			} else {
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			currentExample.WriteString(line)
			currentExample.WriteString("\n")
		}
	}

	return examples
}
