package ast

import (
	"fmt"
	"go/ast"
	"strings"
)

type Annotation struct {
	key   string
	value string
}

type annotations struct {
	Node ast.Node
	Ans  []Annotation
}

func newAnnotations(node ast.Node, commentGroup *ast.CommentGroup) annotations {
	if commentGroup == nil || len(commentGroup.List) == 0 {
		return annotations{Node: node}
	}
	ans := make([]Annotation, 0, len(commentGroup.List))

	for _, comment := range commentGroup.List {
		text, ok := extractContent(comment)
		if !ok {
			continue
		}
		if strings.HasPrefix(text, "@") {
			segues := strings.SplitN(text, " ", 2)
			if len(segues) != 2 {
				continue
			}
			key := segues[0][1:]
			value := segues[1]
			ans = append(ans, Annotation{key, value})
		}
	}
	return annotations{
		Node: node,
		Ans:  ans,
	}
}

func extractContent(comment *ast.Comment) (string, bool) {
	text := comment.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	}
	if strings.HasPrefix(text, "/* ") {
		fmt.Println(text)
		return text[3: len(text)-2], true
	}
	return "", false
}
