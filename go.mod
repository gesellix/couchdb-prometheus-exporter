module github.com/gesellix/couchdb-prometheus-exporter

go 1.13

require (
	github.com/gesellix/couchdb-cluster-config v0.0.0-20200218123558-43a3249a3c3c
	github.com/golang/protobuf v1.3.3-0.20190805180045-4c88cc3f1a34
	github.com/hashicorp/go-version v1.2.1-0.20190424083514-192140e6f3e6
	github.com/okeuday/erlang_go v1.7.5
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.0
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	k8s.io/klog v0.3.3
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
