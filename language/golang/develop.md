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

In this module, we’ll walk through connecting a database to the _extended version_ of the example application from the previous modules. The database will run in its own container with Docker providing the link between that container and the container in which our example application is running. We’ll then learn how to use Docker Compose to manage such multi-container local development environments effectively.

## Local database and containers

First, we’ll take a look at running a database engine in a container and how we use volumes and networking to persist our data and allow our application to talk with the database. Then we’ll pull everything together into a Docker _compose file_ which will allow us to setup and run a local development environment with a single command. Finally, we’ll also take a look at connecting to a running container to run some commands _inside_ a container.

The database engine we are going to use is called [CockroachDB](https://www.cockroachlabs.com/product/). It's a modern, Cloud-native, distributed SQL database.

Instead of downloading CockroachDB, installing, configuring, and then running it using the operating system's native package manager, we can use the [Docker image for CockroachDB](https://hub.docker.com/r/cockroachdb/cockroach) and run it in a container.

### Storage 

Before we run CockroachDB in a container, we are going to to create a _volume_ that Docker can manage to store our persistent data and configuration. Let’s create a managed volume now:

```shell
$ docker volume create roach1
roach1
```

You can view a list of all managed volumes in your Docker instance with the following command:

```shell
$ docker volume list
DRIVER    VOLUME NAME
local     roach1
```

### Networking

Now we’ll create a network that our application and database will use to talk with each other. There are different types of network configurations possible, we are going to use what is called a user-defined _bridge network_. It would provide us with a DNS lookup service which we can use when creating our database connection string.

```shell
$ docker network create -d bridge mynet
51344edd6430b5acd121822cacc99f8bc39be63dd125a3b3cd517b6485ab7709
```

As it was the case with the managed volumes, there is a command to list all networks set up in your Docker instance.

```shell
$ docker network list
NETWORK ID     NAME          DRIVER    SCOPE
0ac2b1819fa4   bridge        bridge    local
51344edd6430   mynet         bridge    local
daed20bbecce   host          host      local
6aee44f40a39   none          null      local
```

Note, that although we've named our managed volume `roach` and the name `godockernet` was given to a network which is going to link the database container to the application container, both names are arbitrary and we could have named them however we wanted. It's useful though to choose a name which is indicative of the intended purpose.

### Start the database engine

Now we can run CockroachDB in a container and attach to the volume and network we had just created. Docker will pull the image from Docker Hub and run it for you locally.

```shell
$ docker run -d \
  --name roach1 \
  --hostname roach1 \
  --network mynet \
  -p 26257:26257 -p 8080:8080 \
  -v roach1:/cockroach/cockroach-data \
  cockroachdb/cockroach:v20.2.5 start-single-node \
  --insecure

# ... output omitted ...
```

### Configure the database engine

Now that the database engine is live, there is some housekeeping to do before we can begin using it. Fortunately, it's not a lot. We must:

1. Create a blank database.
2. Register a new user account with the database engine.
3. Grant that new user access rights to the database.

We can do that with the help of CockroachDB built-in SQL shell. To start the SQL shell in the _same_ container where the database engine is running, type:

```shell
$ docker exec -it roach1 ./cockroach sql --insecure
```

1. In the SQL shell, create the database that our example application is going to use:

   ```sql
   CREATE DATABASE mydb;
   ```

2. Register a new SQL user account with the database engine. We pick the username `totoro`.

   ```sql
   CREATE USER totoro;
   ```

3. Give the new user the necessary permissions:
   
   ```sql
   GRANT ALL ON DATABASE mydb TO totoro;
   ```

4. Type `quit` to exit the shell. 

An example interaction with the SQL shell is presented below.

{% raw %}
```
ofr@hki:~/docker-gs-ping-roach$ sudo docker exec -it roach1 ./cockroach sql --insecure
#
# Welcome to the CockroachDB SQL shell.
# All statements must be terminated by a semicolon.
# To exit, type: \q.
#
# Server version: CockroachDB CCL v20.2.5 (x86_64-unknown-linux-gnu, built 2021/02/16 12:52:58, go1.13.14) (same version as client)
# Cluster ID: 6819a1fc-9da9-48d9-a87f-2b916f0dd0f7
#
# Enter \? for a brief introduction.
#
root@:26257/defaultdb> CREATE DATABASE mydb;
CREATE DATABASE

Time: 22ms total (execution 22ms / network 0ms)

root@:26257/defaultdb> CREATE USER totoro;
CREATE ROLE

Time: 15ms total (execution 15ms / network 0ms)

root@:26257/defaultdb> GRANT ALL ON DATABASE mydb TO totoro;
GRANT

Time: 57ms total (execution 21ms / network 36ms)

root@:26257/defaultdb> quit
ofr@hki:~/docker-gs-ping-roach$
```
{% endraw %}

### Update the application

Now that we have started and configured the database engine, we can swith our attention to the application. 

The example application for this module is an extended version of `docker-gs-ping` application we've used in the previous modules. You have two options:

* You can update your local copy of `docker-gs-ping` to match the new extended version presented in this chapter; or
* You can clone the [olliefr/docker-gs-ping-roach](https://github.com/olliefr/docker-gs-ping-roach) repo. This latter approach is recommended.

To checkout the example application, run:

```shell
$ git clone https://github.com/olliefr/docker-gs-ping-roach.git
# ... output omitted ...
```

The application's `main.go` now includes database initialisation code, as well as the code to implement a new business requirement:

* An HTTP `POST` request to `/send` containing a `{ "value" : string }` JSON must save the value to the database.

We also have an update for another business requirement. The requirement _was_:

* It responds with a text message containing a heart symbol ("`<3`") on requests to `/`.

And _now_ it's going to be:

* It responds with a text message containing a heart symbol ("`<3`") on requests to `/`, preceded by the count of messages, stored in the database, enclosed in the parentheses. Example output: `Hello, Docker! (7) <3`

This is the full text of the `main.go` file:

{% raw %}
```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := initStore()
	if err != nil {
		log.Fatalf("failed to initialise the store: %s", err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		return rootHandler(db, c)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/send", func(c echo.Context) error {
		return sendHandler(db, c)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

type Message struct {
	Value string `json:"value"`
}

func initStore() (*sql.DB, error) {

	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS message (value STRING PRIMARY KEY)"); err != nil {
		return nil, err
	}

	return db, nil
}

func rootHandler(db *sql.DB, c echo.Context) error {
	r, err := countRecords(db)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Hello, Docker! (%d) <3", r))
}

func sendHandler(db *sql.DB, c echo.Context) error {

	m := &Message{}

	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := crdb.ExecuteTx(context.Background(), db, nil,
		func(tx *sql.Tx) error {
			_, err := tx.Exec(
				"INSERT INTO message (value) VALUES ($1) ON CONFLICT (value) DO UPDATE SET value = excluded.value",
				m.Value,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return nil
		})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}

func countRecords(db *sql.DB) (int, error) {

	rows, err := db.Query("SELECT COUNT(*) FROM message")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		rows.Close()
	}

	return count, nil
}
```
{% endraw %}

The repository also includes the `Dockerfile`, which is almost exactly the same as the multi-stage `Dockerfile` introduced in the previous modules.

Regardless of whether we had updated the old example application, or checked out the new one, a new Docker image has to be built to reflect the changes to the application's source code.

### Build the application

We can build the image with the familiar `build` command:

```shell
$ docker build --tag docker-gs-ping-roach .
```

### Run the application

Now, let’s run our container. This time we’ll need to set some environment variables so that our application would know how to access the database. For now, we’ll do this right in the `docker run` command. Later we will see a more convenient method with Docker Compose.

> **Note**
>
> Since we are running our CockroachDB cluster in "insecure" mode, the value for the password can be anything. 
> 
> **Don't run in insecure mode in production, though!**

```shell
$ docker run \
  -it --rm -d \
  --network mynet \
  --name rest-server \
  -p 80:8080 \
  -e PGUSER=totoro \
  -e PGPASSWORD=myfriend \
  -e PGHOST=roach1 \
  -e PGPORT=26257 \
  -e PGDATABASE=mydb \
  docker-gs-ping-roach
```

### Test the application

Let’s test that our application is connected to the database and is able to save key-value pairs. We'll reuse the test input from the [Build your Go image](build-images.md) section:

```shell
$ curl --request POST \
  --url http://localhost:8000/add \
  --header 'content-type: application/json' \
  --data '{"key": "name", "value": "Docker"}'
```

You should receive the following json back from our service.

```json
{"name":"Docker"}
```

## Wind down everything

Remember, that we are running CockroachDB in "insecure" mode. Now that we had built and tested our application, it's time to wind everything down before moving on. You can list the containers that you are running with the `list` command:

```shell
$ docker container list
```

Now that you know the container IDs, you can use `docker container stop` and `docker container rm`, as demonstrated in the previous modules.

Please make sure that you stop the CockroachDB and `docker-gs-ping-roach` containers before moving on.

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
$ curl --request GET --url http://localhost:8000/
```

You should receive the following response:

```json
{"Status":"OK"}
```

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
