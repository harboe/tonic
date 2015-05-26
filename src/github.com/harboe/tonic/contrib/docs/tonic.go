package docs

import (
	"regexp"
	"strings"
)

var (
	// 1: Has?
	// 2: Param|Query
	// 3: Type int|float|uint|bool|all
	// 4: Name
	queryExp = regexp.MustCompile(`.+\.(Has)?(Query|Param)([^\(]+)?\("(\w+)"\)`)
	bindExp  = regexp.MustCompile(`.+\.Bind\("(\w+)"\)`)
	routeExp = regexp.MustCompile(`.+\.(Get|Post|Put|Delete|Patch|Head|Route)`)
	// returnExp = regexp.MustCompile(`.+\.()\(`)
	// initExp  = regexp.MustCompile(`tonic\.(New|Default)\(`)
)

func parseQuery(method string) (Parameter, bool) {
	if m := queryExp.FindAllStringSubmatch(method, -1); len(m) > 0 {
		t := strings.ToLower(m[0][3])
		multiple := false
		optional := true

		if t == "all" {
			multiple = true
			t = "string"
		} else if len(m[0][1]) > 0 { // hasQuery
			t = "bool"
		} else if len(t) == 0 {
			t = "string"
		}

		if m[0][2] == "Param" {
			optional = false
		}

		return Parameter{
			Name:     m[0][4],
			Type:     t,
			Optional: optional,
			Multiple: multiple,
		}, true
	}

	return Parameter{}, false
}

func isConstructor(method string) bool {
	return method == "New" || method == "Default"
}

func isParam(method string) bool {
	return strings.HasPrefix(method, "Param")
}

func isQuery(method string) bool {
	return strings.HasPrefix(method, "Query") || method == "HasQuery"
}

func isBind(method string) bool {
	return method == "Bind"
}

func isResponse(method string) bool {
	responses := []string{"Ok", "Created", "Accepted", "NoContent", "NotImplemented", "Fail", "Error", "Response"}
	for _, r := range responses {
		if r == method {
			return true
		}
	}

	return false
}

func isGroup(method string) bool {
	return method == "Group"
}

func isRoute(method string) bool {
	rotues := []string{"Get", "Post", "Put", "Delete", "Patch", "Route"}
	for _, r := range rotues {
		if r == method {
			return true
		}
	}

	return false
}

func isUse(method string) bool {
	return method == "Use"
}

func parseBind(method string) (string, bool) {
	if m := bindExp.FindStringSubmatch(method); len(m) > 0 {
		return m[1], true
	}
	return "", false
}

func parseRoute(method string) (string, bool) {
	if m := routeExp.FindStringSubmatch(method); len(m) > 0 {
		return m[1], true
	}
	return "", false
}

func parseReturn(method string) ([]int, bool) {
	method = method[strings.Index(method, ".")+1:]

	switch {
	case strings.HasPrefix(method, "Ok"):
		return []int{200}, true
	case strings.HasPrefix(method, "Created"):
		return []int{201}, true
	case strings.HasPrefix(method, "Accepted"):
		return []int{202}, true
	case strings.HasPrefix(method, "NoContent"):
		return []int{204}, true
	case strings.HasPrefix(method, "NotImplemented"):
		return []int{501}, true
	case strings.HasPrefix(method, "Fail"):
		return []int{400, 404, 500}, true
	case strings.HasPrefix(method, "Error"):
		return []int{}, true
	}

	return []int{}, false
}
