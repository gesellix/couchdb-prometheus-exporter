package cluster_config

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type BasicAuth struct {
	Username string
	Password string
}

type CouchdbClient struct {
	BaseUri   string
	basicAuth BasicAuth
	client    *http.Client
}

func (c *CouchdbClient) Request(method string, uri string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
	}
	if len(c.basicAuth.Username) > 0 {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}

	fmt.Printf("[REQ] %s %s@%v\n", req.Method, c.basicAuth.Username, req.URL)
	resp, err = c.client.Do(req)
	if err != nil {
		fmt.Printf("[RES-ERR] %v\n", err)
		return nil, err
	}
	fmt.Printf("[RES-OK] %v\n", resp)
	return resp, nil
}

func (c *CouchdbClient) RequestBody(method string, uri string, body io.Reader) (respBody []byte, err error) {
	resp, err := c.Request(method, uri, body)
	respBody, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			respBody = []byte(err.Error())
		}
		return nil, fmt.Errorf("status %s (%d): %s", resp.Status, resp.StatusCode, respBody)
	}

	return respBody, nil
}

func NewCouchdbClient(uri string, basicAuth BasicAuth, insecure bool) *CouchdbClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	return &CouchdbClient{
		BaseUri:   uri,
		basicAuth: basicAuth,
		client:    httpClient,
	}
}
