---
title: "Use containers for development"
keywords: get started, go, golang, local, development
description: Learn how to develop your application locally.
redirect_from:
- /get-started/golang/develop/
---

{% include_relative nav.html selected="3" %}

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Run your image as a container](run-containers.md).

## Introduction

In this module, we’ll walk through connecting a database to the application we built in the previous modules. The database will run in its own container with Docker providing a link between that container and the container in which our Go application is running. We’ll then learn how to use Docker Compose to manage such multi-container local development environments effectively.

## Local database and containers

First, we’ll take a look at running a database in a container and how we use volumes and networking to persist our data and allow our application to talk with the database. Then we’ll pull everything together into a compose file which will allow us to setup and run a local development environment with one command. Finally, we’ll also take a look at connecting to a running container to run some commands _inside_ a container.

The database engine we are going to use is called [CockroachDB](https://www.cockroachlabs.com/product/). It's a modern, Cloud-native, distributed SQL database.

Instead of downloading CockroachDB, installing, configuring and then running it using the operating system's native package manager, we can use the Docker Official Image for CockroachDB and run it in a container.

### Storage 

Before we run CockroachDB in a container, we want to create a volume that Docker can manage to store our persistent data and configuration. Let’s create a managed volume now.

```shell
$ docker volume create roach1
```

You can view a list of all managed volumes in your Docker instance with the following command.

```shell
$ docker volume list
```

### Networking

Now we’ll create a network that our application and database will use to talk with each other. The network is called a user-defined bridge network and provides us with a DNS lookup service which we can use when creating our database connection string.

```shell
$ docker network create -d bridge godockernet
```

As it was the case with the managed volumes, there is a command to list all networks set up in your Docker instance.

```shell
$ docker network list
```

Note, that although we've named our managed volume `roach` and the name `godockernet` was given to a network which is going to link the database container to the application container, both names are arbitrary and we could have named them however we wanted. It's useful though to choose a name which is indicative of the intended purpose.

### Start the database engine

Now we can run CockroachDB in a container and attach to the volume and network we had just created. Docker will pull the image from Docker Hub and run it for you locally.

```shell
$ docker run -d \
  --name roach1 \
  --hostname roach1 \
  --network godockernet \
  -p 26257:26257 -p 8080:8080 \
  -v roach1:/cockroach/cockroach-data \
  cockroachdb/cockroach:v20.2.5 start-single-node \
  --insecure
```

### Configure the database engine

Now that the database engine is live, we can start the built-in SQL shell in the _same_ container where it is running.

```shell
$ docker exec -it roach1 ./cockroach sql --insecure
```

1. In the SQL shell, create the database that our application will use:

   ```sql
   CREATE DATABASE godocker;
   ```

2. Create a SQL user for our application:

   ```sql
   CREATE USER <username>;
   ```
   
   Take a note of the username. We will use it in our application code later.

3. Give the user the necessary permissions:
   
   ```sql
   GRANT ALL ON DATABASE godocker TO <username>;
   ```

You can type `quit` to exit the shell. 

An example interaction with SQL shell is quoted below.

```sql
ofr@NEON:~$ docker exec -it roach1 ./cockroach sql --insecure
#
# Welcome to the CockroachDB SQL shell.
# All statements must be terminated by a semicolon.
# To exit, type: \q.
#
# Server version: CockroachDB CCL v20.2.5 (x86_64-unknown-linux-gnu, built 2021/02/16 12:52:58, go1.13.14) (same version as client)
# Cluster ID: 42125797-c8d5-4542-9281-18473783abdc
#
# Enter \? for a brief introduction.
#
root@:26257/defaultdb> CREATE DATABASE godocker;
CREATE DATABASE

Time: 83ms total (execution 83ms / network 0ms)

root@:26257/defaultdb> create user wumpus;
CREATE ROLE

Time: 39ms total (execution 39ms / network 0ms)

root@:26257/defaultdb> GRANT ALL ON DATABASE godocker TO wumpus;
GRANT

Time: 161ms total (execution 66ms / network 95ms)

root@:26257/defaultdb> quit
ofr@NEON:~$
```

### Update the application

Now that we have started and configured the database engine, let’s update `main.go` to use it and not an in-memory data store.

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"context"
	"database/sql"
	"log"

	"github.com/cockroachdb/cockroach-go/crdb"
	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Entry is a single key-value pair
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {

	// Read the connection parameters from environment variables,
	// following the Twelve-Factor App philosophy: https://12factor.net/
	username := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	hostname := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	database := os.Getenv("PGDATABASE")

	// Build the connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, hostname, port, database)

	// Initialise the store
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()

	// Create the table to store key-value pairs
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS store (key STRING PRIMARY KEY, value STRING)"); err != nil {
		log.Fatal(err)
	}

	// Initialise the router
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure the routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})
	e.GET("/list", func(c echo.Context) error {
		return listEntries(db, c)
	})
	e.POST("/add", func(c echo.Context) error {
		return addEntry(db, c)
	})

	// Blast off!
	port = os.Getenv("GODOCKER_PORT")
	if port == "" {
		port = "8000"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// listEntries returns a full copy of the store in JSON format to the client.
func listEntries(db *sql.DB, c echo.Context) error {

	rows, err := db.Query("SELECT key, value FROM store")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Temporary storage for the values read from the DB
	store := map[string]string{}

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			log.Fatal(err)
		}
		store[key] = value
	}

	return c.JSON(http.StatusOK, store)
}

// addEntry adds a new entry to the key-value store
// and returns the value in JSON format to the client.
func addEntry(db *sql.DB, c echo.Context) error {

	e := new(Entry)
	if err := c.Bind(e); err != nil {
		return err
	}

	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			"INSERT INTO store (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = excluded.value",
			e.Key, e.Value); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, e)
}
```

Depending on the IDE you use, you might need to `go get` the modules imported in `main.go`:

```shell
$ go get github.com/cockroachdb/cockroach-go/crdb
$ go get github.com/lib/pq
```

Since our application source code has been updated, we now need to rebuild our Docker image to reflect these changes.

### Build the application

We can build an updated image with the following command.

```shell
$ docker build --tag go-docker .
```

### Run the application

Now, let’s run our container. But this time we’ll need to set quite a few environment variables so that our application could make the connection string necessary to access the database. We’ll do this right in the `docker run` command.

Note, that in the following command the value of `PGUSER` must be set to the same value that you had picked before for the _Configure the database engine_ section. In our example, we use the value of `wumpus`.

```shell
$ docker run \
  -it --rm -d \
  --network godockernet \
  --name rest-server \
  -p 8000:8000 \
  -e PGUSER=wumpus -e PGPASSWORD=passwd9 \
  -e PGHOST=roach1 -e PGPORT=26257 -e PGDATABASE=godocker \
  go-docker
```

### Testing the application

Let’s test that our application is connected to the database and is able to save key-value pairs. We'll reuse the test input from the [Build your Go image](build-images.md) section:

```shell
$ curl --request POST \
  --url http://localhost:8000/add
  --header 'content-type: application/json'
  --data '{"key": "name", "value": "Docker"}'
```

You should receive the following json back from our service.

```json
{"name":"Docker"}
```

## Use Compose to develop locally

At this point you might be wondering, if there is a way to avoid having to deal with long lists of argument to the Docker command. The toy example we considered requires five environment variables to define the connection to the database. A real application might need many more. Then there is also a question of dependencies &ndash; ideally, we would like to make sure that the database is started _before_ our application is run. And spinning up the database instance may require another Docker command with many options. But there is a better way to orchestrate these deployments for local development purposes.

In this section, we’ll create a Compose file to start our `go-docker` application and CockroachDB database engine with one command.

### Docker Compose configuration

In our application's directory, create a new text file named `docker-compose.dev.yml` with the following content.

```yaml
version: '3.8'

services:

  godocker:
    depends_on: 
      - cockroachdb
    build:
      context: .
    container_name: rest-server
    hostname: rest-server
    networks:
      - godockernet
    ports:
      - 8000:8000
    environment:
      - PGUSER=${PGUSER:-wumpus}
      - PGPASSWORD=${PGPASSWORD}
      - PGHOST=${PGHOST:-roach1}
      - PGPORT=${PGPORT:-26257}
      - PGDATABASE=${PGDATABASE-godocker}

  cockroachdb:
    image: cockroachdb/cockroach:v20.2.5
    container_name: roach1
    hostname: roach1
    networks:
      - godockernet
    ports:
      - 26257:26257
      - 8080:8080
    volumes:
      - roach1:/cockroach/cockroach-data
    command: start-single-node --insecure

volumes:
  roach1:

networks:
  godockernet:
    driver: bridge
```

This Compose file is super convenient as we do not have to type all the parameters to pass to the `docker run` command. We can declaratively do that in the Compose file. [Docker Compose Documentation](../../compose/index.md) is quite extensive and includes a full reference for the Docker Compose file format.

### Variable substitution in Docker Compose

One of the really cool features of Docker Compose is [variable substitution](https://docs.docker.com/compose/compose-file/compose-file-v3/#variable-substitution). You can see the example in our Docker Compose file, `godocker/environment` section. By means of example:

* `PGUSER=${PGUSER:-wumpus}` means that inside the container, the environment variable `PGUSER` shall be set to the same value as it has on the host machine where Docker Compose is run. If there is no environment variable with such name on the host machine, the variable inside the container gets the default value of `wumpus`. This is what `:-` separator is for.

Other ways of dealing with undefined or empty values exist, check out the above link, if this sounds interesting.

### Validating Docker Compose configuration

Before you apply changes made to a Docker Compose configuration file, there is an opportunity to validate the content of the configuration file with the following command.

```shell
$ docker-compose -f docker-compose.dev.yml config
```

When this command is run, Docker Compose would read the file `docker-compose.dev.yml`, parse it into a data structure in memory, validate where possible, and print back the _reconstruction_ of that configuration file from its internal representation. If this is not possible due to errors, it would print an error message instead.

### Build and run the application using Docker Compose
Let’s start our application and confirm that it is running properly.

```shell
$ docker-compose -f docker-compose.dev.yml up --build
```

We pass the `--build` flag so Docker will compile our image and then starts it.

Since our set-up is now run by Docker Compose, it has assigned it a "project name", so we got a new volume for our CockroachDB instance. This means that our REST server application would fail to connect to the database, because the database does not exist in this new volume. My terminal displays the following output.

```
...
roach1         | CockroachDB node starting at 2021-02-19 06:48:17.3620647 +0000 UTC (took 0.7s)
roach1         | build:               CCL v20.2.5 @ 2021/02/16 12:52:58 (go1.13.14)
roach1         | webui:               http://roach1:8080
roach1         | sql:                 postgresql://root@roach1:26257?sslmode=disable
roach1         | RPC client flags:    /cockroach/cockroach <client cmd> --host=roach1:26257 --insecure
roach1         | logs:                /cockroach/cockroach-data/logs
roach1         | temp dir:            /cockroach/cockroach-data/cockroach-temp062367352
roach1         | external I/O path:   /cockroach/cockroach-data/extern
roach1         | store[0]:            path=/cockroach/cockroach-data
roach1         | storage engine:      pebble
roach1         | status:              restarted pre-existing node
roach1         | clusterID:           e7e2aece-307d-4fc0-b04d-d7f6ccf399cb
roach1         | nodeID:              1
rest-server    | 2021/02/19 06:48:18 pq: password authentication failed for user wumpus
rest-server    | exit status 1
rest-server exited with code 1
```

This is not a big deal. All we have to do is to connect to CockroachDB instance and run the three SQL commands to create the database and the user, as described above in the _Configure the database engine_ section.

It would have been possible to connect the volume that we had previously used, but for the purposes of our example it's more trouble than it's worth. On a side note, it is highly recommended to learn about "pruning" unnecessary Docker resources, such as stopped containers, old volumes, and unused network configurations. 

TODO add link to pruning tutorial?

### Testing the application

Now let’s test our API endpoint. Run the following curl command.

```shell
$ curl --request GET --url http://localhost:8080/
```

You should receive the following response:

```json
{"Status":"OK"}
```

## TODO

TODO Two-stage build (no Go compiler present in the final image, only the app binary)

## Things we didn't cover

There are some tangential, yet interesting points that were purposefully not covered in this chapter. For the more adventurous reader, this section offers some pointers for further study.

### Persistent storage

A _managed volume_ isn't the only way to provide your container with persistent storage. It is highly recommended to get acquainted with available storage options and their use cases, covered in the following part of Docker documentation: [Manage data in Docker](../storage/index.md).

### CockroachDB clusters

We run a single instance of CockroachDB, which was enough for our demonstration. But it is possible to run a CockroachDB _cluster_, which is made of multiple instances of CockroachDB, each instance running in its own container. Since CockroachDB engine is distributed by design, it would have taken us surprisingly little change to our procedure to run a cluster with multiple nodes. For more information, please check out [CockroachDB Documentation](https://www.cockroachlabs.com/docs/v20.2/start-a-local-cluster-in-docker-mac.html).

Such distributed set-up offers interesting possibilities, such as applying _Chaos Engineering_ techniques to simulate parts of the cluster failing and evaluating our application's ability to cope with such failures. 

### Other databases

Since we did not run a cluster of CockroachDB instances, you might be wondering whether we could have used a non-distributed database engine. The answer is 'yes', and if we were to pick a more traditional SQL database, such as [PostgreSQL](https://www.postgresql.org/), the process described in this chapter would have been very similar.

## Next steps

In this module, we set up a containerised development environment with our application and the database engine running in different containers. We also wrote a Docker Compose file which links the two containers together and provides for easy starting up and tearing down of the development environment.

In the next module, we’ll take a look at how to run unit tests in Docker. See:

[Run your tests](run-tests.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs ](https://github.com/docker/docker.github.io/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
