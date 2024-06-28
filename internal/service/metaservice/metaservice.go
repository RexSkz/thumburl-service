package metaservice

import (
	"net/http"

	"github.com/jonlaing/htmlmeta"
)

func GetMeta(url string) (*htmlmeta.HTMLMeta, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	meta := htmlmeta.Extract(resp.Body)
	return &meta, nil
}
