# Creating a local dev cluster using Docker Compose

Prepare a volume to manage the CouchDB data. The network helps in case of debugging.
These steps only need to be performed once.  

    docker volume create --name couchdb-test
    docker network create --attachable couchdb-test

Prepare the container env.

    cat << EOB > .couchdb-env 
    COUCHDB_USERNAME=root
    COUCHDB_PASSWORD=a-secret
    COUCHDB_USER=root
    COUCHDB_PASSWORD=a-secret
    EOB

Run a CouchDB instance along with the CouchDB exporter.

    docker-compose up

Configure the cluster, including the root (admin) user.

    curl -X POST "http://localhost:5984/_cluster_setup" \
         -H "Content-Type: application/json" \
         -d '{"action":"enable_cluster", "username":"root", "password":"a-secret", "bind_address":"0.0.0.0", "port":5984, "node_count": 1}'
    curl -X POST "http://localhost:5984/_session" \
         -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
         -d 'name=root&password=a-secret' \
         --cookie-jar couchdb-cookies.txt
    curl -X POST "http://localhost:5984/_cluster_setup" \
         -H "Content-Type: application/json" \
         -d '{"action":"finish_cluster"}' \
         --cookie couchdb-cookies.txt
