package utils

import (
	"net/http"
	"path"
)

func GetFileFromUrl(url string) string {

	r, _ := http.NewRequest("GET", url, nil)
	file := path.Base(r.URL.Path)

	return file

}
