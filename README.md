# Go Review MCP

A Model Context Protocol (MCP) server that provides real-time access to official Go style guides for intelligent code review.

## Overview

Go Review MCP fetches the latest style guides directly from official and recommended sources:

### Official

- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Best Practices](https://google.github.io/styleguide/go/best-practices)

### Recommended

- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Features

### Tools

- **search_live_guide** - Search through all live-fetched Go style guide content for specific patterns (auto-fetches guides on first use)
- **get_guide_topic** - Extract topic-specific content from all live guides (naming, errors, concurrency, testing, interfaces, formatting, comments, imports, context)
- **get_review_guidelines** - Get curated review guidelines for specific topics

### Prompts

- **review_code** - Comprehensive Go code review
- **check_naming** - Review naming conventions
- **check_error_handling** - Review error handling patterns
- **check_concurrency** - Review concurrency patterns
- **check_testing** - Review test code quality
- **check_interfaces** - Review interface design

## Installation

### Using go install

```bash
go install github.com/BrunoKrugel/go-review-mcp@latest
```

### From source

```bash
git clone https://github.com/BrunoKrugel/go-review-mcp.git
cd go-review-mcp
make build
```

## Configuration

### Claude Desktop

Add to your `claude_desktop_config.json`:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "go-review": {
      "command": "go-review-mcp",
      "args": []
    }
  }
}
```

If installed from source, use the full path to the binary:

```json
{
  "mcpServers": {
    "go-review": {
      "command": "/path/to/go-review-mcp/bin/go-review-mcp",
      "args": []
    }
  }
}
```

### Cursor or VS Code Extension

```json
{
  "mcpServers": {
    "go-review": {
      "command": "go-review-mcp",
      "args": []
    }
  }
}
```

### Crush

```json
{
  "go-review-mcp": {
    "type": "stdio",
    "command": "go-review-mcp",
    "timeout": 120,
    "disabled": false
  },
}
```

### HTTP Mode (for testing with MCP Inspector)

```bash
# Start server in HTTP mode
go-review-mcp --transport http --port 8080

# Or with environment variable
TRANSPORT=http PORT=8080 go-review-mcp
```

Then use [MCP Inspector](https://github.com/modelcontextprotocol/inspector) to test:

```bash
npx @modelcontextprotocol/inspector http://localhost:8080/mcp
```

## Usage

### Basic Workflow

1. **Review your code**:
   ```
   Use the review_code prompt with your Go code
   ```

2. **Get specific guidance**:
   ```
   Use get_guide_topic for focused advice on topics like "naming" or "errors"
   ```
   Guides are automatically fetched on first use and cached for 24 hours.

3. **Search for patterns**:
   ```
   Use search_live_guide to find specific content across all guides
   ```

### Example Interactions

**Comprehensive code review:**
```
Please review this Go code:
[paste your code]
```

**Check specific topics:**
```
Check the naming conventions in this code
```

**Search for guidelines:**
```
Search the style guides for information about error wrapping
```

**Get topic-specific content:**
```
Show me all guidelines about concurrency from the style guides
```

### Available Topics

The following topics are automatically indexed across all guides:
- `naming` - Naming conventions and identifier rules
- `errors` - Error handling and error wrapping patterns
- `concurrency` - Goroutines, channels, mutexes, context usage
- `testing` - Test structure, table-driven tests, best practices
- `interfaces` - Interface design and usage patterns
- `formatting` - Code formatting and style
- `comments` - Documentation and comment conventions
- `imports` - Import organization and grouping
- `context` - Context.Context usage and patterns

## Configuration Options

The server supports configuration through environment variables:

- `TRANSPORT` - Transport mode: `stdio` (default) or `http`
- `PORT` - HTTP port when using http transport (default: 8080)
- `CACHE_TTL` - Cache duration for fetched guides (default: 24h)

Example:

```bash
CACHE_TTL=12h go-review-mcp
```

## Development

### Building

```bash
make build
```

### Running tests

```bash
make test
```

### Running locally

```bash
make run
```

## How It Works

1. **Lazy Fetching**: Style guides are automatically fetched on first tool use (no manual initialization needed)
2. **Content Parsing**: Guides are parsed into structured sections with automatic topic indexing
3. **Unified Search**: All search and topic queries automatically scan across all fetched guides
4. **MCP Integration**: Simple tools and prompts are exposed via the Model Context Protocol

### Key Features

- **Auto-Fetch on Demand**: Guides are fetched lazily on first use — no manual fetch step required
- **Concurrent Fetching**: All guides are fetched in parallel for fast initialization
- **Smart Caching**: 24-hour cache reduces network requests and improves performance
- **Topic Indexing**: Automatic categorization of content by topics (naming, errors, concurrency, etc.)
- **No Configuration**: Works out of the box - no need to choose which guides to use

## Architecture

```
┌─────────────────┐
│   AI Assistant  │
│  (Claude, etc.) │
└────────┬────────┘
         │ MCP Protocol (stdio/http)
┌────────▼──────────────────────┐
│  Go Review MCP Server         │
│  - 3 Tools                    │
│  - 6 Prompts                  │
│  - Concurrent HTTP Fetcher    │
│  - 24h Cache (in-memory)      │
│  - Content Parser             │
│  - Topic Indexer              │
└────────┬──────────────────────┘
         │ HTTPS (parallel)
┌────────▼──────────────────────┐
│  Official Style Guide Sources │
│  ├─ Google Go Guide           │
│  ├─ Google Go Decisions       │
│  ├─ Google Go Best Practices  │
│  └─ Uber Go Style Guide       │
└───────────────────────────────┘
```

### Data Flow

1. **On-Demand Fetch**: First tool call triggers concurrent fetch of all 4 guides
2. **Parsing**: Each guide is parsed into sections with metadata
3. **Indexing**: Sections are automatically categorized by topics
4. **Query**: Search and topic tools scan all indexed content
5. **Caching**: Results are cached for 24 hours to improve performance

## Requirements

- Go 1.21 or later
- Internet connection for fetching live guides (embedded rules work offline)

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

This tool aggregates content from:
- [Google Go Style Guide](https://google.github.io/styleguide/go/) - CC-BY-3.0
- [Uber Go Style Guide](https://github.com/uber-go/guide) - MIT
