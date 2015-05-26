package logging

import "net/http"

type responseWriter struct {
	http   http.ResponseWriter
	status int
}

func (rw *responseWriter) Header() http.Header {
	return rw.http.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.http.Write(b)
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.http.WriteHeader(status)
}
