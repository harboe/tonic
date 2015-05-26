package docs

type (
	Package struct {
		Name        string
		Description []string
		Routes      []*Route
	}
	Route struct {
		Name        string
		Description string
		Path        string
		Method      string
		Body        Input
		Parameters  []Parameter
		Responses   []Response
	}
	Group struct {
		Name        string
		Description []string
		Path        string
		Routes      []*Route
	}
	Input struct {
		Name        string
		Description string
		Type        string
		Tag         string
		Children    []Input
	}
	Parameter struct {
		Name        string
		Description string
		Type        string
		Optional    bool
		Multiple    bool
	}
	Response struct {
		Name        string
		Description string
		StatusCode  int
		Body        Input
		Headers     map[string]string
	}
)

func NewPackage(dir string) *Package {
	// r := &reader{
	// 	fset: token.NewFileSet(),
	// }
	// pkgs, err := parser.ParseDir(r.fset, dir, nil, parser.ParseComments)
	//
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

	return nil
}
