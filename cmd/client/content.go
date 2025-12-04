package main

import "embed"

// content holds web server files

//go:embed templates styles
var content embed.FS
