package cluster_setup

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"github.com/gesellix/couchdb-prometheus-exporter/lib"
)

type ClusterSetup struct {
	Action      string `json:"action"`
	RemoteNode  string `json:"remote_node,omitempty"`
	Port        string `json:"port,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	BindAddress string `json:"bind_address,omitempty"`
	NodeCount   int    `json:"node_count,omitempty"`
}

const (
	adminUsername = "root"
	adminPassword = "a-secret"
)

func SetupClusterNodes() {
	urls := []string{
		"172.16.238.11:5984",
		"172.16.238.12:5984",
		"172.16.238.13:5984",
	}

	err := AwaitNodes(urls)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, u := range urls {
		ip, err := ipAddress(u)
		if err != nil {
			fmt.Println(err)
			return
		}
		client := lib.NewCouchdbClient(fmt.Sprintf("http://%s", u), lib.BasicAuth{}, []string{})

		coreDatabases := []string{"_users", "_replicator"}
		for _, db := range coreDatabases {
			_, err = client.Request("PUT", fmt.Sprintf("%s/%s", client.BaseUri, db), nil)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		_, err = client.Request("PUT", fmt.Sprintf("%s/_node/couchdb@%s/_config/admins/%s", client.BaseUri, ip, adminUsername), strings.NewReader(fmt.Sprintf("\"%s\"", adminPassword)))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	setupNodeUrl := urls[:1]
	otherNodeUrls := urls[1:]
	nodeCount := len(urls)

	client := lib.NewCouchdbClient(fmt.Sprintf("http://%s", setupNodeUrl), lib.BasicAuth{Username: adminUsername, Password: adminPassword}, []string{})

	body, err := json.Marshal(ClusterSetup{
		Action:      "enable_cluster",
		Username:    adminUsername,
		Password:    adminPassword,
		BindAddress: "0.0.0.0",
		NodeCount:   nodeCount})
	if err != nil {
		return
	}
	client.Request("PUT", fmt.Sprintf("http://%s/_cluster_setup", setupNodeUrl), strings.NewReader(string(body)))

	for _, u := range otherNodeUrls {
		body, err = json.Marshal(ClusterSetup{
			Action:      "enable_cluster",
			RemoteNode:  u,
			Port:        "5984",
			Username:    adminUsername,
			Password:    adminPassword,
			BindAddress: "0.0.0.0",
			NodeCount:   nodeCount})
		if err != nil {
			return
		}
		client.Request("PUT", fmt.Sprintf("http://%s/_cluster_setup", setupNodeUrl), strings.NewReader(string(body)))

		body, err = json.Marshal(ClusterSetup{
			Action:     "add_node",
			RemoteNode: u,
			Port:       "5984",
			Username:   adminUsername,
			Password:   adminPassword})
		if err != nil {
			return
		}
		client.Request("PUT", fmt.Sprintf("http://%s/_cluster_setup", setupNodeUrl), strings.NewReader(string(body)))
	}

	body, err = json.Marshal(ClusterSetup{
		Action: "finish_cluster"})
	if err != nil {
		return
	}
	client.Request("PUT", fmt.Sprintf("http://%s/_cluster_setup", setupNodeUrl), strings.NewReader(string(body)))
}

func ipAddress(u string) (string, error) {
	parsed, err := url.Parse(fmt.Sprintf("http://%s", u))
	if err != nil {
		return "", err
	}
	return parsed.Hostname(), nil
}
