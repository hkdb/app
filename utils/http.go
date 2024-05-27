package utils

import (
	"net/http"
	"path"
)

func GetFileFromUrl(url string) (string, string) {

	r, _ := http.Get(url)
	file := path.Base(r.Request.URL.Path)

	return r.Request.URL.String(), file
}
