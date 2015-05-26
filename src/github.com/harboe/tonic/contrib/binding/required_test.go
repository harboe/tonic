package binding

import "testing"

type bar struct {
	A string `binding:"required"`
}

func Test_Required_String(t *testing.T) {
	type test struct {
		A string `binding:"required"`
	}

	if err := Validate(test{"muuha"}); err != nil {
		t.Error(err)
	}

	if err := Validate(test{""}); err == nil {
		t.Error(err)
	}
}

func Test_Required_Float(t *testing.T) {
	type test struct {
		A float32 `binding:"required"`
	}
	var err error

	err = Validate(test{34.1})

	if err != nil {
		t.Error(err)
	}

	err = Validate(test{0})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	exp := "validation: required A"

	if err.Error() != exp {
		t.Errorf("expected: '%s' got '%s'", exp, err.Error())
	}
}

// func Test_Required_Nested(t *testing.T) {
// 	type Foo struct {
// 		Bar string `binding:"required"`
// 		*Foo
// 	}
//
// 	test := Foo{"1", &Foo{"2", &Foo{"", nil}}}
// 	exp := "validation: required Foo.Foo.Bar"
//
// 	if err := Validate(test); err != nil {
// 		if exp != err.Error() {
// 			t.Errorf("expected: '%s' got '%s'", exp, err.Error())
// 		}
// 	} else {
// 		t.Fatal("expected an error, got nil")
// 	}
// }

func Test_Required_Array(t *testing.T) {
	type Bar struct {
		Bar string `binding:"required"`
	}
	type Foo struct {
		Arr []Bar
	}

	test := Foo{[]Bar{Bar{"A"}, Bar{"B"}, Bar{"C"}, Bar{""}}}
	exp := "validation: required Arr.[3].Bar"

	if err := Validate(test); err != nil {
		if err.Error() != exp {
			t.Errorf("expected: '%s' got '%s'", exp, err.Error())
		}
	} else {
		t.Fail()
	}
}

func Test_Required_Map(t *testing.T) {
	type Bar struct {
		Bar string `binding:"required"`
	}
	type Foo struct {
		Dic map[string]Bar
	}

	test := Foo{map[string]Bar{
		"A": Bar{"1"},
		"B": Bar{"2"},
		"C": Bar{"3"},
		"D": Bar{""},
	}}
	exp := "validation: required Dic.[D].Bar"

	if err := Validate(test); err != nil {
		if err.Error() != exp {
			t.Errorf("expected: '%s' got '%s'", exp, err.Error())
		}
	} else {
		t.Fail()
	}

}
