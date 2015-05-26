package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/harboe/tonic"
	"github.com/harboe/tonic/contrib/cors"
	"github.com/harboe/tonic/contrib/docs"
	"github.com/harboe/tonic/contrib/extensions"
)

type appContext struct{}

func main() {
	// knud
	app := tonic.Default()
	app.Use(extensions.Default())
	app.Use(cors.Public())
	app.Use(docs.Default())

	// testing
	// hey palle
	app.Get("/", Auth, Route2, func(ctx *tonic.Context) *tonic.Response {
		log.Println(ctx.QueryBool("boolean"))
		return ctx.Ok("elo, world")
	})
	app.Run("localhost:3000")
}

// Authentication
func Auth(ctx *tonic.Context) *tonic.Response {
	q1, q2 := ctx.Query("test"), ctx.Query("test2") // test param 1
	q3 := ctx.Query("test3")                        /* test param 2 */
	log.Println("test:", q1, q2, q3)

	if true {
		return ctx.Fail(errors.New("true"))
	}

	return ctx.Ok(nil)
}

// hey
func Route2(ctx *tonic.Context) *tonic.Response {
	return &tonic.Response{Context: ctx, Status: http.StatusOK}
}
