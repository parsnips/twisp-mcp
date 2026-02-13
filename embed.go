package main

import "embed"

//go:embed all:embed/schema all:embed/docs all:embed/examples
var embeddedContent embed.FS
