package main

import "embed"

// content holds web server files

//go:embed templates
var content embed.FS
