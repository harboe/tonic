package extensions

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/justinas/alice"
)

var (
	JSON = FileType{".json", "application/json"}
	XML  = FileType{".xml", "application/xml"}
	YAML = FileType{".yml", "text/yaml"}
	TXT  = FileType{".txt", "text/yaml"}
)

type FileType struct {
	Extension   string
	ContentType string
}

// default extensions matching:
// .json -> application/json
// .xml  -> application/xml
// .yml  -> text/yaml
// .txt  -> text/txt
func Default() alice.Constructor {
	return New(JSON, XML, YAML, XML)
}

// New known file extensions and set the correct
// accept header content-type.
func New(files ...FileType) alice.Constructor {
	dic := map[string]string{}
	extRegex := ""

	for _, f := range files {
		dic[f.Extension] = f.ContentType
		extRegex += f.Extension[1:] + "|"
	}

	exp := fmt.Sprintf(`(\.(?:%s))\/?$`, extRegex[:len(extRegex)-1])
	regExt := regexp.MustCompile(exp)

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			ft := ""

			// Get the format extension
			if m := regExt.FindStringSubmatch(req.URL.Path); len(m) > 1 {
				// Rewrite the URL without the format extension
				l := len(req.URL.Path) - len(m[1])
				if strings.HasSuffix(req.URL.Path, "/") {
					l--
				}
				req.URL.Path = req.URL.Path[:l]
				ft = m[1]
			}

			if ct, ok := dic[ft]; ok {
				req.Header.Set("Accept", ct)
			}

			h.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}
