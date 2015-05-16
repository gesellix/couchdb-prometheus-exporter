package main

type StatsDetail struct {
	Description string
	Current     float64
	Sum         float64
	Mean        float64
	Stddev      float64
	Min         float64
	Max         float64
}

type CouchdbStats struct {
	AuthCacheMisses StatsDetail `json:"auth_cache_misses"`
	AuthCacheHits   StatsDetail `json:"auth_cache_hits"`
	DatabaseReads   StatsDetail `json:"database_reads"`
	DatabaseWrites  StatsDetail `json:"database_writes"`
	OpenDatabases   StatsDetail `json:"open_databases"`
	OpenOsFiles     StatsDetail `json:"open_os_files"`
	RequestTime     StatsDetail `json:"request_time"`
}

type HttpdRequestMethods struct {
	COPY   float64
	DELETE float64
	GET    float64
	HEAD   float64
	POST   float64
	PUT    float64
}

type HttpdStatusCodes map[string]StatsDetail

type Httpd struct {
	ClientsRequestingChanges StatsDetail `json:"clients_requesting_changes"`
	TemporaryViewReads       StatsDetail `json:"temporary_view_reads"`
	Requests                 StatsDetail `json:"requests"`
	BulkRequests             StatsDetail `json:"bulk_requests"`
	ViewReads                StatsDetail `json:"view_reads"`
}

type StatsResponse struct {
	Couchdb             CouchdbStats        `json:"couchdb"`
	HttpdRequestMethods HttpdRequestMethods `json:"httpd_request_methods"`
	HttpdStatusCodes    HttpdStatusCodes    `json:"httpd_status_codes"`
	Httpd               Httpd               `json:"httpd"`
}
