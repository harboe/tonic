package docs

import (
	"go/ast"
	"go/token"
	"strings"
)

type (
	namedType struct {
		name, desc string
		local      string
		typeName   string
		isTonic    bool // is a tonic type
		isImported bool // is imported from external
		isMethod   bool // orginal from a method call
		isStruct   bool // is a struct
		isEmbedded bool // only accessable from within the function
	}
	namedTypeList map[string]*namedType
	funcType      struct {
		name       string
		desc       []string
		local      string         // import namespace.
		vars       *namedTypeList // local variables
		funcs      *funcTypeList  // function calls.
		isImported bool           // is imported from external

		params []Parameter
	}
	funcTypeList map[string]*funcType
	commentMap   struct {
		fset     *token.FileSet
		comments map[int]string
	}
)

func (fun *funcType) setDoc(doc string) {
	// log.Println("doc:", doc)
	fun.name, fun.desc = parseComment(doc)
}

func (fun *funcType) lookup(name string) (*namedType, *funcType) {
	funDic := *fun.funcs
	if v, ok := funDic[name]; ok {
		return nil, v
	}

	varDic := *fun.vars
	if v, ok := varDic[name]; ok {
		return v, nil
	}

	return nil, nil
}

func (ntl *namedTypeList) lookup(name string) *namedType {
	dic := *ntl
	if name == "" || name == "_" {
		return nil // no type docs for anoymous types
	}

	if typ, ok := dic[name]; ok {
		return typ
	}

	typ := &namedType{
		name: name,
	}
	dic[name] = typ
	ntl = &dic
	return typ
}

func (ftl *funcTypeList) lookup(name string) *funcType {
	dic := *ftl
	if name == "" || name == "_" {
		return nil // no type docs for anoymous types
	}

	if typ, ok := dic[name]; ok {
		return typ
	}

	typ := &funcType{
		name:  name,
		vars:  &namedTypeList{},
		funcs: &funcTypeList{},
	}
	dic[name] = typ
	ftl = &dic
	return typ
}

func newCommentMap(src *ast.File, fset *token.FileSet) *commentMap {
	// Create an ast.CommentMap from the ast.File's comments.
	// This helps keeping the association between comments
	// and AST nodes.
	cmap := ast.NewCommentMap(fset, src, src.Comments)
	cm := &commentMap{fset, map[int]string{}}

	for _, c := range cmap.Comments() {
		p := c.End()
		l := fset.File(p).Line(p)
		cm.comments[l] = strings.TrimSpace(c.Text())
	}

	return cm
}

// returns the comment on the same line
// or the line above the specific position
func (cm *commentMap) lookup(p token.Pos) string {
	if file := cm.fset.File(p); file != nil {
		l := file.Line(p)

		if d, ok := cm.comments[l]; ok {
			return d
		}

		if d, ok := cm.comments[l-1]; ok {
			return d
		}
	}

	return ""
}

func (cm *commentMap) description(p token.Pos) (string, []string) {
	return parseComment(cm.lookup(p))
}

func parseComment(comment string) (title string, desc []string) {
	for i, c := range strings.Split(comment, "\n") {
		if i == 0 {
			title = strings.TrimSpace(c)
		} else if len(c) > 0 {
			desc = append(desc, strings.TrimSpace(c))
		}
	}
	return
}
