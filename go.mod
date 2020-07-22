module github.com/gesellix/couchdb-prometheus-exporter

go 1.13

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/gesellix/couchdb-cluster-config v0.0.0-20200218123558-43a3249a3c3c
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-version v1.2.1
	github.com/okeuday/erlang_go v2.0.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/sys v0.0.0-20200722175500-76b94024e4b6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/klog v1.0.0
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
