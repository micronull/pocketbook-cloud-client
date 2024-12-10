package pocketbook_cloud_client_test

import (
	"embed"
)

//go:embed testdata
var testdata embed.FS

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func isAllTrue(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}

	return true
}
