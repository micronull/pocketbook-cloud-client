//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -mock_names doer=Doer
package pocketbook_cloud_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	DefaultHost = "https://cloud.pocketbook.digital"
	DefaultPath = "/api/v1.0/"

	login = "auth/login"
)

type Client struct {
	http         doer
	host         string
	path         string
	clientID     string
	clientSecret string
}

func New(opts ...Option) *Client {
	c := &Client{
		host: DefaultHost,
		path: DefaultPath,
		http: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c Client) Providers(ctx context.Context, userName string) ([]Provider, error) {
	type schema struct {
		Providers []struct {
			Alias    string `json:"alias"`
			Name     string `json:"name"`
			ShopID   string `json:"shop_id"`
			Icon     string `json:"icon"`
			IconEink string `json:"icon_eink"`
			LoggedBy string `json:"logged_by"`
		} `json:"providers"`
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.urlUser(login, userName), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	rsp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, httpStatusError{rsp.StatusCode}
	}

	defer func() { _ = rsp.Body.Close() }()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var data schema

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}

	result := make([]Provider, len(data.Providers))
	for i, p := range data.Providers {
		result[i] = Provider{
			Alias:    p.Alias,
			Name:     p.Name,
			ShopID:   p.ShopID,
			Icon:     p.Icon,
			IconEink: p.IconEink,
			LoggedBy: p.LoggedBy,
		}
	}

	return result, nil
}

func (c Client) url(endpoint string) string {
	return c.host + c.path + endpoint + "?client_id=" + c.clientID + "&client_secret=" + c.clientSecret
}

func (c Client) urlUser(endpoint, userName string) string {
	return c.url(endpoint) + "&username=" + userName
}
