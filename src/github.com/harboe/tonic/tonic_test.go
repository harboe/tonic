package tonic

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteOK(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	r.Route("/test", method, func(c *Context) *Response {
		passed = true
		return c.Ok(nil)
	})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == false {
		t.Errorf(method + " route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}
func TestRouterGroupRouteOK(t *testing.T) {
	testRouteOK("GET", t)
	testRouteOK("POST", t)
	testRouteOK("DELETE", t)
	testRouteOK("PATCH", t)
	testRouteOK("PUT", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("HEAD", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	r.Route("/test_2", method, func(c *Context) *Response {
		passed = true
		return c.Ok(nil)
	})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == true {
		t.Errorf(method + " route handler was invoked, when it should not")
	}
	if w.Code != http.StatusNotFound {
		// If this fails, it's because httprouter needs to be updated to at least f78f58a0db
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func TestRouteNotOK(t *testing.T) {
	testRouteNotOK("POST", t)
	testRouteNotOK("DELETE", t)
	testRouteNotOK("PATCH", t)
	testRouteNotOK("PUT", t)
	testRouteNotOK("OPTIONS", t)
	testRouteNotOK("HEAD", t)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK2(method string, t *testing.T) {
	// SETUP
	passed := false
	r := New()
	var methodRoute string
	if method == "POST" {
		methodRoute = "GET"
	} else {
		methodRoute = "POST"
	}
	r.Route("/test", methodRoute, func(c *Context) *Response {
		passed = true
		return c.Ok(nil)
	})

	// RUN
	w := PerformRequest(r, method, "/test")

	// TEST
	if passed == true {
		t.Errorf(method + " route handler was invoked, when it should not")
	}
	if w.Code != http.StatusNotFound {
		// If this fails, it's because httprouter needs to be updated to at least f78f58a0db
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
}
