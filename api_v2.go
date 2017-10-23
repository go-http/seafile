package seafile

import (
	"io"
	"net/http"
)

func (cli *Client) doRequest(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	return cli.requestApi("/api2", method, uri, header, body)
}
