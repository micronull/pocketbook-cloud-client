package pocketbook_cloud_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Books struct {
	Total int
	Books []Book
}

type Book struct {
	ID           string
	Path         string
	Title        string
	MimeType     string
	CreatedAt    time.Time
	Purchased    bool
	Bytes        int
	ClientMtime  time.Time
	FastHash     string
	Favorite     bool
	ReadStatus   string
	Link         string
	HasLinks     bool
	Format       string
	Md5Hash      string
	Mtime        time.Time
	Name         string
	ReadPercent  int
	Percent      string
	IsDrm        bool
	IsLcp        bool
	IsAudioBook  bool
	MetaData     BookMetaData
	Position     BookPosition
	ReadPosition BookReadPosition
	Action       string
	ActionDate   time.Time
}

type BookMetaData struct {
	Title       string
	Authors     string
	Cover       []BookCover
	Lang        string
	Publisher   string
	Updated     time.Time
	Year        int
	Isbn        string
	BookId      []string
	FixedLayout bool
}

type BookCover struct {
	Width  int
	Height int
	Path   string
}

type BookPosition struct {
	Pointer    string
	PointerPb  string
	Percent    int
	Page       string
	PagesTotal int
	Updated    time.Time
	Offs       int
}

type BookReadPosition struct {
	Pointer    string
	PointerPb  string
	Percent    int
	Page       string
	PagesTotal int
	Updated    time.Time
	Offs       int
}

func (c Client) Books(ctx context.Context, token string, limit, offset int) (Books, error) {
	u := c.url(books)
	q := u.Query()

	q.Set("limit", strconv.Itoa(limit))

	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}

	u.RawQuery = q.Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Body:   http.NoBody,
		Header: http.Header{"Authorization": []string{string(TokenTypeBearer) + " " + token}},
	}

	req = req.WithContext(ctx)

	body, err := c.req(req)
	if err != nil {
		return Books{}, fmt.Errorf("get books: %w", err)
	}

	var data struct {
		Total int `json:"total"`
		Items []struct {
			ID          string    `json:"id"`
			Path        string    `json:"path"`
			Title       string    `json:"title"`
			MimeType    string    `json:"mime_type"`
			CreatedAt   time.Time `json:"created_at"`
			Purchased   bool      `json:"purchased"`
			Bytes       int       `json:"bytes"`
			ClientMtime time.Time `json:"client_mtime"`
			FastHash    string    `json:"fast_hash"`
			Favorite    bool      `json:"favorite"`
			ReadStatus  string    `json:"read_status"`
			Link        string    `json:"link"`
			HasLinks    bool      `json:"hasLinks"`
			Format      string    `json:"format"`
			Md5Hash     string    `json:"md5_hash"`
			Mtime       time.Time `json:"mtime"`
			Name        string    `json:"name"`
			ReadPercent int       `json:"read_percent"`
			Percent     string    `json:"percent"`
			IsDrm       bool      `json:"isDrm"`
			IsLcp       bool      `json:"isLcp"`
			IsAudioBook bool      `json:"isAudioBook"`
			Metadata    struct {
				Title   string `json:"title"`
				Authors string `json:"authors"`
				Cover   []struct {
					Width  int    `json:"width"`
					Height int    `json:"height"`
					Path   string `json:"path"`
				} `json:"cover"`
				Lang        string    `json:"lang"`
				Publisher   string    `json:"publisher"`
				Updated     time.Time `json:"updated"`
				Year        int       `json:"year"`
				Isbn        string    `json:"isbn"`
				BookId      []string  `json:"book_id"`
				FixedLayout bool      `json:"fixed_layout"`
			} `json:"metadata"`
			Position struct {
				Pointer    string    `json:"pointer"`
				PointerPb  string    `json:"pointer_pb"`
				Percent    int       `json:"percent"`
				Page       string    `json:"page"`
				PagesTotal int       `json:"pages_total"`
				Updated    time.Time `json:"updated"`
				Offs       int       `json:"offs"`
			} `json:"position"`
			ReadPosition struct {
				Pointer    string    `json:"pointer"`
				PointerPb  string    `json:"pointer_pb"`
				Percent    int       `json:"percent"`
				Page       string    `json:"page"`
				PagesTotal int       `json:"pages_total"`
				Updated    time.Time `json:"updated"`
				Offs       int       `json:"offs"`
			} `json:"read_position"`
			Action     string    `json:"action"`
			ActionDate time.Time `json:"action_date"`
		} `json:"items"`
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return Books{}, fmt.Errorf("unmarshal response body: %w", err)
	}

	books := Books{
		Total: data.Total,
		Books: make([]Book, len(data.Items)),
	}

	for i := 0; i < len(data.Items); i++ {
		item := data.Items[i]

		books.Books[i] = Book{
			ID:          item.ID,
			Path:        item.Path,
			Title:       item.Title,
			MimeType:    item.MimeType,
			CreatedAt:   item.CreatedAt,
			Purchased:   item.Purchased,
			Bytes:       item.Bytes,
			ClientMtime: item.ClientMtime,
			FastHash:    item.FastHash,
			Favorite:    item.Favorite,
			ReadStatus:  item.ReadStatus,
			Link:        item.Link,
			HasLinks:    item.HasLinks,
			Format:      item.Format,
			Md5Hash:     item.Md5Hash,
			Mtime:       item.Mtime,
			Name:        item.Name,
			ReadPercent: item.ReadPercent,
			Percent:     item.Percent,
			IsDrm:       item.IsDrm,
			IsLcp:       item.IsLcp,
			IsAudioBook: item.IsAudioBook,
			MetaData: BookMetaData{
				Title:       item.Metadata.Title,
				Authors:     item.Metadata.Authors,
				Cover:       mappingBookCovers(item.Metadata.Cover),
				Lang:        item.Metadata.Lang,
				Publisher:   item.Metadata.Publisher,
				Updated:     item.Metadata.Updated,
				Year:        item.Metadata.Year,
				Isbn:        item.Metadata.Isbn,
				BookId:      item.Metadata.BookId,
				FixedLayout: item.Metadata.FixedLayout,
			},
			Position: BookPosition{
				Pointer:    item.Position.Pointer,
				PointerPb:  item.Position.PointerPb,
				Percent:    item.Position.Percent,
				Page:       item.Position.Page,
				PagesTotal: item.Position.PagesTotal,
				Updated:    item.Position.Updated,
				Offs:       item.Position.Offs,
			},
			ReadPosition: BookReadPosition{
				Pointer:    item.ReadPosition.Pointer,
				PointerPb:  item.ReadPosition.PointerPb,
				Percent:    item.ReadPosition.Percent,
				Page:       item.ReadPosition.Page,
				PagesTotal: item.ReadPosition.PagesTotal,
				Updated:    item.ReadPosition.Updated,
				Offs:       item.ReadPosition.Offs,
			},
			Action:     item.Action,
			ActionDate: item.ActionDate,
		}
	}

	return books, nil
}

func mappingBookCovers(covers []struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Path   string `json:"path"`
}) []BookCover {
	if len(covers) == 0 {
		return nil
	}

	cs := make([]BookCover, len(covers))

	for i := 0; i < len(covers); i++ {
		cs[i] = BookCover{
			Width:  covers[i].Width,
			Height: covers[i].Height,
			Path:   covers[i].Path,
		}
	}

	return cs
}
