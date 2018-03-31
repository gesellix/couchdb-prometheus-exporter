package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"net/http"
	"strings"
	"io"
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

type MembershipResponse struct {
	AllNodes     []string `json:"all_nodes"`
	ClusterNodes []string `json:"cluster_nodes"`
}

func (c *CouchdbClient) getNodeInfo() (NodeInfo, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/", c.BaseUri), nil)
	if err != nil {
		return NodeInfo{}, err
	}
	var root NodeInfo
	err = json.Unmarshal(data, &root)
	if err != nil {
		return NodeInfo{}, err
	}
	return root, nil
}

func (c *CouchdbClient) getServerVersion() (string, error) {
	nodeInfo, err := c.getNodeInfo()
	if err != nil {
		return "", err
	}
	return nodeInfo.Version, nil
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

	return clusteredCouch.Check(couchDbVersion), nil
}

func (c *CouchdbClient) GetNodeNames() ([]string, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/_membership", c.BaseUri), nil)
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

func (c *CouchdbClient) getNodeBaseUrisByNodeName(baseUri string) (map[string]string, error) {
	names, err := c.GetNodeNames()
	if err != nil {
		return nil, err
	}
	urisByNodeName := make(map[string]string)
	for _, name := range names {
		urisByNodeName[name] = fmt.Sprintf("%s/_node/%s", baseUri, name)
	}
	return urisByNodeName, nil
}

func (c *CouchdbClient) getStatsByNodeName(urisByNodeName map[string]string) (map[string]StatsResponse, error) {
	statsByNodeName := make(map[string]StatsResponse)
	for name, uri := range urisByNodeName {
		var stats StatsResponse
		data, err := c.Request("GET", fmt.Sprintf("%s/_stats", uri), nil)
		if err != nil {
			err = fmt.Errorf("error reading couchdb stats: %v", err)
			if !strings.Contains(err.Error(), "\"error\":\"nodedown\"") {
				return nil, err
			}

			stats.Up = 0
			glog.Error(fmt.Errorf("continuing despite error: %v", err))
			continue
		}

		stats.Up = 1

		err = json.Unmarshal(data, &stats)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling stats: %v", err)
		}

		nodeInfo, err := c.getNodeInfo()
		if err != nil {
			return nil, err
		}
		stats.NodeInfo = nodeInfo
		statsByNodeName[name] = stats
	}

	if len(urisByNodeName) == 0 {
		return nil, fmt.Errorf("all nodes down")
	}

	return statsByNodeName, nil
}

func (c *CouchdbClient) getStats(databases []string) (Stats, error) {
	isCouchDbV2, err := c.isCouchDbV2()
	if err != nil {
		return Stats{}, err
	}
	if isCouchDbV2 {
		urisByNode, err := c.getNodeBaseUrisByNodeName(c.BaseUri)
		if err != nil {
			return Stats{}, err
		}
		nodeStats, err := c.getStatsByNodeName(urisByNode)
		if err != nil {
			return Stats{}, err
		}
		databaseStats, err := c.getDatabasesStatsByNodeName(databases)
		if err != nil {
			return Stats{}, err
		}
		activeTasks, err := c.getActiveTasks()
		if err != nil {
			return Stats{}, err
		}
		return Stats{
			StatsByNodeName:       nodeStats,
			DatabaseStatsByDbName: databaseStats,
			ActiveTasksResponse:   activeTasks,
			ApiVersion:            "2"}, nil
	} else {
		urisByNode := map[string]string{
			"master": c.BaseUri,
		}
		nodeStats, err := c.getStatsByNodeName(urisByNode)
		if err != nil {
			return Stats{}, err
		}
		databaseStats, err := c.getDatabasesStatsByNodeName(databases)
		if err != nil {
			return Stats{}, err
		}
		activeTasks, err := c.getActiveTasks()
		if err != nil {
			return Stats{}, err
		}
		return Stats{
			StatsByNodeName:       nodeStats,
			DatabaseStatsByDbName: databaseStats,
			ActiveTasksResponse:   activeTasks,
			ApiVersion:            "1"}, nil
	}
}

func (c *CouchdbClient) getDatabasesStatsByNodeName(databases []string) (map[string]DatabaseStats, error) {
	dbStatsByDbName := make(map[string]DatabaseStats)
	for _, dbName := range databases {
		data, err := c.Request("GET", fmt.Sprintf("%s/%s", c.BaseUri, dbName), nil)
		if err != nil {
			return nil, fmt.Errorf("error reading database '%s' stats: %v", dbName, err)
		}

		var dbStats DatabaseStats
		err = json.Unmarshal(data, &dbStats)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling database '%s' stats: %v", dbName, err)
		}
		dbStats.DiskSizeOverhead = dbStats.DiskSize - dbStats.DataSize
		dbStatsByDbName[dbName] = dbStats
	}
	return dbStatsByDbName, nil
}

func (c *CouchdbClient) getActiveTasks() (ActiveTasksResponse, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/_active_tasks", c.BaseUri), nil)
	if err != nil {
		return ActiveTasksResponse{}, fmt.Errorf("error reading active tasks: %v", err)
	}

	var activeTasks ActiveTasksResponse
	err = json.Unmarshal(data, &activeTasks)
	if err != nil {
		return ActiveTasksResponse{}, fmt.Errorf("error unmarshalling active tasks: %v", err)
	}
	for _, activeTask := range activeTasks {
		// CouchDB 1.x doesn't know anything about nodes.
		if activeTask.Node == "" {
			activeTask.Node = "master"
		}
	}
	return activeTasks, nil
}

func (c *CouchdbClient) getDatabaseList() ([]string, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/_all_dbs", c.BaseUri), nil)
	if err != nil {
		return nil, err
	}
	var dbs []string
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (c *CouchdbClient) Request(method string, uri string, body io.Reader) (respData []byte, err error) {
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
		return nil, fmt.Errorf("status %s (%d): %s", resp.Status, resp.StatusCode, respData)
	}

	return respData, nil
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
