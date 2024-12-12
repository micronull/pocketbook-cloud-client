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

func TestClient_Books(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)
	client := pbc.New(pbc.WithHTTPClient(httpMock))

	body := must(testdata.Open("testdata/books.json"))

	httpMock.EXPECT().
		Do(mock.MatchedBy(func(req *http.Request) bool {
			return isAllTrue(
				assert.Equal(t, http.MethodGet, req.Method),
				assert.Equal(t, pbc.DefaultScheme, req.URL.Scheme),
				assert.Equal(t, pbc.DefaultHost, req.URL.Host),
				assert.Equal(t, "/api/v1.0/books", req.URL.Path),
				assert.Equal(t, "limit=222", req.URL.Query().Encode()),
				assert.Equal(t, "Bearer some.token", req.Header.Get("Authorization")),
			)
		})).
		Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

	books, err := client.Books(context.Background(), "some.token", 222, 0)
	require.NoError(t, err)

	expected := pbc.Books{
		Total: 2,
		Books: []pbc.Book{
			{
				ID:          "76220203",
				Path:        "/voina-i-mir.epub",
				Title:       "Война и мир",
				MimeType:    "application/epub+zip",
				CreatedAt:   time.Date(2024, time.December, 11, 15, 41, 28, 0, time.UTC),
				Purchased:   false,
				Bytes:       2039555,
				ClientMtime: time.Date(2024, time.December, 10, 15, 41, 28, 0, time.UTC),
				FastHash:    "5c624ec0db399a8f1b99eddabf1e22c1",
				Favorite:    false,
				ReadStatus:  "read",
				Link:        "https://cloud.pocketbook.digital/api/v1.0/files/voina-i-mir.epub?fast_hash=5c624ec0db399a8f1b99eddabf1e22c1&access_token=some.token",
				HasLinks:    true,
				Format:      "epub",
				Md5Hash:     "WW/v6YxXMXC2Zi4a5x71oA==",
				Mtime:       time.Date(2024, time.December, 11, 14, 42, 17, 0, time.UTC),
				Name:        "voina-i-mir.epub",
				ReadPercent: 100,
				Percent:     "100",
				IsDrm:       false,
				IsLcp:       false,
				IsAudioBook: false,
				MetaData: pbc.BookMetaData{
					Title:   "Война и мир",
					Authors: "Толстой Л.Н.",
					Cover: []pbc.BookCover{
						{
							Width:  300,
							Height: 291,
							Path:   "https://cloud.pocketbook.digital/api/v1.0/fileops/cover/voina-i-mir.epub.cover_s.jpg?fast_hash=5c624ec0db399a8f1b99eddabf1e22c1&access_token=some.token",
						},
						{
							Width:  600,
							Height: 582,
							Path:   "https://cloud.pocketbook.digital/api/v1.0/fileops/cover/voina-i-mir.epub.cover_b.jpg?fast_hash=5c624ec0db399a8f1b99eddabf1e22c1&access_token=some.token",
						},
					},
					Lang:        "ru",
					Publisher:   "ДА!Медиа",
					Updated:     time.Date(2024, time.December, 11, 15, 41, 31, 0, time.UTC),
					Year:        2014,
					Isbn:        "9785447237509",
					BookId:      []string{"urn:uuid:95f0109f-dacc-48fc-9bb5-9456c37803b7"},
					FixedLayout: false,
				},
				Position: pbc.BookPosition{
					Pointer:    "some_position_pointer",
					PointerPb:  "some_position_pointer_pb",
					Percent:    10,
					Page:       1,
					PagesTotal: 999,
					Updated:    time.Time{},
					Offs:       0,
				},
				ReadPosition: pbc.BookReadPosition{
					Pointer:    "",
					PointerPb:  "",
					Percent:    0,
					Page:       0,
					PagesTotal: 0,
					Updated:    time.Time{},
					Offs:       0,
				},
				Action:     "create",
				ActionDate: time.Date(2024, time.November, 11, 15, 41, 28, 0, time.UTC),
			},
			{
				ID:          "76220340",
				Path:        "/puteshestvie-iz-peterburga-v-moskvu.epub",
				Title:       "Путешествие из Петербурга в Москву",
				MimeType:    "application/epub+zip",
				CreatedAt:   time.Date(2024, time.December, 11, 15, 44, 46, 0, time.UTC),
				Purchased:   false,
				Bytes:       292816,
				ClientMtime: time.Date(2024, time.December, 11, 15, 44, 45, 0, time.UTC),
				FastHash:    "01882d1bb27a52caba5d8459c80db321",
				Favorite:    true,
				ReadStatus:  "reading",
				Link:        "https://cloud.pocketbook.digital/api/v1.0/files/puteshestvie-iz-peterburga-v-moskvu.epub?fast_hash=01882d1bb27a52caba5d8459c80db321&access_token=some.token",
				HasLinks:    true,
				Format:      "epub",
				Md5Hash:     "6gDHcYaOMWA9qoovZeSUZw==",
				Mtime:       time.Date(2024, time.December, 11, 15, 44, 50, 0, time.UTC),
				Name:        "puteshestvie-iz-peterburga-v-moskvu.epub",
				ReadPercent: 4,
				Percent:     "4",
				IsDrm:       false,
				IsLcp:       false,
				IsAudioBook: false,
				MetaData: pbc.BookMetaData{
					Title:   "Путешествие из Петербурга в Москву",
					Authors: "Радищев А.Н.",
					Cover: []pbc.BookCover{
						{
							Width:  256,
							Height: 400,
							Path:   "https://cloud.pocketbook.digital/api/v1.0/fileops/cover/puteshestvie-iz-peterburga-v-moskvu.epub.cover_s.jpg?fast_hash=01882d1bb27a52caba5d8459c80db321&access_token=some.token",
						},
						{
							Width:  512,
							Height: 800,
							Path:   "https://cloud.pocketbook.digital/api/v1.0/fileops/cover/puteshestvie-iz-peterburga-v-moskvu.epub.cover_b.jpg?fast_hash=01882d1bb27a52caba5d8459c80db321&access_token=some.token",
						},
					},
					Lang:        "ru",
					Publisher:   "",
					Updated:     time.Date(2024, time.December, 11, 15, 44, 50, 0, time.UTC),
					Year:        2019,
					Isbn:        "",
					BookId:      []string{"urn:uuid:8f743510-6b3e-4bbe-9d3f-447ef0788ad0"},
					FixedLayout: false,
				},
				Position: pbc.BookPosition{
					Pointer:    "",
					PointerPb:  "",
					Percent:    0,
					Page:       0,
					PagesTotal: 0,
					Updated:    time.Time{},
					Offs:       0,
				},
				ReadPosition: pbc.BookReadPosition{
					Pointer:    "",
					PointerPb:  "",
					Percent:    0,
					Page:       0,
					PagesTotal: 0,
					Updated:    time.Time{},
					Offs:       0,
				},
				Action:     "create",
				ActionDate: time.Date(2024, time.December, 11, 15, 44, 46, 0, time.UTC),
			},
		},
	}

	assert.Equal(t, expected, books)
}

func TestClient_Books_Error_StatusCode_NoOk(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)
	client := pbc.New(pbc.WithHTTPClient(httpMock))
	code := rand.N(399) + 200 // rand http code > 200 and < 600

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{StatusCode: code}, nil)

	_, err := client.Books(context.Background(), "", 0, 0)
	require.ErrorContains(t, err, "http status code: "+strconv.Itoa(code)+" "+http.StatusText(code))

	var codeError interface{ Code() int }
	require.ErrorAs(t, err, &codeError)

	assert.Equal(t, code, codeError.Code())
}

func TestClient_Books_Error(t *testing.T) {
	t.Parallel()

	ctrlMock := gomock.NewController(t)
	httpMock := mocks.NewMockDoer(ctrlMock)
	client := pbc.New(pbc.WithHTTPClient(httpMock))
	errExpected := errors.New("something went wrong")

	httpMock.EXPECT().
		Do(gomock.Any()).
		Return(nil, errExpected)

	_, err := client.Books(context.Background(), "", 0, 0)
	require.ErrorIs(t, err, errExpected)
}
