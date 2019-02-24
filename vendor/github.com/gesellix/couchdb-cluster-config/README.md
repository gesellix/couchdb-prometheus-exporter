# couchdb-cluster-config

Util to configure a CouchDB 2.x cluster with several nodes.

## Motivation

Maybe you've followed the official [cluster setup docs](http://docs.couchdb.org/en/stable/cluster/setup.html)
when trying to setup a CouchDB cluster.

Maybe, like me, you didn't _really_ follow every tiny bit.
One example where one has to pay attention is the `vm.args` file
to configure a common Erlang cookie, or to set the node
names to something different than `localhost`.

Since I also needed more automation via api to setup an integration test environment
I played with same hacky [shell scripts](https://github.com/gesellix/couchdb-prometheus-exporter/commit/73fae7bc37194a0c8e63107fb16d7993d9cfef25),
but a nice implementation in Golang promises better portability and maintainability.

So here it is, a little utility to configure a CouchDB cluster.

## Usage

The tool currently only aims at the single task of initializing a CouchDB 2.x cluster,
but is designed for further extension by leveraging the [urfave/cli](https://github.com/urfave/cli)
package.

### Download/Installation

Like with most Golang based tools, you only need to have Golang installed locally.

    go get github.com/gesellix/couchdb-cluster-config
    couchdb-cluster-config --help

If you don't want to install a complete Golang package, you can also use the ready to run
Docker image.

    docker pull gesellix/couchdb-cluster-config
    docker run --rm gesellix/couchdb-cluster-config --help

There's no configuration necessary, everything works using command line parameters.

### Perform a fresh cluster setup

The couchdb-cluster-config tool expects to run in the same network like the CouchDB nodes.
You can run a setup like this:

    docker run --rm \
               --network couchdb-cluster \
               gesellix/couchdb-cluster-config \
               setup \
               --nodes 172.16.238.11 \
               --nodes 172.16.238.12 \
               --nodes 172.16.238.13 \
               --username root \
               --password a-secret

There are three nodes listed with their ip addresses, along with the admin credentials.
The tool should work with freshly started nodes, so they usually don't know about an
admin user and not even the core databases `_users` and `_replicator`. For every listed node
the couchdb-cluster-config will ensure that the admin user and core databases exist.
Only then the cluster setup is performed by creating a cluster of all nodes.

### Some more details

To ease running the cluster setup in an automated environment, it waits for every node
to be available on port 5984. This allows to run the tool very early during startup
of a bunch of nodes.

You'll find a `docker-compose.yml` in the root directory of this repository. It allows
to setup an example cluster and also shows how to set node names and the Erlang cookie
before CouchDB startup.

## Contributing

If you found a bug, want to improve something or need a new feature, then
please file an issue or submit a pull request!
You can also contact me at Twitter [@gesellix](https://twitter.com/gesellix).

Never forget: Relax.
