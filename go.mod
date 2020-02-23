module github.com/gesellix/couchdb-prometheus-exporter

go 1.13

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/gesellix/couchdb-cluster-config v0.0.0-20200218123558-43a3249a3c3c
	github.com/golang/protobuf v1.3.3
	github.com/hashicorp/go-version v1.2.1-0.20190424083514-192140e6f3e6
	github.com/okeuday/erlang_go v1.8.0
	github.com/prometheus/client_golang v1.4.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/procfs v0.0.10 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.2
	golang.org/x/sys v0.0.0-20200219091948-cb0a6d8edb6c // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/klog v1.0.0
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
