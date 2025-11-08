package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML in markdown
		),
	)
}

// RenderMarkdown converts markdown to HTML
func RenderMarkdown(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
