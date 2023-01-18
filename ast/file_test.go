package ast

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestFileVisitorGet(t *testing.T) {
	testCases := []struct {
		src   string
		wants File
	}{
		{src: srcStr, wants: wants},
	}

	for _, tc := range testCases {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tc.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		topVisitor := &SingleFileEntryVisitor{}
		ast.Walk(topVisitor, f)
		file := topVisitor.Get()
		assertAnnotations(t, tc.wants.annotations, file.annotations)

		for i, typ := range file.types {
			wantTyp := tc.wants.types[i]
			assertAnnotations(t, wantTyp.annotations, typ.annotations)

			if len(wantTyp.Fields) != len(typ.Fields) {
				t.Fatal()
			}
			for j, field := range wantTyp.Fields {
				assertAnnotations(t, field.annotations, typ.Fields[j].annotations)
			}
		}
	}
}

func TestStrHasPrefix(t *testing.T) {
	str := `/* @multiple first line
second line
*/`
	assert.True(t, strings.HasPrefix(str, "/* "))
}

func assertAnnotations(t *testing.T, wantAns, dst annotations) {

	if len(wantAns.Ans) != len(dst.Ans) {
		t.Fatal()
	}
	for i, ans := range wantAns.Ans {
		assert.Equal(t, ans, dst.Ans[i])
	}
}

var srcStr = `
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
	second line
*/
// @date 2022/04/02
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
		second line
	*/
	// @date 2022/04/02
	FuncType func()
)

type (
	// StructType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
		second line
	*/
	// @date 2022/04/02
	StructType struct {
		// Public is a field
		// @type string
		Public string
	}

	// SecondType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
		second line
	*/
	// @date 2022/04/03
	SecondType struct {
	}
)

type (
	// Interface is a test interface
	// @author Deng Ming
	/* @multiple first line
		second line
	*/
	// @date 2022/04/02
	Interface interface {
		// MyFunc is a test func
		// @parameter arg1 int
		// @parameter arg2 int32
		// @return string
		MyFunc(arg1 int, arg2 int32) string

		// second is a test func
		// @return string
		second() string
	}
)
`

var wants = File{
	annotations: annotations{
		Ans: []Annotation{
			{
				key:   "author",
				value: "Deng Ming",
			},
			{
				key:   "multiple",
				value: "first line\n\tsecond line\n",
			},
			{
				key:   "date",
				value: "2022/04/02",
			},
		},
	},
	types: []Type{
		{
			annotations: annotations{
				Ans: []Annotation{
					{
						key:   "author",
						value: "Deng Ming",
					},
					{
						key:   "multiple",
						value: "first line\n\t\tsecond line\n\t",
					},
					{
						key:   "date",
						value: "2022/04/02",
					},
				},
			},
		},
		{
			annotations: annotations{
				Ans: []Annotation{
					{
						key:   "author",
						value: "Deng Ming",
					},
					{
						key:   "multiple",
						value: "first line\n\t	second line\n\t",
					},
					{
						key:   "date",
						value: "2022/04/02",
					},
				},
			},
			Fields: []Field{
				{
					annotations: annotations{
						Ans: []Annotation{
							{
								key:   "type",
								value: "string",
							},
						},
					},
				},
			},
		},
		{
			annotations: annotations{
				Ans: []Annotation{
					{
						key:   "author",
						value: "Deng Ming",
					},
					{
						key:   "multiple",
						value: "first line\n\t	second line\n\t",
					},
					{
						key:   "date",
						value: "2022/04/03",
					},
				},
			},
		},
		{
			annotations: annotations{
				Ans: []Annotation{
					{
						key:   "author",
						value: "Deng Ming",
					},
					{
						key:   "multiple",
						value: "first line\n\t	second line\n\t",
					},
					{
						key:   "date",
						value: "2022/04/02",
					},
				},
			},
			Fields: []Field{
				{
					annotations: annotations{
						Ans: []Annotation{
							{
								key:   "parameter",
								value: "arg1 int",
							},
							{
								key:   "parameter",
								value: "arg2 int32",
							},
							{
								key:   "return",
								value: "string",
							},
						},
					},
				},
				{
					annotations:annotations{
						Ans: []Annotation{
							{
								key:   "return",
								value: "string",
							},
						},
					},
				},
			},
		},
	},
}
