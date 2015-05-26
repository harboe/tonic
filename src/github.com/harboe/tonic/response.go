package tonic

import "net/http"

type Response struct {
	Error   error
	Data    interface{}
	Status  int
	Headers http.Header
	Context *Context
}

func (r *Response) write(w http.ResponseWriter) {
	if r.Headers != nil {
		for key, headers := range r.Headers {
			for _, value := range headers {
				w.Header().Add(key, value)
			}
		}
	}

	var data interface{}

	if r.Error != nil {
		data = r.Error
	} else {
		data = r.Data
	}

	switch data.(type) {
	case []byte:
		// w.Header().Add("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(data.([]byte))
	default:
		ct, e := r.Context.api.Encoders.Encode(r.Context.Request)
		w.Header().Add("Content-Type", ct)

		if b, err := e.Encode(data); err == nil {
			w.WriteHeader(r.Status)
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
