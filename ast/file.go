package ast

import "go/ast"

type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	if s.file != nil {
		return s.file.Get()
	}
	return File{}
}

func (sv *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	file, ok := node.(*ast.File)
	if !ok {
		return sv
	}

	sv.file = &fileVisitor{
		ans: newAnnotations(node, file.Doc),
	}
	return sv.file
}

type fileVisitor struct {
	ans          annotations
	typeVisitors []*typeVisitor
}

func (fv *fileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	spec, ok := node.(*ast.TypeSpec)
	if !ok {
		return fv
	}
	typV := &typeVisitor{
		ans:    newAnnotations(spec, spec.Doc),
		fields: make([]Field, 0, 0),
	}
	fv.typeVisitors = append(fv.typeVisitors, typV)
	return typV
}

func (fv *fileVisitor) Get() File {
	types := make([]Type, 0, len(fv.typeVisitors))
	for _, tv := range fv.typeVisitors {
		types = append(types, tv.Get())
	}
	return File{
		annotations: fv.ans,
		types:       types,
	}
}

type typeVisitor struct {
	ans    annotations
	fields []Field
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {

	fd, ok := node.(*ast.Field)
	if !ok{
		return t
	}
	t.fields = append(t.fields, Field{
		annotations: newAnnotations(fd, fd.Doc),
	})
	return nil
}

func (t *typeVisitor) Get() Type {
	return Type{
		annotations: t.ans,
		Fields:      t.fields,
	}
}

type Field struct {
	annotations
}

type Type struct {
	annotations
	Fields []Field
}

type File struct {
	annotations
	types []Type
}
