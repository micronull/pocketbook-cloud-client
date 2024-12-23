package pocketbook_cloud_client_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pbc "github.com/micronull/pocketbook-cloud-client"
	"github.com/micronull/pocketbook-cloud-client/mocks"
)

func TestClient_Providers(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)

	const (
		userName     = "some.user.name"
		clientID     = "some.client.id"
		clientSecret = "some.client.secret"
	)

	client := pbc.New(
		pbc.WithHTTPClient(httpMock),
		pbc.WithClientID(clientID),
		pbc.WithClientSecret(clientSecret),
	)

	body := must(testdata.Open("testdata/providers.json"))

	httpMock.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			return isAllTrue(
				assert.Equal(t, http.MethodGet, req.Method),
				assert.Equal(t, pbc.DefaultScheme, req.URL.Scheme),
				assert.Equal(t, pbc.DefaultHost, req.URL.Host),
				assert.Equal(t, "/api/v1.0/auth/login", req.URL.Path),
				assert.Equal(t, "client_id=some.client.id&client_secret=some.client.secret&username=some.user.name", req.URL.Query().Encode()),
			)
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil)

	got, err := client.Providers(context.Background(), userName)
	require.NoError(t, err)

	expected := []pbc.Provider{
		{
			Alias:    "pocketbook_de",
			Name:     "Pocketbook.de",
			ShopID:   "1",
			Icon:     "https://pocketbook.de/static/version1732018432/webapi_rest/_view/de_DE/images/logo.svg",
			IconEink: "https://pocketbook.de/static/version1732018432/webapi_rest/_view/de_DE/images/logo.svg",
			LoggedBy: "password",
		},
		{
			Alias:    "bookland_ru",
			Name:     "bookland new",
			ShopID:   "35",
			Icon:     "https://bookland.ru/static/version1718368600/frontend/Pocketbook/russia/ru_RU/images/logo.svg",
			IconEink: "https://bookland.ru/static/version1718368600/frontend/Pocketbook/russia/ru_RU/images/logo.svg",
			LoggedBy: "facebook|gmail",
		},
	}

	assert.Equal(t, expected, got)
}

func TestClient_Providers_Error_StatusCode_NoOk(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)
	client := pbc.New(pbc.WithHTTPClient(httpMock))
	code := rand.N(399) + 200 // rand http code > 200 and < 600

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{StatusCode: code}, nil)

	_, err := client.Providers(context.Background(), "some")
	require.ErrorContains(t, err, "http status code: "+strconv.Itoa(code)+" "+http.StatusText(code))

	var codeError interface{ Code() int }
	require.ErrorAs(t, err, &codeError)

	assert.Equal(t, code, codeError.Code())
}

func TestClient_Providers_Error(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)
	client := pbc.New(pbc.WithHTTPClient(httpMock))
	errExpected := errors.New("something went wrong")

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(nil, errExpected)

	_, err := client.Providers(context.Background(), "some")
	require.ErrorIs(t, err, errExpected)
}
