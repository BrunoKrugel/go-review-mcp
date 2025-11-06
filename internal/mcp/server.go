package mcp

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/BrunoKrugel/go-review-mcp/internal/config"
	"github.com/BrunoKrugel/go-review-mcp/internal/prompts"
	"github.com/BrunoKrugel/go-review-mcp/internal/styleguide"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	sourceAll = "all"

	// maxContentPreview is the character limit for content previews in search results
	// to keep response sizes manageable while providing useful context
	maxContentPreview = 500
)

// Server wraps the MCP server with style guide data
type Server struct {
	*server.MCPServer
	fetcher   *styleguide.Fetcher
	indices   map[string]*styleguide.ContentIndex
	indicesMu sync.RWMutex
}

// NewServer creates a new MCP server with default configuration
func NewServer() *Server {
	return NewServerWithConfig(&config.Config{
		CacheTTL: config.DefaultCacheTTL,
	})
}

// NewServerWithConfig creates a new MCP server with custom configuration
func NewServerWithConfig(cfg *config.Config) *Server {
	mcpServer := server.NewMCPServer(
		"go-review-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	s := &Server{
		MCPServer: mcpServer,
		fetcher:   styleguide.NewFetcher(cfg.CacheTTL),
		indices:   make(map[string]*styleguide.ContentIndex),
	}

	s.registerTools()
	s.registerPrompts()

	return s
}

func (s *Server) registerTools() {
	// Removed embedded guide tools (get_style_rule, search_style_rules, list_all_rules)
	// All operations now use live-fetched guides only

	s.AddTool(mcp.Tool{
		Name:        "fetch_live_guide",
		Description: "Fetch and parse the latest style guide content from all official URLs (Google and Uber). Call this first before using other tools.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}, s.handleFetchLiveGuide)

	s.AddTool(mcp.Tool{
		Name:        "search_live_guide",
		Description: "Search through all live-fetched style guides for specific topics or patterns",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Search query to find in guide content",
				},
			},
			Required: []string{"query"},
		},
	}, s.handleSearchLiveGuide)

	s.AddTool(mcp.Tool{
		Name:        "get_guide_topic",
		Description: "Get all content related to a specific topic from all live-fetched guides",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"topic": map[string]any{
					"type":        "string",
					"description": "Topic to retrieve: 'naming', 'errors', 'concurrency', 'testing', 'interfaces', 'formatting', 'comments', 'imports', 'context'",
					"enum":        []string{"naming", "errors", "concurrency", "testing", "interfaces", "formatting", "comments", "imports", "context"},
				},
			},
			Required: []string{"topic"},
		},
	}, s.handleGetGuideTopic)

	s.AddTool(mcp.Tool{
		Name:        "get_review_guidelines",
		Description: "Get curated review guidelines covering major topics",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"topic": map[string]any{
					"type":        "string",
					"description": "Specific topic: 'naming', 'errors', 'concurrency', 'testing', 'interfaces', 'all'",
					"enum":        []string{"naming", "errors", "concurrency", "testing", "interfaces", sourceAll},
				},
			},
			Required: []string{"topic"},
		},
	}, s.handleGetReviewGuidelines)
}

func (s *Server) registerPrompts() {
	allPrompts := prompts.GetAllPrompts()

	for _, p := range allPrompts {
		args := make([]mcp.PromptArgument, len(p.Arguments))
		for i, arg := range p.Arguments {
			args[i] = mcp.PromptArgument{
				Name:        arg.Name,
				Description: arg.Description,
				Required:    arg.Required,
			}
		}

		s.AddPrompt(mcp.Prompt{
			Name:        p.Name,
			Description: p.Description,
			Arguments:   args,
		}, createPromptHandler(p))
	}
}

func createPromptHandler(p prompts.Prompt) func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		// Context intentionally unused - this is pure string manipulation
		template := p.Template

		if request.Params.Arguments != nil {
			for k, v := range request.Params.Arguments {
				placeholder := fmt.Sprintf("{{%s}}", k)
				template = strings.ReplaceAll(template, placeholder, v)
			}
		}

		template = strings.ReplaceAll(template, "{{if focus}}", "")
		template = strings.ReplaceAll(template, "{{end}}", "")

		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: template,
					},
				},
			},
		}, nil
	}
}

func (s *Server) handleGetReviewGuidelines(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic := mcp.ParseString(request, "topic", "")

	if topic == "" {
		return mcp.NewToolResultError("missing required parameter: topic"), nil
	}

	output := "# Go Code Review Guidelines\n\n"

	topics := []string{topic}
	if topic == sourceAll {
		topics = []string{"naming", "errors", "concurrency", "testing", "interfaces"}
	}

	for _, t := range topics {
		switch t {
		case "naming":
			output += s.getNamingGuidelines()
		case "errors":
			output += s.getErrorGuidelines()
		case "concurrency":
			output += s.getConcurrencyGuidelines()
		case "testing":
			output += s.getTestingGuidelines()
		case "interfaces":
			output += s.getInterfaceGuidelines()
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown topic: %s", t)), nil
		}
		output += "\n"
	}

	return mcp.NewToolResultText(output), nil
}

func (s *Server) handleFetchLiveGuide(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Always fetch all guides
	guides, err := s.fetcher.FetchAll(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to fetch guides: %v", err)), nil
	}

	s.indicesMu.Lock()
	for name, content := range guides {
		s.indices[name] = styleguide.ParseContent(content)
	}
	s.indicesMu.Unlock()

	output := "Successfully fetched and indexed all style guides:\n"
	for name := range guides {
		output += fmt.Sprintf("- %s\n", name)
	}
	return mcp.NewToolResultText(output), nil
}

func (s *Server) handleSearchLiveGuide(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := mcp.ParseString(request, "query", "")

	s.indicesMu.RLock()
	indicesCount := len(s.indices)
	s.indicesMu.RUnlock()

	if indicesCount == 0 {
		return mcp.NewToolResultError("no guides fetched yet. Use fetch_live_guide first"), nil
	}

	// Always search all guides
	var allResults []styleguide.Section

	s.indicesMu.RLock()
	for _, index := range s.indices {
		results := index.SearchContent(query)
		allResults = append(allResults, results...)
	}
	s.indicesMu.RUnlock()

	if len(allResults) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No results found for query: %s", query)), nil
	}

	output := fmt.Sprintf("Found %d sections matching '%s':\n\n", len(allResults), query)
	for _, section := range allResults {
		output += fmt.Sprintf("## %s\n\n%s\n\n---\n\n", section.Title, truncateContent(section.Content, maxContentPreview))
	}

	return mcp.NewToolResultText(output), nil
}

func (s *Server) handleGetGuideTopic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic := mcp.ParseString(request, "topic", "")

	if topic == "" {
		return mcp.NewToolResultError("missing required parameter: topic"), nil
	}

	s.indicesMu.RLock()
	indicesCount := len(s.indices)
	s.indicesMu.RUnlock()

	if indicesCount == 0 {
		return mcp.NewToolResultError("no guides fetched yet. Use fetch_live_guide first"), nil
	}

	// Always search all guides
	var allSections []styleguide.Section

	s.indicesMu.RLock()
	for _, index := range s.indices {
		sections := index.GetTopic(topic)
		allSections = append(allSections, sections...)
	}
	s.indicesMu.RUnlock()

	if len(allSections) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No content found for topic: %s", topic)), nil
	}

	output := fmt.Sprintf("# %s\n\nFound %d sections:\n\n", cases.Title(language.English).String(topic), len(allSections))
	for _, section := range allSections {
		output += fmt.Sprintf("## %s\n\n%s\n\n---\n\n", section.Title, section.Content)
	}

	return mcp.NewToolResultText(output), nil
}

func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

func (s *Server) getNamingGuidelines() string {
	return `## Naming Conventions

### Key Rules
- Use MixedCaps or mixedCaps (no underscores)
- Receiver names: 1-2 characters, consistent
- Package names: lowercase, single word, no underscores
- Exported names start with uppercase
- Interface names end in "-er" for single methods

### Examples
✓ userConfig, LoadSettings, parseURL
✗ user_config, load_settings, parse_u_r_l

### References
- Google Guide: https://google.github.io/styleguide/go/guide#naming
`
}

func (s *Server) getErrorGuidelines() string {
	return `## Error Handling

### Key Rules
- Always check errors explicitly (no _ discards)
- Return errors as last return value
- Wrap errors with %w for context
- Use errors.Is and errors.As for checking
- Error messages: lowercase, no punctuation

### Examples
✓ if err != nil { return fmt.Errorf("failed: %w", err) }
✗ data, _ := readFile()

### References
- Google Decisions: https://google.github.io/styleguide/go/decisions#errors
- Uber Guide: https://github.com/uber-go/guide/blob/master/style.md#error-wrapping
`
}

func (s *Server) getConcurrencyGuidelines() string {
	return `## Concurrency

### Key Rules
- Use context.Context as first parameter
- Don't store context in structs
- Channels: unbuffered or size 1
- Zero-value mutexes (not pointers)
- Always consider goroutine cleanup

### Examples
✓ func Process(ctx context.Context, data Data) error
✓ mu sync.Mutex
✗ type Handler struct { ctx context.Context }
✗ mu *sync.Mutex

### References
- Google Decisions: https://google.github.io/styleguide/go/decisions#contexts
- Uber Guide: https://github.com/uber-go/guide/blob/master/style.md#channel-size-is-one-or-none
`
}

func (s *Server) getTestingGuidelines() string {
	return `## Testing

### Key Rules
- Use table-driven tests
- Test names: TestXxx, clear descriptions
- Use t.Run for subtests
- Test helpers use t.Helper()
- Cover happy and error paths

### Examples
✓ Table-driven tests with t.Run
✓ Clear test case names
✗ Monolithic tests without subtests
✗ Tests depending on execution order

### References
- Google Best Practices: https://google.github.io/styleguide/go/best-practices#tests
`
}

func (s *Server) getInterfaceGuidelines() string {
	return `## Interfaces

### Key Rules
- Accept interfaces, return concrete types
- Small interfaces (1-3 methods)
- Define at point of use
- No pointers to interfaces
- Use var _ Interface = (*Type)(nil) for compliance

### Examples
✓ func Process(r io.Reader) error
✓ type Reader interface { Read([]byte) (int, error) }
✗ func Process(r *io.Reader) error
✗ Large interfaces with many methods

### References
- Google Best Practices: https://google.github.io/styleguide/go/best-practices#interfaces
- Uber Guide: https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces
`
}
