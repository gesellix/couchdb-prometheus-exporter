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
	app.Name = "CouchDB Cluster Config"
	app.Usage = ""
	app.Description = "Setup a CouchDB 2.x cluster"
	app.Version = ""
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "perform a cluster setup",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "nodes",
					Usage: "list of node ip addresses to participate in the CouchDB cluster",
				},
				cli.StringFlag{
					Name:  "username",
					Usage: "admin username - admin will be created, if missing",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "admin password",
				},
			},
			Action: func(c *cli.Context) error {
				nodes := c.StringSlice("nodes")
				if len(nodes) == 0 {
					return fmt.Errorf("please pass a list of node ip addresses")
				}

				if c.String("username") == "" {
					return fmt.Errorf("please provide an admin username")
				}
				if c.String("password") == "" {
					return fmt.Errorf("please provide an admin password")
				}

				fmt.Printf("Going to setup the following nodes as cluster\n%v\n", nodes)
				return cluster_config.SetupClusterNodes(
					nodes,
					cluster_config.BasicAuth{
						Username: c.String("username"),
						Password: c.String("password")},
					c.Bool("insecure"))
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:  "insecure",
			Usage: "ignore server certificate if using https",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
