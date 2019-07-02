# CouchDB Prometheus Exporter

## Metrics Overview
The following list gives you an overview on the currently exposed metrics.
Please note that beyond the complete CouchDB stats a metric `couchdb_up` has been
added as minimal connection health check.

The `node_name` label defaults to `master` for CouchDB 1.x and `nonode@nohost` for a single CouchDB 2.x setup.
The node name makes more sense in a clustered CouchDB 2.x environment with several nodes to provide separated stats for each node. 

    # HELP couchdb_database_compact_running database compaction running
    # TYPE couchdb_database_compact_running gauge
    couchdb_database_compact_running{db_name="_global_changes"} 0
    couchdb_database_compact_running{db_name="_replicator"} 0
    couchdb_database_compact_running{db_name="_users"} 0
    couchdb_database_compact_running{db_name="db-1"} 0
    couchdb_database_compact_running{db_name="db-2"} 0
    # HELP couchdb_database_data_size data size
    # TYPE couchdb_database_data_size gauge
    couchdb_database_data_size{db_name="_global_changes"} 8005
    couchdb_database_data_size{db_name="_replicator"} 3940
    couchdb_database_data_size{db_name="_users"} 3858
    couchdb_database_data_size{db_name="db-1"} 0
    couchdb_database_data_size{db_name="db-2"} 0
    # HELP couchdb_database_disk_size disk size
    # TYPE couchdb_database_disk_size gauge
    couchdb_database_disk_size{db_name="_global_changes"} 370100
    couchdb_database_disk_size{db_name="_replicator"} 70958
    couchdb_database_disk_size{db_name="_users"} 70958
    couchdb_database_disk_size{db_name="db-1"} 34024
    couchdb_database_disk_size{db_name="db-2"} 34024
    # HELP couchdb_database_doc_count document count
    # TYPE couchdb_database_doc_count gauge
    couchdb_database_doc_count{db_name="_global_changes"} 6
    couchdb_database_doc_count{db_name="_replicator"} 1
    couchdb_database_doc_count{db_name="_users"} 1
    couchdb_database_doc_count{db_name="db-1"} 0
    couchdb_database_doc_count{db_name="db-2"} 0
    # HELP couchdb_database_doc_del_count deleted document count
    # TYPE couchdb_database_doc_del_count gauge
    couchdb_database_doc_del_count{db_name="_global_changes"} 0
    couchdb_database_doc_del_count{db_name="_replicator"} 0
    couchdb_database_doc_del_count{db_name="_users"} 0
    couchdb_database_doc_del_count{db_name="db-1"} 0
    couchdb_database_doc_del_count{db_name="db-2"} 0
    # HELP couchdb_database_overhead disk size overhead
    # TYPE couchdb_database_overhead gauge
    couchdb_database_overhead{db_name="_global_changes"} 362095
    couchdb_database_overhead{db_name="_replicator"} 67018
    couchdb_database_overhead{db_name="_users"} 67100
    couchdb_database_overhead{db_name="db-1"} 34024
    couchdb_database_overhead{db_name="db-2"} 34024
    # HELP couchdb_erlang_memory_atom erlang memory counters - atom
    # TYPE couchdb_erlang_memory_atom gauge
    couchdb_erlang_memory_atom{node_name="couchdb@172.16.238.11"} 504433
    couchdb_erlang_memory_atom{node_name="couchdb@172.16.238.12"} 504433
    couchdb_erlang_memory_atom{node_name="couchdb@172.16.238.13"} 504433
    # HELP couchdb_erlang_memory_atom_used erlang memory counters - atom_used
    # TYPE couchdb_erlang_memory_atom_used gauge
    couchdb_erlang_memory_atom_used{node_name="couchdb@172.16.238.11"} 487705
    couchdb_erlang_memory_atom_used{node_name="couchdb@172.16.238.12"} 477916
    couchdb_erlang_memory_atom_used{node_name="couchdb@172.16.238.13"} 479589
    # HELP couchdb_erlang_memory_binary erlang memory counters - binary
    # TYPE couchdb_erlang_memory_binary gauge
    couchdb_erlang_memory_binary{node_name="couchdb@172.16.238.11"} 312944
    couchdb_erlang_memory_binary{node_name="couchdb@172.16.238.12"} 273120
    couchdb_erlang_memory_binary{node_name="couchdb@172.16.238.13"} 173704
    # HELP couchdb_erlang_memory_code erlang memory counters - code
    # TYPE couchdb_erlang_memory_code gauge
    couchdb_erlang_memory_code{node_name="couchdb@172.16.238.11"} 1.1470019e+07
    couchdb_erlang_memory_code{node_name="couchdb@172.16.238.12"} 1.1022704e+07
    couchdb_erlang_memory_code{node_name="couchdb@172.16.238.13"} 1.1068815e+07
    # HELP couchdb_erlang_memory_ets erlang memory counters - ets
    # TYPE couchdb_erlang_memory_ets gauge
    couchdb_erlang_memory_ets{node_name="couchdb@172.16.238.11"} 1.560672e+06
    couchdb_erlang_memory_ets{node_name="couchdb@172.16.238.12"} 1.521624e+06
    couchdb_erlang_memory_ets{node_name="couchdb@172.16.238.13"} 1.593168e+06
    # HELP couchdb_erlang_memory_other erlang memory counters - other
    # TYPE couchdb_erlang_memory_other gauge
    couchdb_erlang_memory_other{node_name="couchdb@172.16.238.11"} 1.7356164e+07
    couchdb_erlang_memory_other{node_name="couchdb@172.16.238.12"} 1.7199255e+07
    couchdb_erlang_memory_other{node_name="couchdb@172.16.238.13"} 1.7248264e+07
    # HELP couchdb_erlang_memory_processes erlang memory counters - processes
    # TYPE couchdb_erlang_memory_processes gauge
    couchdb_erlang_memory_processes{node_name="couchdb@172.16.238.11"} 1.1010032e+07
    couchdb_erlang_memory_processes{node_name="couchdb@172.16.238.12"} 1.0854232e+07
    couchdb_erlang_memory_processes{node_name="couchdb@172.16.238.13"} 1.1285792e+07
    # HELP couchdb_erlang_memory_processes_used erlang memory counters - processes_used
    # TYPE couchdb_erlang_memory_processes_used gauge
    couchdb_erlang_memory_processes_used{node_name="couchdb@172.16.238.11"} 1.098528e+07
    couchdb_erlang_memory_processes_used{node_name="couchdb@172.16.238.12"} 1.0838344e+07
    couchdb_erlang_memory_processes_used{node_name="couchdb@172.16.238.13"} 1.125456e+07
    # HELP couchdb_exporter_request_count Number of CouchDB requests for this scrape.
    # TYPE couchdb_exporter_request_count gauge
    couchdb_exporter_request_count 25
    # HELP couchdb_fabric_doc_update doc update metrics
    # TYPE couchdb_fabric_doc_update gauge
    couchdb_fabric_doc_update{metric="errors",node_name="couchdb@172.16.238.11"} 1
    couchdb_fabric_doc_update{metric="errors",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_doc_update{metric="errors",node_name="couchdb@172.16.238.13"} 2
    couchdb_fabric_doc_update{metric="mismatched_errors",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_doc_update{metric="mismatched_errors",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_doc_update{metric="mismatched_errors",node_name="couchdb@172.16.238.13"} 0
    couchdb_fabric_doc_update{metric="write_quorum_errors",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_doc_update{metric="write_quorum_errors",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_doc_update{metric="write_quorum_errors",node_name="couchdb@172.16.238.13"} 2
    # HELP couchdb_fabric_open_shard open_shard metrics
    # TYPE couchdb_fabric_open_shard gauge
    couchdb_fabric_open_shard{metric="timeout",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_open_shard{metric="timeout",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_open_shard{metric="timeout",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_fabric_read_repairs read repair metrics
    # TYPE couchdb_fabric_read_repairs gauge
    couchdb_fabric_read_repairs{metric="failure",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_read_repairs{metric="failure",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_read_repairs{metric="failure",node_name="couchdb@172.16.238.13"} 0
    couchdb_fabric_read_repairs{metric="success",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_read_repairs{metric="success",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_read_repairs{metric="success",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_fabric_worker worker metrics
    # TYPE couchdb_fabric_worker gauge
    couchdb_fabric_worker{metric="timeout",node_name="couchdb@172.16.238.11"} 0
    couchdb_fabric_worker{metric="timeout",node_name="couchdb@172.16.238.12"} 0
    couchdb_fabric_worker{metric="timeout",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_auth_cache_hits number of authentication cache hits
    # TYPE couchdb_httpd_auth_cache_hits gauge
    couchdb_httpd_auth_cache_hits{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_auth_cache_hits{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_auth_cache_hits{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_auth_cache_misses number of authentication cache misses
    # TYPE couchdb_httpd_auth_cache_misses gauge
    couchdb_httpd_auth_cache_misses{node_name="couchdb@172.16.238.11"} 2
    couchdb_httpd_auth_cache_misses{node_name="couchdb@172.16.238.12"} 1
    couchdb_httpd_auth_cache_misses{node_name="couchdb@172.16.238.13"} 1
    # HELP couchdb_httpd_bulk_requests number of bulk requests
    # TYPE couchdb_httpd_bulk_requests gauge
    couchdb_httpd_bulk_requests{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_bulk_requests{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_bulk_requests{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_clients_requesting_changes number of clients for continuous _changes
    # TYPE couchdb_httpd_clients_requesting_changes gauge
    couchdb_httpd_clients_requesting_changes{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_clients_requesting_changes{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_clients_requesting_changes{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_database_reads number of times a document was read from a database
    # TYPE couchdb_httpd_database_reads gauge
    couchdb_httpd_database_reads{node_name="couchdb@172.16.238.11"} 53
    couchdb_httpd_database_reads{node_name="couchdb@172.16.238.12"} 52
    couchdb_httpd_database_reads{node_name="couchdb@172.16.238.13"} 53
    # HELP couchdb_httpd_database_writes number of times a database was changed
    # TYPE couchdb_httpd_database_writes gauge
    couchdb_httpd_database_writes{node_name="couchdb@172.16.238.11"} 90
    couchdb_httpd_database_writes{node_name="couchdb@172.16.238.12"} 90
    couchdb_httpd_database_writes{node_name="couchdb@172.16.238.13"} 90
    # HELP couchdb_httpd_databases_total Total number of databases in the cluster
    # TYPE couchdb_httpd_databases_total gauge
    couchdb_httpd_databases_total 5
    # HELP couchdb_httpd_node_up Is the node available.
    # TYPE couchdb_httpd_node_up gauge
    couchdb_httpd_node_up{node_name="couchdb@172.16.238.11"} 1
    couchdb_httpd_node_up{node_name="couchdb@172.16.238.12"} 1
    couchdb_httpd_node_up{node_name="couchdb@172.16.238.13"} 1
    # HELP couchdb_httpd_open_databases number of open databases
    # TYPE couchdb_httpd_open_databases gauge
    couchdb_httpd_open_databases{node_name="couchdb@172.16.238.11"} 21
    couchdb_httpd_open_databases{node_name="couchdb@172.16.238.12"} 21
    couchdb_httpd_open_databases{node_name="couchdb@172.16.238.13"} 22
    # HELP couchdb_httpd_open_os_files number of file descriptors CouchDB has open
    # TYPE couchdb_httpd_open_os_files gauge
    couchdb_httpd_open_os_files{node_name="couchdb@172.16.238.11"} 23
    couchdb_httpd_open_os_files{node_name="couchdb@172.16.238.12"} 23
    couchdb_httpd_open_os_files{node_name="couchdb@172.16.238.13"} 25
    # HELP couchdb_httpd_request_methods number of HTTP requests by method
    # TYPE couchdb_httpd_request_methods gauge
    couchdb_httpd_request_methods{method="COPY",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_request_methods{method="COPY",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_request_methods{method="COPY",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_request_methods{method="DELETE",node_name="couchdb@172.16.238.11"} 1
    couchdb_httpd_request_methods{method="DELETE",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_request_methods{method="DELETE",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_request_methods{method="GET",node_name="couchdb@172.16.238.11"} 525
    couchdb_httpd_request_methods{method="GET",node_name="couchdb@172.16.238.12"} 152
    couchdb_httpd_request_methods{method="GET",node_name="couchdb@172.16.238.13"} 2
    couchdb_httpd_request_methods{method="HEAD",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_request_methods{method="HEAD",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_request_methods{method="HEAD",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_request_methods{method="POST",node_name="couchdb@172.16.238.11"} 12
    couchdb_httpd_request_methods{method="POST",node_name="couchdb@172.16.238.12"} 3
    couchdb_httpd_request_methods{method="POST",node_name="couchdb@172.16.238.13"} 3
    couchdb_httpd_request_methods{method="PUT",node_name="couchdb@172.16.238.11"} 7
    couchdb_httpd_request_methods{method="PUT",node_name="couchdb@172.16.238.12"} 3
    couchdb_httpd_request_methods{method="PUT",node_name="couchdb@172.16.238.13"} 3
    # HELP couchdb_httpd_request_time length of a request inside CouchDB without MochiWeb
    # TYPE couchdb_httpd_request_time gauge
    couchdb_httpd_request_time{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_request_time{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_request_time{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_requests number of HTTP requests
    # TYPE couchdb_httpd_requests gauge
    couchdb_httpd_requests{node_name="couchdb@172.16.238.11"} 545
    couchdb_httpd_requests{node_name="couchdb@172.16.238.12"} 158
    couchdb_httpd_requests{node_name="couchdb@172.16.238.13"} 8
    # HELP couchdb_httpd_status_codes number of HTTP responses by status code
    # TYPE couchdb_httpd_status_codes gauge
    couchdb_httpd_status_codes{code="200",node_name="couchdb@172.16.238.11"} 495
    couchdb_httpd_status_codes{code="200",node_name="couchdb@172.16.238.12"} 151
    couchdb_httpd_status_codes{code="200",node_name="couchdb@172.16.238.13"} 1
    couchdb_httpd_status_codes{code="201",node_name="couchdb@172.16.238.11"} 11
    couchdb_httpd_status_codes{code="201",node_name="couchdb@172.16.238.12"} 3
    couchdb_httpd_status_codes{code="201",node_name="couchdb@172.16.238.13"} 3
    couchdb_httpd_status_codes{code="202",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="202",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="202",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="301",node_name="couchdb@172.16.238.11"} 2
    couchdb_httpd_status_codes{code="301",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="301",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="304",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="304",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="304",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="400",node_name="couchdb@172.16.238.11"} 1
    couchdb_httpd_status_codes{code="400",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="400",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="401",node_name="couchdb@172.16.238.11"} 3
    couchdb_httpd_status_codes{code="401",node_name="couchdb@172.16.238.12"} 2
    couchdb_httpd_status_codes{code="401",node_name="couchdb@172.16.238.13"} 2
    couchdb_httpd_status_codes{code="403",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="403",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="403",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="404",node_name="couchdb@172.16.238.11"} 2
    couchdb_httpd_status_codes{code="404",node_name="couchdb@172.16.238.12"} 2
    couchdb_httpd_status_codes{code="404",node_name="couchdb@172.16.238.13"} 2
    couchdb_httpd_status_codes{code="405",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="405",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="405",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="409",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="409",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="409",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="412",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="412",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="412",node_name="couchdb@172.16.238.13"} 0
    couchdb_httpd_status_codes{code="500",node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_status_codes{code="500",node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_status_codes{code="500",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_temporary_view_reads number of temporary view reads
    # TYPE couchdb_httpd_temporary_view_reads gauge
    couchdb_httpd_temporary_view_reads{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_temporary_view_reads{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_temporary_view_reads{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_httpd_up Was the last query of CouchDB stats successful.
    # TYPE couchdb_httpd_up gauge
    couchdb_httpd_up 1
    # HELP couchdb_httpd_view_reads number of view reads
    # TYPE couchdb_httpd_view_reads gauge
    couchdb_httpd_view_reads{node_name="couchdb@172.16.238.11"} 0
    couchdb_httpd_view_reads{node_name="couchdb@172.16.238.12"} 0
    couchdb_httpd_view_reads{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_changes_manager_deaths number of failed replicator changes managers
    # TYPE couchdb_replicator_changes_manager_deaths gauge
    couchdb_replicator_changes_manager_deaths{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_changes_manager_deaths{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_changes_manager_deaths{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_changes_queue_deaths number of failed replicator changes work queues
    # TYPE couchdb_replicator_changes_queue_deaths gauge
    couchdb_replicator_changes_queue_deaths{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_changes_queue_deaths{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_changes_queue_deaths{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_changes_read_failures number of failed replicator changes read failures
    # TYPE couchdb_replicator_changes_read_failures gauge
    couchdb_replicator_changes_read_failures{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_changes_read_failures{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_changes_read_failures{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_changes_reader_deaths number of failed replicator changes readers
    # TYPE couchdb_replicator_changes_reader_deaths gauge
    couchdb_replicator_changes_reader_deaths{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_changes_reader_deaths{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_changes_reader_deaths{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_checkpoints replicator checkpoint counters
    # TYPE couchdb_replicator_checkpoints gauge
    couchdb_replicator_checkpoints{metric="failure",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_checkpoints{metric="failure",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_checkpoints{metric="failure",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_checkpoints{metric="success",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_checkpoints{metric="success",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_checkpoints{metric="success",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_cluster_is_stable 1 if cluster is stable, 0 if unstable
    # TYPE couchdb_replicator_cluster_is_stable gauge
    couchdb_replicator_cluster_is_stable{node_name="couchdb@172.16.238.11"} 1
    couchdb_replicator_cluster_is_stable{node_name="couchdb@172.16.238.12"} 1
    couchdb_replicator_cluster_is_stable{node_name="couchdb@172.16.238.13"} 1
    # HELP couchdb_replicator_connections replicator connection metrics shown by type
    # TYPE couchdb_replicator_connections gauge
    couchdb_replicator_connections{metric="acquires",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="acquires",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="acquires",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_connections{metric="closes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="closes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="closes",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_connections{metric="creates",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="creates",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="creates",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_connections{metric="owner_crashes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="owner_crashes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="owner_crashes",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_connections{metric="releases",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="releases",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="releases",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_connections{metric="worker_crashes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_connections{metric="worker_crashes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_connections{metric="worker_crashes",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_db_scans number of times replicator db scans have been started
    # TYPE couchdb_replicator_db_scans gauge
    couchdb_replicator_db_scans{node_name="couchdb@172.16.238.11"} 2
    couchdb_replicator_db_scans{node_name="couchdb@172.16.238.12"} 2
    couchdb_replicator_db_scans{node_name="couchdb@172.16.238.13"} 2
    # HELP couchdb_replicator_docs replicator metrics shown by type
    # TYPE couchdb_replicator_docs gauge
    couchdb_replicator_docs{metric="completed_state_updates",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_docs{metric="completed_state_updates",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_docs{metric="completed_state_updates",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_docs{metric="dbs_changes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_docs{metric="dbs_changes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_docs{metric="dbs_changes",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_docs{metric="dbs_created",node_name="couchdb@172.16.238.11"} 8
    couchdb_replicator_docs{metric="dbs_created",node_name="couchdb@172.16.238.12"} 8
    couchdb_replicator_docs{metric="dbs_created",node_name="couchdb@172.16.238.13"} 8
    couchdb_replicator_docs{metric="dbs_deleted",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_docs{metric="dbs_deleted",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_docs{metric="dbs_deleted",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_docs{metric="dbs_found",node_name="couchdb@172.16.238.11"} 3
    couchdb_replicator_docs{metric="dbs_found",node_name="couchdb@172.16.238.12"} 3
    couchdb_replicator_docs{metric="dbs_found",node_name="couchdb@172.16.238.13"} 11
    couchdb_replicator_docs{metric="failed_state_updates",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_docs{metric="failed_state_updates",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_docs{metric="failed_state_updates",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_failed_starts number of replications that have failed to start
    # TYPE couchdb_replicator_failed_starts gauge
    couchdb_replicator_failed_starts{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_failed_starts{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_failed_starts{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_jobs replicator jobs shown by type
    # TYPE couchdb_replicator_jobs gauge
    couchdb_replicator_jobs{metric="adds",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="adds",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="adds",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="crashed",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="crashed",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="crashed",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="crashes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="crashes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="crashes",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="duplicate_adds",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="duplicate_adds",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="duplicate_adds",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="pending",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="pending",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="pending",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="removes",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="removes",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="removes",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="running",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="running",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="running",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="starts",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="starts",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="starts",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="stops",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="stops",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="stops",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_jobs{metric="total",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_jobs{metric="total",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_jobs{metric="total",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_requests number of HTTP requests made by the replicator
    # TYPE couchdb_replicator_requests gauge
    couchdb_replicator_requests{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_requests{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_requests{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_responses number of HTTP responses by state
    # TYPE couchdb_replicator_responses gauge
    couchdb_replicator_responses{metric="failure",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_responses{metric="failure",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_responses{metric="failure",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_responses{metric="success",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_responses{metric="success",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_responses{metric="success",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_stream_responses number of streaming HTTP responses by state
    # TYPE couchdb_replicator_stream_responses gauge
    couchdb_replicator_stream_responses{metric="failure",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_stream_responses{metric="failure",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_stream_responses{metric="failure",node_name="couchdb@172.16.238.13"} 0
    couchdb_replicator_stream_responses{metric="success",node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_stream_responses{metric="success",node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_stream_responses{metric="success",node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_worker_deaths number of failed replicator workers
    # TYPE couchdb_replicator_worker_deaths gauge
    couchdb_replicator_worker_deaths{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_worker_deaths{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_worker_deaths{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_replicator_workers_started number of replicator workers started
    # TYPE couchdb_replicator_workers_started gauge
    couchdb_replicator_workers_started{node_name="couchdb@172.16.238.11"} 0
    couchdb_replicator_workers_started{node_name="couchdb@172.16.238.12"} 0
    couchdb_replicator_workers_started{node_name="couchdb@172.16.238.13"} 0
    # HELP couchdb_server_couch_log number of messages logged by log level
    # TYPE couchdb_server_couch_log gauge
    couchdb_server_couch_log{level="alert",node_name="couchdb@172.16.238.11"} 0
    couchdb_server_couch_log{level="alert",node_name="couchdb@172.16.238.12"} 0
    couchdb_server_couch_log{level="alert",node_name="couchdb@172.16.238.13"} 0
    couchdb_server_couch_log{level="critical",node_name="couchdb@172.16.238.11"} 0
    couchdb_server_couch_log{level="critical",node_name="couchdb@172.16.238.12"} 0
    couchdb_server_couch_log{level="critical",node_name="couchdb@172.16.238.13"} 0
    couchdb_server_couch_log{level="debug",node_name="couchdb@172.16.238.11"} 0
    couchdb_server_couch_log{level="debug",node_name="couchdb@172.16.238.12"} 0
    couchdb_server_couch_log{level="debug",node_name="couchdb@172.16.238.13"} 0
    couchdb_server_couch_log{level="emergency",node_name="couchdb@172.16.238.11"} 0
    couchdb_server_couch_log{level="emergency",node_name="couchdb@172.16.238.12"} 0
    couchdb_server_couch_log{level="emergency",node_name="couchdb@172.16.238.13"} 0
    couchdb_server_couch_log{level="error",node_name="couchdb@172.16.238.11"} 7
    couchdb_server_couch_log{level="error",node_name="couchdb@172.16.238.12"} 2
    couchdb_server_couch_log{level="error",node_name="couchdb@172.16.238.13"} 2
    couchdb_server_couch_log{level="info",node_name="couchdb@172.16.238.11"} 8
    couchdb_server_couch_log{level="info",node_name="couchdb@172.16.238.12"} 8
    couchdb_server_couch_log{level="info",node_name="couchdb@172.16.238.13"} 8
    couchdb_server_couch_log{level="notice",node_name="couchdb@172.16.238.11"} 589
    couchdb_server_couch_log{level="notice",node_name="couchdb@172.16.238.12"} 192
    couchdb_server_couch_log{level="notice",node_name="couchdb@172.16.238.13"} 42
    couchdb_server_couch_log{level="warning",node_name="couchdb@172.16.238.11"} 4
    couchdb_server_couch_log{level="warning",node_name="couchdb@172.16.238.12"} 4
    couchdb_server_couch_log{level="warning",node_name="couchdb@172.16.238.13"} 20
    # HELP couchdb_server_node_info General info about a node.
    # TYPE couchdb_server_node_info gauge
    couchdb_server_node_info{node_name="couchdb@172.16.238.11",vendor_name="The Apache Software Foundation",version="2.3.0"} 1
    couchdb_server_node_info{node_name="couchdb@172.16.238.12",vendor_name="The Apache Software Foundation",version="2.3.0"} 1
    couchdb_server_node_info{node_name="couchdb@172.16.238.13",vendor_name="The Apache Software Foundation",version="2.3.0"} 1
    ... (more go internal stats)
