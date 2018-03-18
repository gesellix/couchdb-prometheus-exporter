package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gesellix/couchdb-cluster-config/pkg"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "nodes",
			Usage: "list of node ip addresses to participate in the CouchDB cluster",
		},
		cli.BoolTFlag{
			Name:  "insecure",
			Usage: "Ignore server certificate if using https",
		},
	}

	app.Action = func(c *cli.Context) error {
		nodes := c.StringSlice("nodes")
		if len(nodes) == 0 {
			return fmt.Errorf("please pass a list of node ip addresses")
		}

		fmt.Printf("Going to setup the following nodes as cluster\n%v\n", nodes)
		return cluster_config.SetupClusterNodes(nodes, c.Bool("insecure"))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
