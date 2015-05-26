package docs

import (
	"html/template"
	"net/http"
	"path"
	"runtime"

	"github.com/justinas/alice"
)

type (
	Configuration struct {
		Template  string
		Route     string
		Directory string
	}
)

func Default() alice.Constructor {
	_, file, _, ok := runtime.Caller(1)

	if !ok {
		return nil
	}

	return New(Configuration{"", "_help", path.Dir(file)})
}

func New(config Configuration) alice.Constructor {
	if len(config.Directory) == 0 {
		if _, dir, _, ok := runtime.Caller(1); !ok {
			return nil
		} else {
			config.Directory = path.Dir(dir)
		}
	}

	pkg := NewPackage(config.Directory)
	t := template.New("doc")

	if len(config.Template) == 0 {
		t = template.Must(t.Parse(text_template))
	} else {
		t = template.Must(t.ParseFiles(config.Template))
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" && req.URL.Path == "/"+config.Route {
				if err := t.Execute(w, pkg); err != nil {
					w.Write([]byte(err.Error()))
					w.WriteHeader(500)
				} else {
					w.WriteHeader(200)
				}
				return
			}

			h.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}
