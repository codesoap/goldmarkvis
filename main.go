package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

type visualizer struct{}

func (v *visualizer) AddOptions(...renderer.Option) {}
func (v *visualizer) Render(w io.Writer, source []byte, n ast.Node) error {
	return renderWithIndent(w, source, n, 0)
}

func renderWithIndent(w io.Writer, source []byte, n ast.Node, indent int) error {
	fmt.Fprintln(w, strings.Repeat("\t", indent), toNodeSummary(source, n))
	if n.HasChildren() {
		siblingCount := n.ChildCount()
		for child := n.FirstChild(); siblingCount > 0; siblingCount-- {
			renderWithIndent(w, source, child, indent+1)
			child = child.NextSibling()
		}
	}
	return nil
}

func toNodeSummary(source []byte, n ast.Node) string {
	if text := string(n.Text(source)); text != "" {
		return fmt.Sprintf("Type: %s, Kind: %s, Text: %s", toNodeTypeName(n.Type()), n.Kind().String(), text)
	}
	return fmt.Sprintf("Type: %s, Kind: %s, (no text)", toNodeTypeName(n.Type()), n.Kind().String())
}

func toNodeTypeName(t ast.NodeType) string {
	switch t {
	case ast.TypeBlock:
		return "Block"
	case ast.TypeInline:
		return "Inline"
	case ast.TypeDocument:
		return "Doc"
	}
	return "?"
}

func main() {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	VisualizeMarkdown(in)
}

func VisualizeMarkdown(in []byte) {
	md := goldmark.New(goldmark.WithRenderer(&visualizer{}))
	if err := md.Convert(in, os.Stdout); err != nil {
		panic(err)
	}
}
