# Using compose with UCP

## Creating Applications using the UI

The 'Create Application' option in the Applications screen of the UI allows a user to perform a one-off compose deployment by entering a `docker-compose.yml` file.  The `docker-compose.yml` can be provided by entering it into the text area or uploading a file using the upload option below it.

When a compose application is created using the UI, a `docker/ucp-compose` user container is launched with several environment variables:

  - `DOCKER_COMPOSE_YML`: base64 encoded `docker-compose.yml` string
  - `PROJECT_NAME`: The 'Application Name' that is passed to compose's `--project-name` option
  - `CONTROLLER_HOST`: The UCP hostname for compose to connect to the cluster
  - `CONTROLLER_CA_CERT`: base64 encoded CA cert that compose should use to validate the server cert of `CONTROLLER_HOST`
  - `JWT`: The user's authentication token that is used for client authentication with the UCP controller, this is written to `~/.docker/config.json` as a custom HTTP header

the `docker-compose.yml`.  The container launches `docker-compose` that is configured to connect to the UCP cluster and launch the application.

When launching an application using the UI, there are some restrictions on Compose keywords, this is due to receiving the docker-compose.yml in isolation.  The UI will highlight unsupported keywords in red to remind the user that the keyword is unavailable.

The restrictions are because compose is running inside a container on your cluster and it will not have access to:

 - The working directory, so a build context cannot be created using 'build'
 - A Dockerfile on your local filesystem specified with `dockerfile`
 - An env file on your local filesystem specified with `env_file`


### FAQ

#### What does a 'one-off' deployment mean?

UCP does not save the `docker-compose.yml` file, and it doesn't save any state regarding the application.  At the moment it will not monitor an application or make updates to it without the user performing another deployment.


#### Are application names namespaced?

Applications names are not currently namespaced, so the user will need to ensure that when deploying an application that they use a unique application name to prevent compose from attempting to mmanipulate the configuration of another application running on the cluster with the same name.


#### Why are environment variables used by the UI to pass in authentication information?  Does this mean sensitive information could be leaked?

Since the `docker/ucp-compose` container runs as a user container, the assumption is that only the user that launched the container (and admin users) should have be able to view the container while it's running.

We had investigated a couple of other options to authenticate the compose container with the UCP controller using the current middleware:

 - **Stream a client bundle into the container once it starts, and have the wrapper script wait for the file on stdin:**  The issue with this is that client bundle cleanup is non-trivial, since we can only identify a client bundle in a user's account via the public key.  This means that the frontend would need to cache the public key from inside the client bundle zip, for when the compose container finishes, and then use that to make relevant API calls to remove the temporary bundle.
 - **Stream the JWT into the container:** Streaming content into the container requires exec privileges, which means that users with "write" privileges would be unable to deploy an application

There is room for improvement around how the compose data is provided to the container, ideas have been extending our middleware to allow for one-time (or short TTL) JWTs and client bundles.


#### Why does the UI launch a compose container?

The current understanding is that libcompose isn't ready to be integrated into our backend yet so that it can be ran natively as part of the controller, instead of a sideband container.  We'll be looking to integrate this natively in future.


#### What's in the `ucp/compose` image?

The `ucp/compose` image is based on the official `docker/compose` image that's released by the core team.  This is based on `alpine:edge` and contains python plus python dependencies that compose requires to run.  The compose team have investigated making this a standalone static linked binary already, but this is probably something we could re-visit to slim down this container.

Our `ucp/compose` image adds a wrapper script that expects environment variables to be set, and then decodes and writes them to their expected locations on the filesystem, then runs compose against the UCP controller.

With regards to size, the `docker/ucp-compose` image is of comparable size to our other alpine based images at 58 MB.


