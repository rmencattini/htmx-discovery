package templates

import "embed"

// Embed the entire directory.
//go:embed *
var TemplatesFolder embed.FS
