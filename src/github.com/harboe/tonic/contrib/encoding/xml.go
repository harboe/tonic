package encoding

import (
	"encoding/xml"
	"net/http"

	"github.com/harboe/tonic/contrib/binding"
)

type xmlEncoder struct{}

func (e xmlEncoder) ContentType() string {
	return "application/xml; charset=utf-8"
}

func (e xmlEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	if err, ok := v.(error); ok {
		return []byte("<error>" + err.Error() + "</error>"), nil
	}

	return xml.Marshal(v)
}

func (e xmlEncoder) Decode(req *http.Request, v interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return binding.Validate(v)
}
