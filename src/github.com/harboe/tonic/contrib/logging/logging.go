package logging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
)

type Configuration struct {
	AppName       string
	Logger        *log.Logger
	IgnoreFavicon bool // dont log favicon request
}

func Default(appName string) alice.Constructor {
	return New(Configuration{
		appName,
		log.New(os.Stderr, "", 0),
		true,
	})
}

func New(config Configuration) alice.Constructor {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			// ignore favicons.
			if config.IgnoreFavicon && IsFavicon(req) {
				next.ServeHTTP(w, req)
				return
			}

			start := time.Now()
			rw := &responseWriter{http: w}
			next.ServeHTTP(rw, req)

			end := time.Now()
			latency := end.Sub(start)

			str := fmt.Sprintf("%v [%s] |%s %3d %s| %12v | %s | %-7s | %s",
				end.Format("2006/01/02 15:04:05"),
				config.AppName,
				ColorForStatus(rw.status), rw.status, Reset,
				latency,
				ClientIP(req),
				req.Method,
				req.URL.Path,
			)
			config.Logger.Println(str)
		}

		return http.HandlerFunc(fn)
	}
}

func ClientIP(req *http.Request) string {
	// save the IP of the requester
	ip := req.Header.Get("X-Real-IP")

	// if the requester-header is empty,
	// check the forwarded-header
	if len(ip) == 0 {
		ip = req.Header.Get("X-Forwarded-For")
	}

	// if the requester is still empty,
	// use the hard-coded address from the socket
	if len(ip) == 0 {
		ip = req.RemoteAddr
	}

	return ip
}
