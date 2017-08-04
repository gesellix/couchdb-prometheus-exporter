# Using Docker Stack

Docker Stack is similar to Docker Compose, but uses Swarm services instead of plain containers.
Swarm, and thereby Stack, also allows you to configure secrets and arbitrary configs.
The example stacks showcase those features.

You need to ensure that the Docker Swarm mode is enabled before deploying a stack:

    docker swarm init

The admin credentials need to be made available with a simple text file,
which will then be used by the example stacks to create Docker Swarm secrets:

    cat << EOB > config.properties 
    couchdb.username=root
    couchdb.password=a-secret
    EOB

## Run a single CouchDB instance along with the CouchDB Prometheus exporter

The following command reads the stack YAML and creates both a CouchDB node and a Prometheus exporter node.
Missing Docker images will be pulled from the Docker Hub. A one-shot `couchdb-config` service
will use the previously created properties file to configure the admin credentials. 

    docker stack deploy -c docker-stack-single.yml couchdb-single


## Run a clustered CouchDB setup and the CouchDB Prometheus exporter

Similar to the single node stack the following command deploys a multi-node CouchDB dev cluster including
a HAProxy as load balancer:

    docker stack deploy -c docker-stack-cluster.yml couchdb-cluster
