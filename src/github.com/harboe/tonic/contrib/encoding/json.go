package encoding

import (
	"encoding/json"
	"net/http"

	"github.com/harboe/tonic/contrib/binding"
)

type jsonEncoder struct{}

func (e jsonEncoder) ContentType() string {
	return "application/json; charset=utf-8"
}

func (e jsonEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	if err, ok := v.(error); ok {
		return []byte(err.Error()), nil
	}

	return json.Marshal(v)
}

func (e jsonEncoder) Decode(req *http.Request, v interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return binding.Validate(v)
}
