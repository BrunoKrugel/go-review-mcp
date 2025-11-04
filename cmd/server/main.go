package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BrunoKrugel/go-review-mcp/internal/config"
	"github.com/BrunoKrugel/go-review-mcp/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command-line flags
	transport := flag.String("transport", "", "Transport type (stdio or http)")
	port := flag.Int("port", 0, "HTTP port (when transport=http)")
	flag.Parse()

	// Load configuration
	cfg := config.LoadFromEnv()

	// Override with command-line flags if provided
	if *transport != "" {
		cfg.Transport = *transport
	}
	if *port != 0 {
		cfg.Port = *port
	}

	log.SetOutput(os.Stderr)
	log.Printf("Starting Go Review MCP Server (transport=%s)...", cfg.Transport)

	mcpServer := mcp.NewServerWithConfig(cfg)

	switch cfg.Transport {
	case "http":
		addr := fmt.Sprintf(":%d", cfg.Port)
		log.Printf("HTTP server listening on %s", addr)
		httpServer := server.NewStreamableHTTPServer(mcpServer.MCPServer)
		if err := httpServer.Start(addr); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	case "stdio":
		if err := server.ServeStdio(mcpServer.MCPServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf("Unknown transport: %s (use 'stdio' or 'http')", cfg.Transport)
	}
}
