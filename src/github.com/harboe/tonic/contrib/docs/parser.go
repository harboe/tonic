package docs

import (
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type (
	// Parser looks though sourcecode to
	// find all documentation.
	reader struct {
		fset *token.FileSet
	}
	FileDesc struct {
		Funcs   []FuncDesc
		Imports map[string]string
	}
	// FuncDesc represents the result of func :-> lol
	FuncDesc struct {
		File string
		Line int
		// route parameters found
		Params []ParamDesc
		// query paramters found.
		Queries []ParamDesc
	}
	ParamDesc struct {
		// Name of the param used
		Name string
		// Comment assosiated with the Param
		Desc string
		// the primitive type it (string,int,float,bool,date?)
		Type string
	}
	ReturnDesc struct {
	}
)

func newFileDesc(file string) (FileDesc, error) {
	return FileDesc{}, nil
}

func newReader() *reader {
	return &reader{
		fset: token.NewFileSet(),
	}
}

func ParseFunc(fn interface{}) FuncDesc {
	r := newReader()
	return r.Func(fn)
}

// Func parse func body and return
func (r *reader) Func(fn interface{}) FuncDesc {
	pc := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	file, line := pc.FileLine(pc.Entry())

	fd := r.readFile(file)

	for _, d := range fd.Funcs {
		if d.Line == line {
			return d
		}
	}

	return FuncDesc{}
}

func (r *reader) readFile(file string) FileDesc {
	fd := FileDesc{
		Funcs:   []FuncDesc{},
		Imports: map[string]string{},
	}

	src, err := parser.ParseFile(r.fset, file, nil, parser.ParseComments)

	if err != nil {
		log.Println(err)
		return fd
	}

	// init file imports
	for _, i := range src.Imports {
		if path, err := strconv.Unquote(i.Path.Value); err == nil {
			if i.Name != nil {
				fd.Imports[path] = i.Name.Name
			} else if index := strings.LastIndex(path, "/"); index > 0 {
				fd.Imports[path] = path[index+1:]
			} else {
				fd.Imports[path] = path
			}
		}
	}

	log.Println("imports:", fd.Imports)

	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil
	// }
	//
	// for _, pkg := range pkgs {
	// 	r.readPackage(pkg)
	//
	// 	name := pkg.Name
	//
	// 	if len(r.name) > 0 {
	// 		name = r.name
	// 	}
	//
	// 	return &Package{
	// 		Name:        name,
	// 		Description: r.desc,
	// 		// Routes:      r.sortedRoutes(r.routes),
	// 	}
	// }

	return fd
}
