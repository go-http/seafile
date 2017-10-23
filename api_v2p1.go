package seafile

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func (cli *Client) apiRequestV2p1(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	return cli.requestApi("/api/v2.1", method, uri, header, body)
}

func (cli *Client) apiGET(uri string) (*http.Response, error) {
	return cli.apiRequestV2p1("GET", uri, nil, nil)
}

func (cli *Client) apiPOSTForm(uri string, form url.Values) (*http.Response, error) {
	header := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}
	return cli.apiRequestV2p1("POST", uri, header, bytes.NewBufferString(form.Encode()))
}

func (cli *Client) apiPOST(uri string, header http.Header, body io.Reader) (*http.Response, error) {
	return cli.apiRequestV2p1("POST", uri, header, body)
}

func (cli *Client) apiPUT(uri string, header http.Header, body io.Reader) (*http.Response, error) {
	return cli.apiRequestV2p1("PUT", uri, header, body)
}

func (cli *Client) apiDELETE(uri string) (*http.Response, error) {
	return cli.apiRequestV2p1("DELETE", uri, nil, nil)
}
