package pocketbook_cloud_client

import (
	"fmt"
	"net/http"
)

type httpStatusError struct {
	code int
}

func (e httpStatusError) Error() string {
	return fmt.Sprintf("http status code: %d %s", e.code, http.StatusText(e.code))
}

func (e httpStatusError) Code() int {
	return e.code
}
