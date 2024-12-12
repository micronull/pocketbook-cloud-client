//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -mock_names doer=Doer
package pocketbook_cloud_client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	DefaultScheme = "https"
	DefaultHost   = "cloud.pocketbook.digital"
	DefaultPath   = "/api/v1.0/"

	login = "auth/login"
	books = "books"
)

type Client struct {
	http         doer
	scheme       string
	host         string
	path         string
	clientID     string
	clientSecret string
}

func New(opts ...Option) *Client {
	c := &Client{
		scheme: DefaultScheme,
		host:   DefaultHost,
		path:   DefaultPath,
		http:   http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c Client) req(req *http.Request) ([]byte, error) {
	rsp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, httpStatusError{rsp.StatusCode}
	}

	defer func() { _ = rsp.Body.Close() }()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return body, nil
}

func (c Client) url(endpoint string) *url.URL {
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.path,
	}

	return u.JoinPath(endpoint)
}
