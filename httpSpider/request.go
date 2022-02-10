package httpSpider

import (
	"context"
	"io"
	"net/http"
)

// 在原有功能上增加修改header
func (hs *HttpSpider) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	req.Header = defaultHeader
	return req, err
}
