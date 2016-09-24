package lib

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/hashicorp/go-version"
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

type RootResponse struct {
	Couchdb string `json:"couchdb"`
	Version string `json:"version"`
}

type MembershipResponse struct {
	AllNodes     []string `json:"all_nodes"`
	ClusterNodes []string `json:"cluster_nodes"`
}

func (c *CouchdbClient) getServerVersion() (string, error) {
	data, err := c.request("GET", fmt.Sprintf("%s/", c.baseUri))
	if err != nil {
		return "", err
	}
	var root RootResponse
	err = json.Unmarshal(data, &root)
	if err != nil {
		return "", err
	}
	return root.Version, nil
}

func (c *CouchdbClient) isCouchDbV2() (bool, error) {
	clusteredCouch, err := version.NewConstraint(">= 2.0")
	if err != nil {
		return false, err
	}

	serverVersion, err := c.getServerVersion()
	if err != nil {
		return false, err
	}

	couchDbVersion, err := version.NewVersion(serverVersion)
	if err != nil {
		return false, err
	}

	glog.Infof("relaxing on couch@%s", couchDbVersion)
	//fmt.Printf("relaxing on couch@%s\n", couchDbVersion)
	return clusteredCouch.Check(couchDbVersion), nil
}

func (c *CouchdbClient) getNodeNames() ([]string, error) {
	data, err := c.request("GET", fmt.Sprintf("%s/_membership", c.baseUri))
	if err != nil {
		return nil, err
	}
	var membership MembershipResponse
	err = json.Unmarshal(data, &membership)
	if err != nil {
		return nil, err
	}
	for i, name := range membership.ClusterNodes {
		glog.Infof("node[%d]: %s\n", i, name)
	}
	return membership.ClusterNodes, nil
}

func (c *CouchdbClient) getStatsUrisByNodeName(baseUri string) (map[string]string, error) {
	names, err := c.getNodeNames()
	if err != nil {
		return nil, err
	}
	urisByNodeName := make(map[string]string)
	for _, name := range names {
		urisByNodeName[name] = fmt.Sprintf("%s/_node/%s/_stats", baseUri, name)
	}
	return urisByNodeName, nil
}

func (c *CouchdbClient) getStatsByNodeName(urisByNodeName map[string]string) (StatsByNodeName, error) {
	statsByNodeName := make(StatsByNodeName)
	for name, uri := range urisByNodeName {
		data, err := c.request("GET", uri)
		if err != nil {
			return nil, fmt.Errorf("Error reading couchdb stats: %v", err)
		}

		var stats StatsResponse
		err = json.Unmarshal(data, &stats)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling stats: %v", err)
		}
		if stats.Httpd == (Httpd{}) {
			stats.Httpd = stats.Couchdb.Httpd
		}
		if stats.HttpdRequestMethods == (HttpdRequestMethods{}) {
			stats.HttpdRequestMethods = stats.Couchdb.HttpdRequestMethods
		}
		if stats.HttpdStatusCodes == nil {
			stats.HttpdStatusCodes = stats.Couchdb.HttpdStatusCodes
		}
		statsByNodeName[name] = stats
	}
	return statsByNodeName, nil
}

func (c *CouchdbClient) getStats() (StatsByNodeName, error) {
	isCouchDbV2, err := c.isCouchDbV2()
	if err != nil {
		return nil, err
	}
	if isCouchDbV2 {
		urisByNode, err := c.getStatsUrisByNodeName(c.baseUri)
		if err != nil {
			return nil, err
		}
		return c.getStatsByNodeName(urisByNode)
	} else {
		urisByNode := map[string]string{
			"master": c.statsUri,
		}
		return c.getStatsByNodeName(urisByNode)
	}
}

func (c *CouchdbClient) request(method string, uri string) (respData []byte, err error) {
	req, err := http.NewRequest(method, uri, nil)
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

	respData, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			respData = []byte(err.Error())
		}
		return nil, fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, respData)
	}

	return respData, nil
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
