package binding

import "testing"

func TestMin(t *testing.T) {
	type test struct {
		A int `binding:"min:10,max:20"`
	}

	if err := Validate(test{19}); err != nil {
		t.Error(err)
	}

}
