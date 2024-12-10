package pocketbook_cloud_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LoginRequest struct {
	ShopID   string
	UserName string
	Password string
	Provider string
}

type tokenType string

const TokenTypeBearer tokenType = "Bearer"

type Token struct {
	AccessToken  string
	TokenType    tokenType
	ExpiresIn    time.Time
	RefreshToken string
}

func (c Client) Login(ctx context.Context, lreq LoginRequest) (Token, error) {
	q := url.Values{}
	q.Set("store_id", lreq.ShopID)
	q.Set("username", lreq.UserName)
	q.Set("password", lreq.Password)
	q.Set("client_id", c.clientID)
	q.Set("client_secret", c.clientSecret)
	q.Set("grant_type", "password")

	req := &http.Request{
		Method: http.MethodPost,
		URL:    c.url(login).JoinPath(lreq.Provider),
		Header: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
		Body: io.NopCloser(strings.NewReader(q.Encode())),
	}

	req = req.WithContext(ctx)

	body, err := c.req(req)
	if err != nil {
		return Token{}, fmt.Errorf("%s userName=%s provider=%s: %w", login, lreq.UserName, lreq.Provider, err)
	}

	var data struct {
		AccessToken  string    `json:"access_token"`
		TokenType    tokenType `json:"token_type"`
		ExpiresIn    int       `json:"expires_in"`
		RefreshToken string    `json:"refresh_token"`
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return Token{}, fmt.Errorf("unmarshal response body: %w", err)
	}

	t := Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		ExpiresIn:    time.Now().Add(time.Duration(data.ExpiresIn) * time.Second),
		RefreshToken: data.RefreshToken,
	}

	return t, nil
}
