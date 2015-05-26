package docs

import "testing"

var queryTests = map[string]Parameter{
	"ctx.Query(\"test\")": Parameter{
		Name:     "test",
		Type:     "string",
		Optional: true,
		Multiple: false,
	},
	"ctx.QueryInt(\"test\")": Parameter{
		Name:     "test",
		Type:     "int",
		Optional: true,
		Multiple: false,
	},
	"ctx.QueryUint(\"test\")": Parameter{
		Name:     "test",
		Type:     "uint",
		Optional: true,
		Multiple: false,
	},
	"ctx.QueryFloat(\"test\")": Parameter{
		Name:     "test",
		Type:     "float",
		Optional: true,
		Multiple: false,
	},
	"ctx.QueryBool(\"test\")": Parameter{
		Name:     "test",
		Type:     "bool",
		Optional: true,
		Multiple: false,
	},
	"ctx.QueryAll(\"test\")": Parameter{
		Name:     "test",
		Type:     "string",
		Optional: true,
		Multiple: true,
	},
	"ctx.HasQuery(\"test\")": Parameter{
		Name:     "test",
		Type:     "bool",
		Optional: true,
	},
	"ctx.Param(\"test\")": Parameter{
		Name: "test",
		Type: "string",
	},
	"ctx.ParamInt(\"test\")": Parameter{
		Name: "test",
		Type: "int",
	},
	"ctx.ParamUint(\"test\")": Parameter{
		Name: "test",
		Type: "uint",
	},
	"ctx.ParamFloat(\"test\")": Parameter{
		Name: "test",
		Type: "float",
	},
}

func TestParseQuery(t *testing.T) {
	for k, v := range queryTests {
		actual, _ := parseQuery(k)

		if actual != v {
			t.Errorf("expected %#v got %#v", v, actual)
			return
		}
	}
}
