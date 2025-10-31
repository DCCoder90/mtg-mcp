package main

import (
	"embed"
)

// Embed all resource files from the res directory
// https://pkg.go.dev/embed
//go:embed res/*
var embeddedResources embed.FS
