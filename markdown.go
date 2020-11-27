package markdown_server

import (
	"bytes"
	"github.com/gomarkdown/markdown"
	"io"
)

func renderMarkdown(input []byte) io.ReadSeeker {
	return bytes.NewReader(markdown.ToHTML(input, nil, nil))
}