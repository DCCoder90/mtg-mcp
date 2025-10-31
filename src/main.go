package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupLogging(config *Config) {
	if !config.LogToFile {
		log.SetOutput(os.Stderr)
		log.Println("Logging to stderr only (file logging disabled)")
		return
	}

	logFile, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file %s: %v. Logging to stderr only.", config.LogFilePath, err)
	} else {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
		log.Printf("Logging initialized. Outputting to stderr and %s", config.LogFilePath)
	}
}

func main() {
	config := LoadConfig()
	setupLogging(config)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    config.ServerName,
		Version: config.ServerVersion}, nil)

	registerTools(server)

	switch config.Transport {
	case TransportStdio:
		runStdioServer(server)
	case TransportSSE:
		runSSEServer(config, server)
	default:
		log.Fatalf("Unknown transport type: %s", config.Transport)
	}
}

func runStdioServer(server *mcp.Server) {
	log.Println("Starting MCP server in STDIO mode for local execution...")

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped.")
}

