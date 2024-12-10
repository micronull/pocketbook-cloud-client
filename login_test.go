package pocketbook_cloud_client_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pbc "github.com/micronull/pocketbook-cloud-client"
	"github.com/micronull/pocketbook-cloud-client/mocks"
)

func TestClient_Login(t *testing.T) {
	t.Parallel()

	const (
		userName     = "some.user.name"
		password     = "some.password"
		clientID     = "some.client.id"
		clientSecret = "some.client.secret"
		storeID      = "some.store.id"
		provideAlias = "some_provider"
	)

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)

	client := pbc.New(
		pbc.WithHTTPClient(httpMock),
		pbc.WithClientID(clientID),
		pbc.WithClientSecret(clientSecret),
	)

	body := must(testdata.Open("testdata/token.json"))

	httpMock.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			return isAllTrue(
				assert.Equal(t, http.MethodPost, req.Method),
				assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type")),
				assert.Equal(t, storeID, req.FormValue("store_id")),
				assert.Equal(t, userName, req.FormValue("username")),
				assert.Equal(t, password, req.FormValue("password")),
				assert.Equal(t, clientID, req.FormValue("client_id")),
				assert.Equal(t, clientSecret, req.FormValue("client_secret")),
				assert.Equal(t, "password", req.FormValue("grant_type")),
				assert.Equal(t, "https://cloud.pocketbook.digital/api/v1.0/auth/login/some_provider", req.URL.String()),
			)
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil)

	req := pbc.LoginRequest{
		ShopID:   storeID,
		UserName: userName,
		Password: password,
		Provider: provideAlias,
	}

	got, err := client.Login(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "some.access.token", got.AccessToken)
	assert.Equal(t, pbc.TokenTypeBearer, got.TokenType)
	assert.Equal(t, "some.refresh.token", got.RefreshToken)
	assert.WithinDuration(t, time.Now().Add(time.Second*7200), got.ExpiresIn, time.Second)
}

func TestClient_Login_HTTPCode_NoOk(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)

	client := pbc.New(pbc.WithHTTPClient(httpMock))

	code := rand.N(399) + 200 // rand http code > 200 and < 600

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{StatusCode: code}, nil)

	_, err := client.Login(context.Background(), pbc.LoginRequest{})
	require.ErrorContains(t, err, "http status code: "+strconv.Itoa(code)+" "+http.StatusText(code))

	var codeError interface{ Code() int }
	require.ErrorAs(t, err, &codeError)

	assert.Equal(t, code, codeError.Code())
}

func TestClient_Login_Error(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)

	client := pbc.New(pbc.WithHTTPClient(httpMock))

	errExpected := errors.New("something went wrong")

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(nil, errExpected)

	_, err := client.Login(context.Background(), pbc.LoginRequest{})
	require.ErrorIs(t, err, errExpected)
}
