package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-version"
	"k8s.io/klog"
)

type BasicAuth struct {
	Username string
	Password string
}

type CouchdbClient struct {
	BaseUri           string
	basicAuth         BasicAuth
	client            *http.Client
	ResetRequestCount func()
	GetRequestCount   func() int
}

type MembershipResponse struct {
	AllNodes     []string `json:"all_nodes"`
	ClusterNodes []string `json:"cluster_nodes"`
}

func (c *CouchdbClient) getHealthByNodeName(urisByNodeName map[string]string) (map[string]HealthResponse, error) {
	healthByNodeName := make(map[string]HealthResponse)
	for name, uri := range urisByNodeName {
		var health HealthResponse

		// Instead of http://localhost:15984/_node/couchdb@172.16.238.13/_up
		// we have to request http://localhost:15984/_up,
		// but how do we find the correct base uris, say: ip and port, for each node?
		data, err := c.Request("GET", fmt.Sprintf("%s/_up", uri), nil)
		if err != nil {
			err = fmt.Errorf("error reading couchdb health: %v", err)
			//if !strings.Contains(err.Error(), "\"error\":\"nodedown\"") {
			//	return nil, err
			//}
			klog.Error(fmt.Errorf("continuing despite error: %v", err))
			continue
		}

		err = json.Unmarshal(data, &health)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling health: %v", err)
		}

		healthByNodeName[name] = health
	}

	if len(urisByNodeName) == 0 {
		return nil, fmt.Errorf("all nodes down")
	}

	return healthByNodeName, nil
}

func (c *CouchdbClient) getNodeInfo(uri string) (NodeInfo, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/", uri), nil)
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
	nodeInfo, err := c.getNodeInfo(c.BaseUri)
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
	// for i, name := range membership.ClusterNodes {
	// 	klog.Infof("node[%d]: %s\n", i, name)
	// }
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
			klog.Error(fmt.Errorf("continuing despite error: %v", err))
			continue
		}

		stats.Up = 1

		err = json.Unmarshal(data, &stats)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling stats: %v", err)
		}

		// TODO this one is expected to retrieve other nodes' info
		nodeInfo, err := c.getNodeInfo(c.BaseUri)
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

func (c *CouchdbClient) getSystemByNodeName(urisByNodeName map[string]string) (map[string]SystemResponse, error) {
	systemByNodeName := make(map[string]SystemResponse)
	for name, uri := range urisByNodeName {
		var stats SystemResponse

		data, err := c.Request("GET", fmt.Sprintf("%s/_system", uri), nil)
		if err != nil {
			err = fmt.Errorf("error reading couchdb system stats: %v", err)
			if !strings.Contains(err.Error(), "\"error\":\"nodedown\"") {
				return nil, err
			}
			klog.Error(fmt.Errorf("continuing despite error: %v", err))
			continue
		}

		err = json.Unmarshal(data, &stats)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling stats: %v", err)
		}

		systemByNodeName[name] = stats
	}

	if len(urisByNodeName) == 0 {
		return nil, fmt.Errorf("all nodes down")
	}

	return systemByNodeName, nil
}

func (c *CouchdbClient) getStats(config CollectorConfig) (Stats, error) {
	isCouchDbV2, err := c.isCouchDbV2()
	if err != nil {
		return Stats{}, err
	}
	if isCouchDbV2 {
		urisByNode, err := c.getNodeBaseUrisByNodeName(c.BaseUri)
		if err != nil {
			return Stats{}, err
		}
		//nodeHealth, err := c.getHealthByNodeName(urisByNode)
		//if err != nil {
		//	return Stats{}, err
		//}
		nodeStats, err := c.getStatsByNodeName(urisByNode)
		if err != nil {
			return Stats{}, err
		}
		databaseStats, err := c.getDatabasesStatsByDbName(config.ObservedDatabases)
		if err != nil {
			return Stats{}, err
		}
		if config.CollectViews {
			err := c.enhanceWithViewUpdateSeq(databaseStats)
			if err != nil {
				return Stats{}, err
			}
		}
		schedulerJobs := SchedulerJobsResponse{}
		if config.CollectSchedulerJobs {
			schedulerJobs, err = c.getSchedulerJobs()
		}
		activeTasks, err := c.getActiveTasks()
		if err != nil {
			return Stats{}, err
		}
		databasesList, err := c.getDatabaseList()
		if err != nil {
			return Stats{}, err
		}
		systemStats, err := c.getSystemByNodeName(urisByNode)
		if err != nil {
			return Stats{}, err
		}

		return Stats{
			//HealthByNodeName:      nodeHealth,
			StatsByNodeName:       nodeStats,
			DatabasesTotal:        len(databasesList),
			DatabaseStatsByDbName: databaseStats,
			ActiveTasksResponse:   activeTasks,
			SchedulerJobsResponse: schedulerJobs,
			SystemByNodeName:      systemStats,
			ApiVersion:            "2"}, nil
	} else {
		urisByNode := map[string]string{
			"master": c.BaseUri,
		}
		nodeStats, err := c.getStatsByNodeName(urisByNode)
		if err != nil {
			return Stats{}, err
		}
		databaseStats, err := c.getDatabasesStatsByDbName(config.ObservedDatabases)
		if err != nil {
			return Stats{}, err
		}
		if config.CollectViews {
			err := c.enhanceWithViewUpdateSeq(databaseStats)
			if err != nil {
				return Stats{}, err
			}
		}
		activeTasks, err := c.getActiveTasks()
		if err != nil {
			return Stats{}, err
		}
		databasesList, err := c.getDatabaseList()
		if err != nil {
			return Stats{}, err
		}
		return Stats{
			StatsByNodeName:       nodeStats,
			DatabasesTotal:        len(databasesList),
			DatabaseStatsByDbName: databaseStats,
			ActiveTasksResponse:   activeTasks,
			ApiVersion:            "1"}, nil
	}
}

func (c *CouchdbClient) getDatabasesStatsByDbName(databases []string) (map[string]DatabaseStats, error) {
	dbStatsByDbName := make(map[string]DatabaseStats)
	type result struct {
		dbName  string
		dbStats DatabaseStats
		err     error
	}
	r := make(chan result, len(databases))
	// scatter
	for _, dbName := range databases {
		dbName := dbName // rebind for closure to capture the value
		go func() {
			var dbStats DatabaseStats
			data, err := c.Request("GET", fmt.Sprintf("%s/%s", c.BaseUri, dbName), nil)
			if err != nil {
				r <- result{err: fmt.Errorf("error reading database '%s' stats: %v", dbName, err)}
				return
			}

			err = json.Unmarshal(data, &dbStats)
			if err != nil {
				r <- result{err: fmt.Errorf("error unmarshalling database '%s' stats: %v", dbName, err)}
				return
			}
			dbStats.DiskSizeOverhead = dbStats.DiskSize - dbStats.DataSize
			if dbStats.CompactRunningBool {
				dbStats.CompactRunning = 1
			} else {
				dbStats.CompactRunning = 0
			}
			r <- result{dbName, dbStats, nil}
		}()
	}
	// gather
	for range databases {
		res := <-r
		if res.err != nil {
			return nil, res.err
		}
		dbStatsByDbName[res.dbName] = res.dbStats
	}
	return dbStatsByDbName, nil
}

func (c *CouchdbClient) enhanceWithViewUpdateSeq(dbStatsByDbName map[string]DatabaseStats) error {
	type result struct {
		dbName  string
		dbStats DatabaseStats
		err     error
	}
	r := make(chan result, len(dbStatsByDbName))
	// scatter
	for dbName, dbStats := range dbStatsByDbName {
		dbName := dbName // rebind for closure to capture the value
		dbStats := dbStats
		go func() {
			query := strings.Join([]string{
				"startkey=\"_design/\"",
				"endkey=\"_design0\"",
				"include_docs=true",
			}, "&")
			designDocData, err := c.Request("GET", fmt.Sprintf("%s/%s/_all_docs?%s", c.BaseUri, dbName, query), nil)
			if err != nil {
				r <- result{err: fmt.Errorf("error reading database '%s' stats: %v", dbName, err)}
				return
			}

			var designDocs DocsResponse
			err = json.Unmarshal(designDocData, &designDocs)
			if err != nil {
				r <- result{err: fmt.Errorf("error unmarshalling design docs for database '%s': %v", dbName, err)}
				return
			}
			views := make(ViewStatsByDesignDocName)
			for _, row := range designDocs.Rows {
				updateSeqByView := make(ViewStats)
				type viewresult struct {
					viewName  string
					updateSeq string
					err       error
				}
				v := make(chan viewresult, len(row.Doc.Views))
				for viewName := range row.Doc.Views {
					viewName := viewName
					go func() {
						//klog.Infof("/%s/%s/_view/%s\n", dbName, row.Doc.Id, viewName)
						query := strings.Join([]string{
							"stale=ok",
							"update=false",
							"stable=true",
							"update_seq=true",
							"include_docs=false",
							"limit=0",
						}, "&")
						var viewDoc ViewResponse
						viewDocData, err := c.Request("GET", fmt.Sprintf("%s/%s/%s/_view/%s?%s", c.BaseUri, dbName, row.Doc.Id, viewName, query), nil)
						err = json.Unmarshal(viewDocData, &viewDoc)
						if err != nil {
							v <- viewresult{err: fmt.Errorf("error unmarshalling view doc for view '%s/%s/_view/%s': %v", dbName, row.Doc.Id, viewName, err)}
							return
						}
						v <- viewresult{viewName, viewDoc.UpdateSeq.String(), nil}
					}()
				}
				for range row.Doc.Views {
					res := <-v
					if res.err != nil {
						r <- result{err: res.err}
						return
					}
					updateSeqByView[res.viewName] = res.updateSeq
				}
				views[row.Doc.Id] = updateSeqByView
			}
			dbStats.Views = views
			r <- result{dbName, dbStats, nil}
		}()
	}
	// gather
	for range dbStatsByDbName {
		resp := <-r
		dbName, dbStats, err := resp.dbName, resp.dbStats, resp.err
		if err != nil {
			return err
		}
		dbStatsByDbName[dbName] = dbStats
	}
	return nil
}

// CouchDB 2.x+ only
func (c *CouchdbClient) getSchedulerJobs() (SchedulerJobsResponse, error) {
	data, err := c.Request("GET", fmt.Sprintf("%s/_scheduler/jobs", c.BaseUri), nil)
	if err != nil {
		return SchedulerJobsResponse{}, fmt.Errorf("error reading scheduler jobs: %v", err)
	}

	var schedulerJobs SchedulerJobsResponse
	err = json.Unmarshal(data, &schedulerJobs)
	if err != nil {
		return SchedulerJobsResponse{}, fmt.Errorf("error unmarshalling scheduler jobs: %v", err)
	}
	//for _, job := range schedulerJobs.Jobs {
	//	replDoc, err := c.Request("GET", fmt.Sprintf("%s/%s/%s", c.BaseUri, job.Database, job.DocID), nil)
	//	if err != nil {
	//		return SchedulerJobsResponse{}, fmt.Errorf("error reading replication doc '%s/%s': %v", job.Database, job.DocID, err)
	//	}
	//}
	return schedulerJobs, nil
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
	data, err := c.Request("GET", fmt.Sprintf("%s/%s", c.BaseUri, AllDbs), nil)
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
	if resp != nil {
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				err = cerr
			}
		}()
	}

	respData, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			respData = []byte(err.Error())
		}
		return nil, fmt.Errorf("status %s (%d): %s", resp.Status, resp.StatusCode, respData)
	}
	if err != nil {
		return nil, err
	}
	return respData, nil
}

type requestCountingRoundTripper struct {
	RequestCount int64
	rt           http.RoundTripper
}

func (rt *requestCountingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt64(&rt.RequestCount, 1)
	//klog.Infof("req[%d] %s", atomic.LoadInt64(&rt.RequestCount), req.URL.String())
	return rt.rt.RoundTrip(req)
}

func NewCouchdbClient(uri string, basicAuth BasicAuth, insecure bool) *CouchdbClient {
	countingRoundTripper := &requestCountingRoundTripper{
		0,
		&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	httpClient := &http.Client{
		Transport: countingRoundTripper,
	}

	return &CouchdbClient{
		BaseUri:   uri,
		basicAuth: basicAuth,
		client:    httpClient,
		ResetRequestCount: func() {
			atomic.StoreInt64(&countingRoundTripper.RequestCount, 0)
		},
		GetRequestCount: func() int {
			return int(atomic.LoadInt64(&countingRoundTripper.RequestCount))
		},
	}
}
