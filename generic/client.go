package generic

import "github.com/tonygcs/capo/client"

// Request is a http request.
type Request[T any, U any] struct {
	r *client.Request
}

// NewRequest returns a new http request instance.
func NewRequest[T any, U any]() *Request[T, U] {
	return &Request[T, U]{
		r: client.NewRequest(),
	}
}

// AddHeader include a new header in the request.
func (r *Request[T, U]) AddHeader(key string, value string) *Request[T, U] {
	r.r.AddHeader(key, value)
	return r
}

// URL sets the url the request will be sent.
func (r *Request[T, U]) URL(url string) *Request[T, U] {
	r.r.URL(url)
	return r
}

// RelativePath sets the relative path for the url where the request will be
// sent.
func (r *Request[T, U]) RelativePath(url string) *Request[T, U] {
	r.r.RelativePath(url)
	return r
}

// Method sets the request http method.
func (r *Request[T, U]) Method(method string) *Request[T, U] {
	r.r.Method(method)
	return r
}

// Data sets the request data that will be sent.
func (r *Request[T, U]) Data(data *T) *Request[T, U] {
	r.r.Data(data)
	return r
}

// Do performs the http request.
func (r *Request[T, U]) Do() (*U, error) {
	res := new(U)
	err := r.r.Do(res)
	return res, err
}
