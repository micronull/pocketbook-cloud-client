package pocketbook_cloud_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Provider struct {
	Alias    string
	Name     string
	ShopID   string
	Icon     string
	IconEink string
	LoggedBy string
}

// Providers getting allowed providers list.
func (c Client) Providers(ctx context.Context, userName string) ([]Provider, error) {
	u := c.url(login)

	q := u.Query()
	q.Set("client_id", c.clientID)
	q.Set("client_secret", c.clientSecret)
	q.Set("username", userName)

	u.RawQuery = q.Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Body:   http.NoBody,
	}

	req = req.WithContext(ctx)

	body, err := c.req(req)
	if err != nil {
		return nil, fmt.Errorf("%s userName=%s: %w", login, userName, err)
	}

	var data struct {
		Providers []struct {
			Alias    string `json:"alias"`
			Name     string `json:"name"`
			ShopID   string `json:"shop_id"`
			Icon     string `json:"icon"`
			IconEink string `json:"icon_eink"`
			LoggedBy string `json:"logged_by"`
		} `json:"providers"`
	}

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
