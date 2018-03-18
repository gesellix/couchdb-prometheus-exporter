package cluster_config

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ClusterSetup struct {
	Action      string `json:"action"`
	RemoteNode  string `json:"remote_node,omitempty"`
	Host        string `json:"host,omitempty"`
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

func SetupClusterNodes(ipAddresses []string, insecure bool) error {
	hosts := make([]string, len(ipAddresses))
	for i, ip := range ipAddresses {
		hosts[i] = fmt.Sprintf("%s:5984", ip)
	}
	err := AwaitNodes(hosts)
	if err != nil {
		return err
	}

	// TODO extract node setup into dedicated functions
	for _, ip := range ipAddresses {
		client := NewCouchdbClient(fmt.Sprintf("http://%s:5984", ip), BasicAuth{}, []string{}, insecure)

		databaseNames := []string{"_users", "_replicator"}
		for _, dbName := range databaseNames {
			_, err = client.Request("PUT", fmt.Sprintf("%s/%s", client.BaseUri, dbName), nil)
			if err != nil {
				return err
			}
		}

		_, err = client.Request("PUT", fmt.Sprintf("%s/_node/couchdb@%s/_config/admins/%s", client.BaseUri, ip, adminUsername), strings.NewReader(fmt.Sprintf("\"%s\"", adminPassword)))
		if err != nil {
			return err
		}
	}

	setupNodeIp := ipAddresses[:1]
	otherNodeIps := ipAddresses[1:]
	nodeCount := len(ipAddresses)

	client := NewCouchdbClient(fmt.Sprintf("http://%s:5984", setupNodeIp), BasicAuth{Username: adminUsername, Password: adminPassword}, []string{}, insecure)

	body, err := json.Marshal(ClusterSetup{
		Action:      "enable_cluster",
		Username:    adminUsername,
		Password:    adminPassword,
		BindAddress: "0.0.0.0",
		NodeCount:   nodeCount})
	if err != nil {
		return err
	}
	client.Request("POST", fmt.Sprintf("http://%s:5984/_cluster_setup", setupNodeIp), strings.NewReader(string(body)))

	for _, ip := range otherNodeIps {
		body, err = json.Marshal(ClusterSetup{
			Action:      "enable_cluster",
			RemoteNode:  ip,
			Port:        "5984",
			Username:    adminUsername,
			Password:    adminPassword,
			BindAddress: "0.0.0.0",
			NodeCount:   nodeCount})
		if err != nil {
			return err
		}
		client.Request("POST", fmt.Sprintf("http://%s:5984/_cluster_setup", setupNodeIp), strings.NewReader(string(body)))

		body, err = json.Marshal(ClusterSetup{
			Action:   "add_node",
			Host:     ip,
			Port:     "5984",
			Username: adminUsername,
			Password: adminPassword})
		if err != nil {
			return err
		}
		client.Request("POST", fmt.Sprintf("http://%s:5984/_cluster_setup", setupNodeIp), strings.NewReader(string(body)))
	}

	body, err = json.Marshal(ClusterSetup{
		Action: "finish_cluster"})
	if err != nil {
		return err
	}
	client.Request("POST", fmt.Sprintf("http://%s:5984/_cluster_setup", setupNodeIp), strings.NewReader(string(body)))

	return nil
}
