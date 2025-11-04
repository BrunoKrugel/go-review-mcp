# Go Review MCP (WIP)

A Model Context Protocol (MCP) server that provides real-time access to official Go style guides for intelligent code review.

## Overview

Go Review MCP fetches the latest style guides directly from official sources:

- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Best Practices](https://google.github.io/styleguide/go/best-practices)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

The server provides MCP tools and prompts that AI assistants can use to review Go code against these official guidelines.

## Features

### Tools

- **get_style_rule** - Retrieve specific style guide rules by ID
- **search_style_rules** - Search for rules matching a query
- **list_all_rules** - List all available style guide rules
- **get_review_guidelines** - Get comprehensive guidelines for specific topics (naming, errors, concurrency, testing, interfaces)
- **fetch_live_guide** - Fetch and index the latest style guides from official URLs
- **search_live_guide** - Search through live-fetched guide content
- **get_guide_topic** - Extract topic-specific content from live guides

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
go install github.com/bruno-krugel/go-review-mcp@latest
```

### From source

```bash
git clone https://github.com/bruno-krugel/go-review-mcp.git
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

### VS Code with Cline Extension

Add to Cline MCP settings:

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

### Cursor

Add to Cursor's MCP settings (Settings > Features > MCP):

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

1. Start by fetching the latest guides:
   ```
   Use the fetch_live_guide tool with source="all"
   ```

2. Review your code:
   ```
   Use the review_code prompt with your Go code
   ```

3. Get specific guidance:
   ```
   Use get_guide_topic for focused advice on topics like "naming" or "errors"
   ```

### Example Prompts

**For comprehensive review:**
```
Please review this Go code using the review_code prompt
```

**For specific topics:**
```
Check the naming conventions in this code
```

**For searching guidelines:**
```
Search for guidelines about error handling
```

## Configuration Options

The server supports configuration through environment variables:

- `TRANSPORT` - Transport mode: `stdio` (default) or `http`
- `PORT` - HTTP port when using http transport (default: 8080)
- `CACHE_TTL` - Cache duration for fetched guides (default: 24h)
- `LOG_LEVEL` - Logging level: debug, info, warn, error (default: info)

Example:

```bash
CACHE_TTL=12h LOG_LEVEL=debug go-review-mcp
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

1. **Live Fetching**: The server fetches style guides from official URLs and caches them for 24 hours
2. **Content Parsing**: Guides are parsed into structured sections with automatic topic indexing
3. **Embedded Fallback**: Essential rules are embedded for offline use
4. **MCP Integration**: Tools and prompts are exposed via the Model Context Protocol

## Architecture

```
┌─────────────────┐
│   AI Assistant  │
│  (Claude, etc.) │
└────────┬────────┘
         │ MCP Protocol (stdio/http)
┌────────▼──────────────────────┐
│  Go Review MCP Server         │
│  - 7 Tools                    │
│  - 6 Prompts                  │
│  - HTTP Fetcher (24h cache)   │
│  - Content Parser             │
│  - Embedded Rules             │
└────────┬──────────────────────┘
         │ HTTPS
┌────────▼──────────────────────┐
│  Official Style Guide URLs    │
└───────────────────────────────┘
```

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
