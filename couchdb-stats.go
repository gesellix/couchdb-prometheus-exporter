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
	AuthCacheHits   StatsDetail `json:"auth_cache_hits"`
	AuthCacheMisses StatsDetail `json:"auth_cache_misses"`
	DatabaseReads   StatsDetail `json:"database_reads"`
	DatabaseWrites  StatsDetail `json:"database_writes"`
	OpenDatabases   StatsDetail `json:"open_databases"`
	OpenOsFiles     StatsDetail `json:"open_os_files"`
	RequestTime     StatsDetail `json:"request_time"`
}

type HttpdRequestMethods struct {
	COPY   StatsDetail
	DELETE StatsDetail
	GET    StatsDetail
	HEAD   StatsDetail
	POST   StatsDetail
	PUT    StatsDetail
}

type HttpdStatusCodes map[string]StatsDetail

type Httpd struct {
	BulkRequests             StatsDetail `json:"bulk_requests"`
	ClientsRequestingChanges StatsDetail `json:"clients_requesting_changes"`
	Requests                 StatsDetail `json:"requests"`
	TemporaryViewReads       StatsDetail `json:"temporary_view_reads"`
	ViewReads                StatsDetail `json:"view_reads"`
}

type StatsResponse struct {
	Couchdb             CouchdbStats        `json:"couchdb"`
	HttpdRequestMethods HttpdRequestMethods `json:"httpd_request_methods"`
	HttpdStatusCodes    HttpdStatusCodes    `json:"httpd_status_codes"`
	Httpd               Httpd               `json:"httpd"`
}
