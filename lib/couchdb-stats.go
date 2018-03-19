package lib

type Counter struct {
	// v1.x api
	Description string
	Current     float64
	// v2.x api
	Value float64
	Type  string
	Desc  string
}

// v2.x api
type HistogramValue struct {
	Min               float64 `json:"min"`
	Max               float64 `json:"max"`
	ArithmeticMean    float64 `json:"arithmetic_mean"`
	GeometricMean     float64 `json:"geometric_mean"`
	HarmonicMean      float64 `json:"harmonic_mean"`
	Median            float64 `json:"median"`
	Variance          float64 `json:"variance"`
	StandardDeviation float64 `json:"standard_deviation"`
	Skewness          float64 `json:"skewness"`
	Kurtosis          float64 `json:"kurtosis"`
	N                 float64 `json:"n"`
}

type Histogram struct {
	// v1.x api
	Description string
	Current     float64 `json:"current"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	Mean        float64 `json:"mean"`
	Sum         float64 `json:"sum"`
	Stddev      float64 `json:"stddev"`

	// v2.x api
	Value HistogramValue
	Type  string
	Desc  string
}

type CouchdbStats struct {
	// v1.x, and v2.x api
	AuthCacheHits   Counter   `json:"auth_cache_hits"`
	AuthCacheMisses Counter   `json:"auth_cache_misses"`
	DatabaseReads   Counter   `json:"database_reads"`
	DatabaseWrites  Counter   `json:"database_writes"`
	OpenDatabases   Counter   `json:"open_databases"`
	OpenOsFiles     Counter   `json:"open_os_files"`
	RequestTime     Histogram `json:"request_time"`
	// v2.x api
	Httpd               Httpd               `json:"httpd"`
	HttpdRequestMethods HttpdRequestMethods `json:"httpd_request_methods"`
	HttpdStatusCodes    HttpdStatusCodes    `json:"httpd_status_codes"`
}

type HttpdRequestMethods struct {
	COPY   Counter `json:"COPY"`
	DELETE Counter `json:"DELETE"`
	GET    Counter `json:"GET"`
	HEAD   Counter `json:"HEAD"`
	POST   Counter `json:"POST"`
	PUT    Counter `json:"PUT"`
}

type HttpdStatusCodes map[string]Counter

type Httpd struct {
	BulkRequests             Counter `json:"bulk_requests"`
	ClientsRequestingChanges Counter `json:"clients_requesting_changes"`
	Requests                 Counter `json:"requests"`
	TemporaryViewReads       Counter `json:"temporary_view_reads"`
	ViewReads                Counter `json:"view_reads"`
}

type StatsResponse struct {
	Couchdb CouchdbStats `json:"couchdb"`
	Up      float64      `json:"-"`
	// v1.x api
	Httpd               Httpd               `json:"httpd"`
	HttpdRequestMethods HttpdRequestMethods `json:"httpd_request_methods"`
	HttpdStatusCodes    HttpdStatusCodes    `json:"httpd_status_codes"`
}

type DatabaseStats struct {
	DiskSize         float64 `json:"disk_size"`
	DataSize         float64 `json:"data_size"`
	DiskSizeOverhead float64
}

type DatabaseStatsByDbName map[string]DatabaseStats

type ActiveTask struct {
	Type       string  `json:"type"`
	Node       string  `json:"node,omitempty"`
	Continuous bool    `json:"continuous,omitempty"`
	UpdatedOn  float64 `json:"updated_on,omitempty"`
	Source     string  `json:"source,omitempty"`
	Target     string  `json:"target,omitempty"`
	DocId      string  `json:"doc_id,omitempty"`
}

type ActiveTasksResponse []ActiveTask

type Stats struct {
	StatsByNodeName       map[string]StatsResponse
	DatabaseStatsByDbName DatabaseStatsByDbName
	ActiveTasksResponse   ActiveTasksResponse
	ApiVersion            string
}
