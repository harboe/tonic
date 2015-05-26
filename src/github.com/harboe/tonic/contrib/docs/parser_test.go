package docs

import (
	"log"
	"testing"

	"github.com/harboe/tonic"
)

func TestParserFunc(t *testing.T) {
	fn := ParseFunc(func(ctx *tonic.Context) *tonic.Response {
		ctx.Param("Aichata") // hej skat
		return nil
	})

	if e := len(fn.Params); e != 1 {
		t.Fatalf("expected %v got %v", 1, e)
	}

	if n := fn.Params[0].Name; n != "Aichata" {
		t.Fatalf("expected %v got %v", "Aichata", n)
	}

	if d := fn.Params[0].Desc; d != "Hej Skat" {
		t.Fatalf("expected %v got %v", "Hej Skat", d)
	}

	log.Println("func:", fn)
}
