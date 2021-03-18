package markdown_server

import (
	"bytes"
	"io"

	"github.com/gomarkdown/markdown"
)

func renderMarkdown(input []byte) io.ReadSeeker {
	return bytes.NewReader(markdown.ToHTML(input, nil, nil))
}
