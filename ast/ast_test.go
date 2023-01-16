package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var srcContent = `
package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)
// printVisitor 访问器模式
type printVisitor struct {
}

func (t *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println(nil)
		return t
	}
	val := reflect.ValueOf(node)
	typ := reflect.TypeOf(node)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	fmt.Printf("val: %+v, type: %s \n", val.Interface(), typ.Name())
	return t
}

`

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", srcContent, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ast.Walk(&printVisitor{}, f)
}
