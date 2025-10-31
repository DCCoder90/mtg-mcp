package main

import (
	"os"
	"strconv"
)

type TransportType string

const (
	TransportStdio TransportType = "stdio"
	TransportSSE   TransportType = "sse"
)

type Config struct {
	ServerName    string
	ServerVersion string
	LogToFile     bool
	LogFilePath   string
	Transport     TransportType
	SSEHost       string
	SSEPort       string
	SSEPath       string
	SSLCertFile   string
	SSLKeyFile    string
}

func LoadConfig() *Config {
	logToFile := true
	if val := os.Getenv("MCP_LOG_TO_FILE"); val != "" {
		logToFile, _ = strconv.ParseBool(val)
	}

	SSLCertFile := ""
	if val := os.Getenv("MCP_SSL_CERT_FILE"); val != "" {
		SSLCertFile = val
	}

	SSLKeyFile := ""
	if val := os.Getenv("MCP_SSL_KEY_FILE"); val != "" {
		SSLKeyFile = val
	}

	logFilePath := "mcp-server.log"
	if val := os.Getenv("MCP_LOG_FILE"); val != "" {
		logFilePath = val
	}

	serverName := "scryfall-card-search-server"
	if val := os.Getenv("MCP_SERVER_NAME"); val != "" {
		serverName = val
	}

	serverVersion := "v1.0.0"
	if val := os.Getenv("MCP_SERVER_VERSION"); val != "" {
		serverVersion = val
	}

	transport := TransportStdio
	if val := os.Getenv("MCP_TRANSPORT"); val != "" {
		if val == "sse" || val == "SSE" {
			transport = TransportSSE
		}
	}

	sseHost := "0.0.0.0"
	if val := os.Getenv("MCP_SSE_HOST"); val != "" {
		sseHost = val
	}

	ssePort := "3000"
	if val := os.Getenv("MCP_SSE_PORT"); val != "" {
		ssePort = val
	}

	ssePath := "/sse"
	if val := os.Getenv("MCP_SSE_PATH"); val != "" {
		ssePath = val
	}

	return &Config{
		ServerName:    serverName,
		ServerVersion: serverVersion,
		LogToFile:     logToFile,
		LogFilePath:   logFilePath,
		Transport:     transport,
		SSEHost:       sseHost,
		SSEPort:       ssePort,
		SSEPath:       ssePath,
		SSLCertFile:   SSLCertFile,
		SSLKeyFile:    SSLKeyFile,
	}
}
