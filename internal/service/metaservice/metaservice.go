package metaservice

import (
	"net/http"

	"github.com/jonlaing/htmlmeta"
)

func GetMeta(url string) (*htmlmeta.HTMLMeta, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	meta := htmlmeta.Extract(resp.Body)
	return &meta, nil
}
