package tonic

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Context represents context for the current request. It holds request
// references, path parameters, data and registered handler for
// the route.
type Context struct {
	*http.Request
	params httprouter.Params
	api    *API
	store  map[string]interface{}
}

func (ctx *Context) Param(name string) string {
	return ctx.params.ByName(name)
}

func (ctx *Context) ParamUint(name string) uint64 {
	i, _ := strconv.ParseUint(ctx.Param(name), 0, 10)
	return i
}

func (ctx *Context) ParamInt(name string) int64 {
	i, _ := strconv.ParseInt(ctx.Param(name), 0, 10)
	return i
}

func (ctx *Context) ParamFloat(name string) float64 {
	f, _ := strconv.ParseFloat(ctx.Param(name), 10)
	return f
}

// func (ctx *Context) HasAcceptType(mimeType string) bool {
// 	for _, c := range ctx.Request.Header["Accept"] {
// 		if mimeType == c {
// 			return true
// 		}
// 	}
//
// 	return false
// }

//shorthand query params.
func (ctx *Context) Query(name string) string {
	if q := ctx.QueryAll(name); len(q) > 0 {
		return q[0]
	}
	return ""
}

func (ctx *Context) QueryUint(name string) uint64 {
	i, _ := strconv.ParseUint(ctx.Query(name), 0, 10)
	return i
}

func (ctx *Context) QueryInt(name string) int64 {
	i, _ := strconv.ParseInt(ctx.Query(name), 0, 10)
	return i
}

func (ctx *Context) QueryFloat(name string) float64 {
	f, _ := strconv.ParseFloat(ctx.Query(name), 10)
	return f
}

func (ctx *Context) QueryBool(name string) bool {
	b, _ := strconv.ParseBool(ctx.Query(name))
	return b
}

func (ctx *Context) QueryAll(name string) []string {
	if q, ok := ctx.Request.URL.Query()[name]; ok {
		return q
	}
	return []string{}
}

func (ctx *Context) HasQuery(name string) bool {
	_, ok := ctx.Request.URL.Query()[name]
	return ok
}

// returns 200 Ok with data
func (ctx *Context) Ok(data interface{}) *Response {
	return &Response{Context: ctx, Data: data, Status: http.StatusOK}
}

// returns 201 Created with location to the new resource
func (ctx *Context) Created(location string) *Response {
	return &Response{Context: ctx, Status: http.StatusCreated, Headers: http.Header{"Location": []string{location}}}
}

// returns 202 Accepted with location to the new resource
func (ctx *Context) Accepted(location string) *Response {
	return &Response{Context: ctx, Status: http.StatusAccepted, Headers: http.Header{"Location": []string{location}}}
}

// returns 204 NoContent
func (ctx *Context) NoContent() *Response {
	return &Response{Context: ctx, Status: http.StatusNoContent}
}

// returns 501 NotImplemented
func (ctx *Context) NotImplemented() *Response {
	return &Response{Context: ctx, Status: http.StatusNotImplemented}
}

// depending on the error it return
// 404 NotFound - if the error message is "not found"
// 400 BadRequest - if the error message starts with "valiation: "
// 500 InternalServerError - it fallsback
func (ctx *Context) Fail(err error) *Response {
	msg := err.Error()

	if msg == "not found" {
		return &Response{Context: ctx, Error: nil, Status: http.StatusNotFound}
	} else if strings.Index(msg, "validation: ") == 0 {
		return &Response{Context: ctx, Error: err, Status: http.StatusBadRequest}
	}

	return &Response{Context: ctx, Error: err, Status: http.StatusInternalServerError}
}

func (ctx *Context) Error(status int, err error) *Response {
	return &Response{Context: ctx, Error: err, Status: status}
}

func (ctx *Context) Errorf(status int, format string, args ...interface{}) *Response {
	return ctx.Error(status, errors.New(fmt.Sprintf(format, args...)))
}

func (ctx *Context) Bind(v interface{}) error {
	return ctx.api.Encoders.Decode(ctx.Request, v)
}

func (ctx *Context) Get(key string) interface{} {
	return ctx.store[key]
}

func (ctx *Context) Set(key string, val interface{}) {
	ctx.store[key] = val
}
