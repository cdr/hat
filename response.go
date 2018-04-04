package hat

import (
	"net/http"
)

// ResponseAssertion asserts a quality of the response.
type ResponseAssertion func(t T, r Response)

// CombineResponseAssertions returns a new ResponseAssertion which internally
// calls each member of asserts in the provided order.
func CombineResponseAssertions(asserts ...ResponseAssertion) ResponseAssertion {
	return func(t T, r Response) {
		for _, a := range asserts {
			a(t, r)
		}
	}
}

// Response represents an HTTP response generated by hat.Request.
type Response struct {
	*http.Response
}

// Assert runs each assertion against the response.
// It closes the response body after all of the assertions have ran.
func (r Response) Assert(t T, assertions ...ResponseAssertion) Response {
	defer r.Body.Close()
	for _, a := range assertions {
		a(t, r)
	}
	return r
}
