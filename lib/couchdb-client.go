package lib

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	insecure = flag.Bool("insecure", true, "Ignore server certificate if using https")
)

type BasicAuth struct {
	Username string
	Password string
}

type CouchdbClient struct {
	baseUri   string
	statsUri  string
	basicAuth BasicAuth
	client    *http.Client
}

func (c *CouchdbClient) getStats() (respData []byte, err error) {
	req, err := http.NewRequest("GET", c.statsUri, nil)
	if err != nil {
		return nil, err
	}
	if len(c.basicAuth.Username) > 0 {
		req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return nil, fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	return data, nil
}

func NewCouchdbClient(uri string, basicAuth BasicAuth) *CouchdbClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		},
	}
	return &CouchdbClient{
		baseUri:   uri,
		statsUri:  fmt.Sprintf("%s/_stats", uri),
		basicAuth: basicAuth,
		client:    httpClient,
	}
}
