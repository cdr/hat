package hat

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

// RequestOption modifies a request.
// Use the passed t to fail if the option cannot be set.
type RequestOption func(t testing.TB, req *http.Request)

// URLParams sets the URL parameters of the request.
func URLParams(v url.Values) RequestOption {
	return func(_ testing.TB, req *http.Request) {
		req.URL.RawQuery += v.Encode()
	}
}

// Path joins elem on to the URL.
func Path(elem string) RequestOption {
	return func(_ testing.TB, req *http.Request) {
		req.URL.Path = path.Join(req.URL.Path, elem)
		// preserve trailing slash
		if elem[len(elem)-1] == '/' && req.URL.Path != "/" {
			req.URL.Path += "/"
		}
	}
}

// Body sets the body of a request.
func Body(r io.Reader) RequestOption {
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}

	return func(_ testing.TB, req *http.Request) {
		req.Body = rc
	}
}

// CombineRequestOptions returns a new RequestOption which internally
// calls each member of options in the provided order.
func CombineRequestOptions(opts ...RequestOption) RequestOption {
	return func(t testing.TB, req *http.Request) {
		t.Helper()
		for _, o := range opts {
			o(t, req)
		}
	}
}

// Request represents a pending HTTP request.
type Request struct {
	r *http.Request

	// copy creates an exact copy of the request.
	copy func() *http.Request
}

func makeRequest(t testing.TB, copy func() *http.Request) Request {
	t.Helper()
	req := Request{
		r:    copy(),
		copy: copy,
	}
	return req
}

// Send dispatches the HTTP request.
func (r Request) Send(t *T) *Response {
	t.Helper()
	t.Logf("%v %v", r.r.Method, r.r.URL)

	resp, err := t.Client.Do(r.r)
	require.NoError(t, err, "failed to send request")

	return &Response{
		Response: resp,
	}
}

// Clone creates a duplicate HTTP request and applies opts to it.
func (r Request) Clone(t testing.TB, opts ...RequestOption) Request {
	t.Helper()
	return makeRequest(t, func() *http.Request {
		t.Helper()
		req := r.copy()
		for _, opt := range opts {
			opt(t, req)
		}
		return req
	})
}

// Request creates an HTTP request to the endpoint.
func (t T) Request(method string, opts ...RequestOption) Request {
	return makeRequest(t.T,
		func() *http.Request {
			req, err := http.NewRequest(method, t.URL.String(), nil)
			require.NoError(t, err, "failed to create request")

			for _, opt := range opts {
				opt(t, req)
			}

			return req
		},
	)
}

func (t *T) Get(opts ...RequestOption) Request {
	return t.Request(http.MethodGet, opts...)
}

func (t *T) Head(opts ...RequestOption) Request {
	return t.Request(http.MethodHead, opts...)
}

func (t *T) Post(opts ...RequestOption) Request {
	return t.Request(http.MethodPost, opts...)
}

func (t *T) Put(opts ...RequestOption) Request {
	return t.Request(http.MethodPut, opts...)
}

func (t *T) Patch(opts ...RequestOption) Request {
	return t.Request(http.MethodPatch, opts...)
}

func (t *T) Delete(opts ...RequestOption) Request {
	return t.Request(http.MethodDelete, opts...)
}
