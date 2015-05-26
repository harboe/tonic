package timeout

import (
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

type Configuration struct {
	HandlerTimeout time.Duration
	Message        string
}

func Default(duration string) alice.Constructor {
	timeout, err := time.ParseDuration(duration)

	if err != nil {
		log.Panic(err)
	}

	return New(Configuration{timeout, "request processing timed out"})
}

func New(config Configuration) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.TimeoutHandler(h, config.HandlerTimeout, config.Message)
	}
}
