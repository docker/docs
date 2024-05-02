`healthcheck` declares a check that's run to determine whether or not the service containers are "healthy". It works in the same way, and has the same default values, as the
[HEALTHCHECK Dockerfile instruction](https://docs.docker.com/reference/dockerfile/#healthcheck)
set by the service's Docker image. Your Compose file can override the values set in the Dockerfile. 