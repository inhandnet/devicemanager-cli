package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, string(e.Body))
}

type APIClient struct {
	inner *resty.Client
}

func NewAPIClient(baseURL string, transport http.RoundTripper) *APIClient {
	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetTransport(transport)
	return &APIClient{inner: c}
}

func (c *APIClient) Get(path string, query url.Values) ([]byte, error) {
	r := c.inner.R()
	if clean := cleanValues(query); len(clean) > 0 {
		r.SetQueryParamsFromValues(clean)
	}
	return c.execute(r, resty.MethodGet, path)
}

func (c *APIClient) Post(path string, body interface{}) ([]byte, error) {
	r := c.inner.R()
	if body != nil {
		r.SetBody(body)
	}
	return c.execute(r, resty.MethodPost, path)
}

func (c *APIClient) Put(path string, body interface{}) ([]byte, error) {
	r := c.inner.R()
	if body != nil {
		r.SetBody(body)
	}
	return c.execute(r, resty.MethodPut, path)
}

func (c *APIClient) Delete(path string) ([]byte, error) {
	return c.execute(c.inner.R(), resty.MethodDelete, path)
}

func (c *APIClient) Upload(path, fieldName, fileName string, reader io.Reader) ([]byte, error) {
	r := c.inner.R().SetFileReader(fieldName, fileName, reader)
	return c.execute(r, resty.MethodPost, path)
}

type RequestOptions struct {
	Query       url.Values
	Body        interface{}
	RawBody     io.Reader
	Headers     map[string]string
	ContentType string
}

func (c *APIClient) Do(method, path string, opts *RequestOptions) ([]byte, error) {
	r := c.inner.R()
	if opts != nil {
		if clean := cleanValues(opts.Query); len(clean) > 0 {
			r.SetQueryParamsFromValues(clean)
		}
		if len(opts.Headers) > 0 {
			r.SetHeaders(opts.Headers)
		}
		if opts.RawBody != nil {
			r.SetBody(opts.RawBody)
			if opts.ContentType != "" {
				r.SetHeader("Content-Type", opts.ContentType)
			}
		} else if opts.Body != nil {
			r.SetBody(opts.Body)
		}
	}
	return c.execute(r, method, path)
}

func (c *APIClient) execute(r *resty.Request, method, path string) ([]byte, error) {
	resp, err := r.Execute(method, path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	body := resp.Body()
	if resp.IsError() {
		return body, &HTTPError{StatusCode: resp.StatusCode(), Body: body}
	}
	return body, nil
}

func (c *APIClient) Download(path, destFile string) error {
	reqURL := path
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		reqURL = c.inner.BaseURL + path
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("creating download request: %w", err)
	}
	resp, err := c.inner.GetClient().Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	f, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = f.Close()
		_ = os.Remove(destFile)
		return fmt.Errorf("writing file: %w", err)
	}
	return f.Close()
}

func (c *APIClient) HTTPClient() *http.Client {
	return c.inner.GetClient()
}

func (c *APIClient) BaseURL() string {
	return c.inner.BaseURL
}

func cleanValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}
	clean := make(url.Values)
	for k, vals := range v {
		for _, val := range vals {
			if val == "" {
				continue
			}
			clean.Add(k, val)
		}
	}
	return clean
}
