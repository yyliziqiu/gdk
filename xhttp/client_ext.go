package xhttp

import (
	"bytes"
	"net/http"
	"net/url"
	"path/filepath"
)

func (cli *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, field string, filepath string, out interface{}) error {
	files := map[string]string{field: filepath}
	return cli.PostData(path, query, header, values, files, out)
}

func (cli *Client) ForwardBinary(path string, query url.Values, header http.Header, src string, out interface{}) error {
	data, typ, err := cli.GetBinary(src, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostBinary(path, query, header, typ, bytes.NewReader(data), out)
}

func (cli *Client) ForwardStream(path string, query url.Values, header http.Header, values map[string]string, field string, mimeTyp string, src string, out interface{}) error {
	data, _, err := cli.GetBinary(src, nil, nil)
	if err != nil {
		return err
	}
	return cli.PostStream(path, query, header, values, field, filepath.Base(src), mimeTyp, bytes.NewReader(data), out)
}
