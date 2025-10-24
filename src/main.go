package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupLogging() {
	logFileName := "mcp-server.log"
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file %s: %v. Logging to stderr only.", logFileName, err)
	} else {
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
		log.Printf("Logging initialized. Outputting to stderr and %s", logFileName)
	}
}

func main() {
	setupLogging()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "scryfall-card-search-server",
		Version: "v1.0.0"}, nil)

	registerTools(server)

	log.Println("Starting MCP server for MTG card search...")

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped.")
}