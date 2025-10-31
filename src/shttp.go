package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//https://modelcontextprotocol.io/specification/2025-06-18/basic/transports#streamable-http
func runSSEServer(config *Config, server *mcp.Server) {
	addr := fmt.Sprintf("%s:%s", config.SSEHost, config.SSEPort)

	log.Printf("Starting MCP server in SSE mode...")
	log.Printf("Server listening on http://%s%s", addr, config.SSEPath)
	log.Printf("Clients can connect to: http://%s%s", addr, config.SSEPath)

	handler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle(config.SSEPath, handler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"name":"%s","version":"%s","transport":"sse","endpoint":"%s"}`,
			config.ServerName, config.ServerVersion, config.SSEPath)
	})

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutdown signal received, gracefully stopping...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	log.Printf("Server listening...")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}

	log.Println("Server stopped.")
}