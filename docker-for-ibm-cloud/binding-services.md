---
description: Bind IBM Cloud services to swarms
keywords: ibm, ibm cloud, services, watson, AI, IoT, iaas, tutorial
title: Bind IBM Cloud services to swarms
---

With Docker EE for IBM Cloud, you can easily bind services to your cluster to enhance your apps with Watson, AI, Internet of Things, and other services available in the IBM Cloud catalog.

## Bind IBM Cloud services
Before you begin:

* Ensure that you have [set up your IBM Cloud account](/docker-for-ibm-cloud/index.md).
* [Install the IBM Cloud CLI and plug-ins](/docker-for-ibm-cloud/index.md#install-the-clis).
* [Create a cluster](administering-swarms.md).
* Get the name of the cluster to which you want to bind the service by running `bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key`.
* Identify an existing or create a new IBM Cloud service.  To list existing services, run `bx service list`.
* Identify an existing registry namespace or create a registry namespace. [IBM Cloud Container Registry example](https://console.bluemix.net/docs/services/Registry/registry_setup_cli_namespace.html#registry_namespace_add).
* Review example files that are customized for an IBM Watson Conversation service, to give you an idea of how you might develop your own service files. **Note**: The steps include examples for another type of service to give you ideas for other ways you might build your files.
   * Example [Dockerfile](https://github.com/docker/docker.github.io/tree/master/docker-for-ibm-cloud/scripts/Dockerfile)
   * Example [docker-service.yaml file](https://github.com/docker/docker.github.io/tree/master/docker-for-ibm-cloud/scripts/docker-stack.yaml)

There are three main steps in binding IBM Cloud services to your Docker EE for IBM Cloud cluster:

1. Create a Docker secret.
2. Build a Docker image that uses the IBM Cloud service.
3. Create a Docker service.

### Step 1: Create a Docker secret

1. Log in to IBM Cloud. If you have a federated account, use the `--sso` option.

   ```bash
   $ bx login [--sso]
   ```

2. Target the org and space that has the service:

   ```bash
   $ bx target --cf
   ```

3. Create the Docker secret for the service. The `--swarm-name` is the cluster that you're binding the service to. The `--service-name` flag must match the name of your IBM Cloud service. The `--service-key` flag is used to create the Docker service YAML file. The `--cert-path` is the filepath to your cluster's UCP client bundle certificates. Include your IBM Cloud infrastructure credentials if you have not set the environment variables.

   ```bash
   $ bx d4ic key-create --swarm-name my_swarm \
   --service-name my_ibm_service \
   --service-key my_secret \
   --cert-path filepath/to/certificate/repo \
   --sl-user user.name.1234567 \
   --sl-api-key api_key
   ```

4. Verify the secret is created:

   ```bash
   $ docker secret ls
   ```

5. Update your service code to use the secret that you created. For example:

   ```none
   {% raw %}
   ...
   // WatsonSecret holds Watson VR service keys
   type WatsonSecret struct {
	   URL    string `json:"url"`
	   Note   string `json:"note"`
	   APIKey string `json:"api_key"`
   }
   ...
   var watsonSecretName = "watson-secret"
   var watsonSecrets WatsonSecret
   ...
	   watsonSecretFile, err := ioutil.ReadFile("/run/secrets/" + watsonSecretName)
	   if err != nil {
		   fmt.Println(err.Error())
		   os.Exit(1)
	   }

	   json.Unmarshal(watsonSecretFile, &watsonSecrets)
	   fmt.Println("Watson URL: ", watsonSecrets.URL)
   ...
		   msgQ.Add("api_key", watsonSecrets.APIKey)
   ...
   {% endraw %}
   ```

### Step 2: Build a Docker image
1. Log in to the registry that you are using to store the image.

2. Create a Dockerfile following [Dockerfile best practices](/engine/userguide/eng-image/dockerfile_best-practices/).

   > Docker images
   >
   > If you are unfamiliar with Docker images, try the [Getting Started](/get-started/).

   **Example** snippet for a Dockerfile that uses _mmssearch_ service.

   ```none
   {% raw %}
   FROM golang:latest
   WORKDIR /go/src/mmssearch
   COPY . /go/src/mmssearch
   RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=0 go/src/mmssearch/main .
   CMD ["./main"]
   LABEL version=demo-3
   {% endraw %}
   ```

3. Navigate to the directory of the Dockerfile, and build the image. Don't forget the period in the `docker build` command.

   ```bash
   $ cd directory/path && docker build -t my_image_name .
   ```

4. Test the image locally before pushing to your registry.

5. Tag the image:

   ```bash
   $ docker tag my_image_name registry-path/namespace/image:tag
   ```

6. Push the image to your registry:

   ```bash
   $ docker push registry-path/namespace/image:tag
   ```

### Step 3: Create a Docker service

1. Develop a `docker-service.yaml` file using the [compose file reference](/compose/compose-file/).

   * Save the file in an easily accessible directory, such as the one that has the Dockerfile that you used in the previous step.
   * For the `image` field, use the same `registry/namespace/image:tag` path that you made in the previous step for the for the service `image` field.
   * For the service `environment` field, use a service environment, such as a workspace ID, from the IBM Cloud service that you made before you began.
   * **Example** snippet for a `docker-service.yaml` that uses _mmssearch_ with a Watson secret.

   ```none
   {% raw %}
   mmssearch:
     image: mmssearch:latest
     build: .
     ports:
       - "8080:8080"
     secrets:
       - source: watson-secret
         target: watson-secret
   secrets:
     watson-secret:
     external: true
   {% endraw %}
   ```

2. Connect to your cluster by setting the environment variables from the [client certificate bundle that you downloaded](administering-swarms.md#download-client-certificates).

   ```bash
   $ cd filepath/to/certificate/repo && source env.sh
   ```

3. Navigate to the directory of the `docker-service.yaml` file.

4. Deploy the service:

   ```bash
   $ docker stack deploy my_service_name \
   --with-registry-auth \
   --compose-file docker-stack.yaml
   ```

5. Verify that the service has been deployed to your cluster:

   ```bash
   $ docker service ls
   ```
