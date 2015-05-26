package encoding

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/harboe/tonic/contrib/binding"
)

type yamlEncoder struct{}

func (e yamlEncoder) ContentType() string {
	return "text/yaml; charset=utf-8"
}

func (e yamlEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	if err, ok := v.(error); ok {
		return []byte(err.Error()), nil
	}

	return yaml.Marshal(v)
}

func (e yamlEncoder) Decode(req *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, b); err != nil {
		return err
	}

	return binding.Validate(v)
}
