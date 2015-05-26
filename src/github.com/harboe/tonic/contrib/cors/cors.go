package cors

import (
	"net/http"

	"github.com/justinas/alice"
)

type Configuration struct {
	Origin  string
	Methods string
	Headers string
}

func Public() alice.Constructor {
	return New(Configuration{
		Origin:  "*",
		Methods: "GET",
		Headers: "Accept, Accept-Encoding",
	})
}

func Default(site string) alice.Constructor {
	return New(Configuration{
		Origin:  site,
		Methods: "GET, POST, PUT, DELETE, OPTIONS",
		Headers: "Accept, Accept-Encoding, Authorization, Content-Type, Content-Length",
	})
}

func New(config Configuration) alice.Constructor {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {

			if origin := req.Header.Get("Origin"); origin == config.Origin || config.Origin == "*" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", config.Methods)
				w.Header().Set("Access-Control-Allow-Headers", config.Headers)
			}
			// Stop here if its Preflighted OPTIONS request
			if req.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}
