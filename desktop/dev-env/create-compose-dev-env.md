



## Start a sample Compose Dev Environment

You can also use Dev Environments to collaborate on any Docker Compose-based projects. For example, let's use the `compose-dev-env` project from the [Docker Samples](https://github.com/dockersamples/compose-dev-env){:target="_blank" rel="noopener" class="_"} GitHub repository.

  > **Note**
  >
  > When cloning a Git repository using SSH, ensure you've added your SSH key to the ssh-agent. To do this, open a terminal and run `ssh-add <path to your private ssh key>`.

1. Under Dev Environments, click **Create**. This opens the **Create a Dev Environment** dialog. Click **Get Started** and then copy `https://github.com/dockersamples/compose-dev-env.git` and add it to the **Enter the Git Repository** field on the **Existing Git repo** source.
2. Click **Continue**. This initializes the project and clones the Git code and builds the Compose application. This:

    - Builds local images for services that are defined in the Compose file
    - Pulls images required for other services
    - Creates volumes and networks
    - Starts the Compose stack

![Compose Dev environment](images/dev-env-compose.png){:width="700px"}

Now your application is up and running, you can check by opening [http://localhost:8080](http://localhost:8080) in your browser.

The time taken to start the Compose application depends on how your application is configured, whether the images have been built, and the number of services you have defined, etc.

You'll also notice that VS Code doesn't open directly (unlike the single container Dev Environment) as there are multiple services configured. You can hover over a service and then click on the **Open in VS Code** button to open a specific service in Visual Studio Code. This stops the existing container and creates a new container which allows you to develop and update your service in VS Code.

You can now update your service and test it against your Compose application.

## Set up your own Compose Dev Environment

In the previous section, we've learnt how to start a sample Compose Dev Environment. To set up a Dev Environment for your own Compose-based project, you'll need some extra configuration that tells Docker Desktop how to build, start, and use the right Dev Environment image for your services.

Dev Environments use an additional `docker-compose.yaml` file located in the `.docker` directory at the root of your project. This file allows you to define the image required for a dedicated service, the ports you'd like to expose, along with additional configuration options dedicated to Dev Environments coming in the future.

Let's take a detailed look at the `docker-compose.yaml` file we used in the [compose-dev-env](https://github.com/dockersamples/compose-dev-env/blob/main/.docker/docker-compose.yaml){:target="_blank" rel="noopener" class="_"} sample project.

```yaml
version: "3.7"
services:
  backend:
    build:
      context: backend
      target: development
    secrets:
      - db-password
    depends_on:
      - db
  db:
    image: mariadb
    restart: always
    healthcheck:
      test: [ "CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "--silent" ]
      interval: 3s
      retries: 5
      start_period: 30s
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/mysql
    environment:
      - MYSQL_DATABASE=example
      - MYSQL_ROOT_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 3306
  proxy:
    build: proxy
    ports:
      - 8080:80
    depends_on:
      - backend
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

In the above yaml file, the build context `backend` specifies that that the container should be built using the `development` stage (`target` attribute) of the Dockerfile located in the `backend` directory (`context` attribute)

The `development` stage of the Dockerfile is defined as follows:

```dockerfile
FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

RUN go build -o server .

FROM build AS development
RUN apk update \
    && apk add git
CMD ["go", "run", "main.go"]

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]
```

The `development`target uses a `golang:1.16-alpine` image with all dependencies you need for development. You can start your project directly from VS Code and interact with the others applications or services such as the database or the frontend.

In our example, the Docker Compose files are the same. However, they could be different and the services defined in the main Compose file may use other targets to build or directly reference other images.



