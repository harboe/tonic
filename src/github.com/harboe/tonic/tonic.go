package tonic

import (
	"errors"
	"log"
	"net/http"

	"github.com/harboe/tonic/contrib/encoding"
	"github.com/harboe/tonic/contrib/logging"
	"github.com/harboe/tonic/contrib/timeout"

	router "github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	MIME_JSON = "application/json"
	MIME_XML  = "application/xml"
	MINE_YML  = "text/yaml"
	MINI_TXT  = "text/plain"
)

type (
	API struct {
		router   *router.Router
		chain    []alice.Constructor
		Encoders *encoding.Encoders
		*RouteGroup
	}
	HandleFunc func(*Context) *Response
)

// Returns a new blank API instance without any middleware.
// Default encoder is JSON. The most basic configuration.
func New() *API {
	api := &API{
		router:   router.New(),
		chain:    []alice.Constructor{},
		Encoders: encoding.New(),
	}
	api.RouteGroup = &RouteGroup{[]HandleFunc{}, "/", nil, api}
	return api
}

// Returns a new API instance with Logging and Timeout middleware
// The default encoders are JSON, XML, YAML & FORM.
func Default() *API {
	api := New()
	api.Use(
		timeout.Default("12s"),
		logging.Default("TONIC"))

	// adding default encoders..
	api.Encoder(
		encoding.JSON,
		encoding.XML,
		encoding.YAML,
		encoding.FORM)

	return api
}

func (api *API) Encoder(encoders ...encoding.Encoding) {
	api.Encoders.Add(encoders...)
}

func (api *API) Use(constructors ...alice.Constructor) {
	api.chain = append(api.chain, constructors...)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (api *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	chain := alice.New(api.chain...).Then(api.router)
	chain.ServeHTTP(w, req)
}

func (api *API) Run(addr string) {
	log.Println("Tonic server ready at:", addr)
	log.Fatal(http.ListenAndServe(addr, api))
}

func (api *API) handleRoute(route, method string, handlers []HandleFunc) {
	if len(handlers) == 0 {
		log.Panic(errors.New("Missing handler(s) for route [" + method + "] " + route))
	}

	// log.Println(api.Docs.Parse(route, method, handlers...))

	api.router.Handle(method, route, func(w http.ResponseWriter, req *http.Request, ps router.Params) {
		ctx := &Context{req, ps, api, make(map[string]interface{}, 0)}

		var resp *Response

		// invoke route handlers
		for _, h := range handlers {
			resp = h(ctx)

			if resp != nil && resp.Error != nil {
				break
			}
		}

		if resp == nil {
			resp = ctx.NoContent()
		}

		resp.write(w)
	})
}
