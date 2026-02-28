package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) get2(method string, path string, header http.Header, out any) ([]byte, error) {
	req, err := c.newRequest(method, path, nil, header, nil)
	if err != nil {
		return nil, err
	}

	_, resbody, err := c.doRequest(req, out, nil, true)
	if err != nil {
		return nil, err
	}

	return resbody, nil
}

func (c *Client) post2(method string, path string, header http.Header, in any, out any) ([]byte, error) {
	if in == nil {
		in = struct{}{}
	}

	reqbody, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := c.newRequest(method, path, nil, header, bytes.NewReader(reqbody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, resbody, err := c.doRequest(req, out, reqbody, true)
	if err != nil {
		return nil, err
	}

	return resbody, nil
}

func (c *Client) Get2(path string, header http.Header, out any) ([]byte, error) {
	return c.get2(http.MethodGet, path, header, out)
}

func (c *Client) Post2(path string, header http.Header, in any, out any) ([]byte, error) {
	return c.post2(http.MethodPost, path, header, in, out)
}

func (c *Client) Put2(path string, header http.Header, in any, out any) ([]byte, error) {
	return c.post2(http.MethodPut, path, header, in, out)
}

func (c *Client) Patch2(path string, header http.Header, in any, out any) ([]byte, error) {
	return c.post2(http.MethodPatch, path, header, in, out)
}

func (c *Client) Delete2(path string, header http.Header, out any) ([]byte, error) {
	return c.get2(http.MethodDelete, path, header, out)
}

func (c *Client) post3(method string, path string, header http.Header, in []byte, out any) ([]byte, error) {
	req, err := c.newRequest(method, path, nil, header, bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, resbody, err := c.doRequest(req, out, in, true)
	if err != nil {
		return nil, err
	}

	return resbody, nil
}

func (c *Client) Post3(path string, header http.Header, in []byte, out any) ([]byte, error) {
	return c.post3(http.MethodPost, path, header, in, out)
}

func (c *Client) Put3(path string, header http.Header, in []byte, out any) ([]byte, error) {
	return c.post3(http.MethodPut, path, header, in, out)
}

func (c *Client) Patch3(path string, header http.Header, in []byte, out any) ([]byte, error) {
	return c.post3(http.MethodPatch, path, header, in, out)
}
