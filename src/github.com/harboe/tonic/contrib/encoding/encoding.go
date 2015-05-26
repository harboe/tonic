package encoding

import (
	"errors"
	"net/http"
	"strings"
)

var (
	JSON = jsonEncoder{}
	XML  = xmlEncoder{}
	YAML = yamlEncoder{}
	FORM = formEncoder{}
)

type Encoding interface {
	ContentType() string
}

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(req *http.Request, v interface{}) error
}

type Encoders map[string]Encoding

func New() *Encoders {
	return &Encoders{}
}

func (e *Encoders) Add(encoders ...Encoding) {
	dic := *e

	for _, enc := range encoders {
		key := removeContentTypeEncoding(enc.ContentType())
		dic[key] = enc
	}
}

// determine the decoder to use, it looks though the content-
// type headers. if no sutible accept content-type is found
// default back to the first added encoding.Encoder.
func (e *Encoders) Decode(req *http.Request, v interface{}) error {
	ct := req.Header.Get("content-type")
	encoder := e.getEncoder(ct)

	if d, ok := encoder.(Decoder); !ok {
		return errors.New("no decoder found for content-type: " + ct)
	} else {
		return d.Decode(req, v)
	}
}

// determine the encoder to use, it looks though the accept
// headers. if no sutible accept content-type is found
// default back to the JSON encoder.
func (e *Encoders) Encode(req *http.Request) (string, Encoder) {
	encoder := e.getEncoder(req.Header.Get("accept"))

	if e, ok := encoder.(Encoder); !ok {
		return JSON.ContentType(), JSON
	} else {
		return encoder.ContentType(), e
	}
}

func (e *Encoders) getEncoder(contentTypes ...string) Encoding {
	dic := *e

	for _, ct := range contentTypes {
		ct = removeContentTypeEncoding(ct)
		if encoder, ok := dic[ct]; ok {
			return encoder
		}
	}

	return JSON
}

func removeContentTypeEncoding(contentType string) string {
	if index := strings.IndexRune(contentType, ';'); index > 0 {
		contentType = contentType[0:index]
	}

	return contentType
}
