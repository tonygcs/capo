package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/tonygcs/capo/marshaler"
)

var (
	client httpClient
	m      marshaler.Marshaler
)

func init() {
	// Set default http client and marshaler.
	client = &http.Client{}
	m = marshaler.GetMarshaler()
}

// SetHTTPClient sets the http client that will request the server.
func SetHTTPClient(c httpClient) {
	client = c
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request is a http request.
type Request struct {
	mu sync.Mutex

	headers      map[string]string
	method       string
	url          string
	relativePath string
	data         interface{}
}

// NewRequest returns a new http request instance.
func NewRequest() *Request {
	return &Request{
		headers:      make(map[string]string),
		method:       http.MethodGet,
		url:          "/",
		relativePath: "",
	}
}

// AddHeader include a new header in the request.
func (r *Request) AddHeader(key string, value string) *Request {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.headers[key] = value
	return r
}

// URL sets the url the request will be sent.
func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

// RelativePath sets the relative path for the url where the request will be
// sent.
func (r *Request) RelativePath(url string) *Request {
	r.relativePath = url
	return r
}

// Method sets the request http method.
func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// Data sets the request data that will be sent.
func (r *Request) Data(data interface{}) *Request {
	r.data = data
	return r
}

// Do performs the http request.
func (r *Request) Do(response interface{}) error {
	// Create request body if it is needed.
	var body io.Reader
	if r.data != nil {
		bodyContent, err := m.Marshal(r.data)
		if err != nil {
			return fmt.Errorf("cannot marshal the request body :: %w", err)
		}

		body = bytes.NewReader(bodyContent)
	}

	// Create request instance.
	url, err := joinPaths(r.url, r.relativePath)
	if err != nil {
		return fmt.Errorf("invalid url format :: %w", err)
	}

	req, err := http.NewRequest(r.method, url, body)
	if err != nil {
		return fmt.Errorf("cannot create the request entity :: %w", err)
	}

	// Set headers.
	for key, value := range r.headers {
		req.Header.Add(key, value)
	}

	// Send request.
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot perform the http request :: %w", err)
	}

	// Handle error from server side.
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("cannot read the error response body :: %w", err)
		}

		return newServerError(res.StatusCode, data)
	}

	// Read response if it is needed.
	if response != nil {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("cannot read the response body :: %w", err)
		}

		err = m.Unmarshal(data, response)
		if err != nil {
			return fmt.Errorf("invalid response format :: %w", err)
		}
	}

	return nil
}

func joinPaths(path1 string, path2 string) (string, error) {
	u, err := url.Parse(path1)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, path2)
	return u.String(), nil
}
