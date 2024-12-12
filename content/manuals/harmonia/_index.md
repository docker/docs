---
title: Project Harmonia
description: Learn how you can run your applications in the cloud with Project Harmonia
keywords: run, cloud, docker desktop, resources
sitemap: false
params:
  sidebar:
    group: Products
aliases:
- /run-cloud/
---

{{% restricted title="Private preview" %}}
Project Harmonia is in Private preview.
{{% /restricted %}}

Project Harmonia brings the power of the cloud to your local development workflow. You can now run your applications in the cloud whilst continuing to use your existing tools and workflows and without worrying about local resource limitations. Project Harmonia also lets you share previews of your cloud-based applications for real-time feedback. 

## Set up 

To get started with Project Harmonia, you need to:

- Have a Docker account that's part of a Docker organization
- Email `run.cloud@docker.com` to get help with onboarding

## Quickstart

You can use Project Harmonia from the Docker Desktop Dashboard or from the CLI. 

This guide introduces you to essential commands and steps for creating, managing, and sharing a cloud engine. 

### Step one: Create a cloud engine

{{< tabs group="method" >}}
{{< tab name="Docker Desktop">}}

1. In the Docker Desktop Dashboard, navigate to the **Project Harmonia** tab. 
2. In the top right-hand corner, select **Create Cloud Engine**.
3. Fill out the creation form:
   - Enter `cloudengine` as the name
   - Choose an organization to associate the cloud engine with
   - Select the engine size and architecture

   Note that the **Switch Docker Context to use remote engine** is selected by default. The automatically switches you to your new cloud engine once it has been created. 
4. Select **Create**.

To verify creation, check the context switcher in the top-left corner of the Docker Desktop Dashboard; it should display `cloudengine`. You’re now ready to use it.

{{< /tab >}}
{{< tab name="CLI">}}

Run the following command: 

```console
$ docker harmonia engine create cloudengine --type "standard-amd64"  --use
```

This creates an engine called `cloudengine` and:
- Immediately switches you to the new cloud engine with the `--use` flag.
- Sets the engine size to standard and the engine's CPU architecture to amd64 with the `--type` flag. 

Project Harmonia supports the following values for `--type`:
- `standard-arm64`
- `standard-amd64` (default)
- `large-arm64`
- `large-amd64`
- `aiml-amd64`

Standard size engines have 2 CPU cores and 4GB RAM, large and AI/ML engines have 4 CPU cores and 8GB RAM.

To verify you're using the newly created cloud engine, run:

```console
$ docker context inspect
```

You should see the following:

```text
[
    {
        "Name": "cloudengine2",
...
```

{{< /tab >}}
{{< /tabs >}}

### Step two: Run and remove containers with the newly created cloud engine 

1.  Run an Nginx container in the cloud engine:
    ```console
    $ docker run -d --name cloudnginx -p 8080:80 nginx
    ```
    This maps the container's port `80` to the host's port `8080`. If port `8080` is already in use on your host, you can specify a different port.
2.  View the Nginx welcome page. Navigate to [`http://localhost:8080/`](http://localhost:8080/).
3.  Verify the running container:
    -   In the **Containers** tab in the Docker Desktop Dashboard, you should see your Nginx container listed.
    -   Alternatively, list all running containers in the cloud engine via the terminal:
        ```console
        $ docker ps
        ```
4.  Shut down the container:
    ```console
    $ docker kill cloudnginx
    ```

Running a container with a cloud engine is just as straightforward as running it locally.

### Step three: Create and switch to a new cloud engine

{{< tabs group="method" >}}
{{< tab name="Docker Desktop">}}

1. Create a new cloud engine:
   - Enter `cloudengine2` as the name
   - Choose an organization to associate the cloud engine with
   - Select the **Standard** engine size with the **AMD-64** architecture
   In the **Project Harmonia** view you should now see both `cloudengine` and `cloudengine2`. 
2. Switch between engines, also known as your Docker contexts. Use the context switcher in the top-left corner of the Docker Desktop Dashboard to toggle between your cloud engines or switch from your local engine (`desktop-linux`) to a cloud engine. 

{{< /tab >}}
{{< tab name="CLI">}}

1. Create a new cloud engine. Run:
   ```console
   $ docker harmonia engine create cloudengine2
   ```
   Docker automatically switches you to your new cloud engine. 
2. Switch between engines, also known as your Docker contexts. Either switch to your first cloud engine:
   ```console
   $ docker context use cloudengine
   ``` 
   Or switch back to your local engine: 
   ```console
   $ docker context use desktop-linux
   ```

{{< /tab >}}
{{< /tabs >}}

### Step four: Use a file sync for your cloud engine

Project Harmonia takes advantage of [Synchronized file shares](/manuals/desktop/features/synchronized-file-sharing.md) to enable local-to-remote file shares and port mappings. 

{{< tabs group="method" >}}
{{< tab name="Docker Desktop">}}

1. Clone the [Awesome Compose](https://github.com/docker/awesome-compose) repository.
2. In the Docker Desktop Dashboard, navigate to the **Project Harmonia** view.
3. For the `cloudengine` cloud engine, select the **Actions** menu and then **Manage file syncs**.
4. Select **Create file sync**.
5. Navigate to the `awesome-compose/react-express-mysql` folder and select **Open**.
6. In your terminal, navigate to the `awesome-compose/react-express-mysql` directory.
7. Run the project in the cloud engine with:
   ```console
   $ docker compose up -d
   ```
8. Test the application by visiting [`http://localhost:3000`](http://localhost:3000/).
   You should see the home page. The code for this page is located in `react-express-mysql/frontend/src/App.js`.
9. In an IDE or text editor, open the `App.js` file, change some text, and save. Watch as the code reloads live in your browser.

{{< /tab >}}
{{< tab name="CLI">}}

1. Clone the [Awesome Compose](https://github.com/docker/awesome-compose) repository.
2. In your terminal, change into the `awesome-compose/react-express-mysql` directory.
3. Create a file sync for `cloudengine`:
   ```console
   $ docker harmonia file-sync create --engine cloudengine $PWD
4. Run the project in the cloud engine with:
   ```console
   $ docker compose up -d
   ```
5. Test the application by visiting [`http://localhost:3000`](http://localhost:3000/).
   You should see the home page. The code for this page is located in `react-express-mysql/frontend/src/App.js`.
6. In an IDE or text editor, open the `App.js` file, change some text, and save. Watch as the code reloads live in your browser.

{{< /tab >}}
{{< /tabs >}}

### Step five: Share a container port

{{< tabs group="method" >}}
{{< tab name="Docker Desktop">}}

1.  Make sure your Docker context is set to `cloudengine`.
2.  In the Docker Desktop Dashboard, navigate to the **Containers** view.
3.  If necessary, expand the application listing to show all of its containers.
4.  Select the **lock** icon in the **Ports** column of your running container next to `3000:3000`.
    This creates a publicly accessible URL that you can share with teammates.
5.  Select the **copy** icon, to copy this URL.

To view all shared ports for your Docker context, select the **Shared ports** icon in the bottom-right corner of the Docker Desktop Dashboard.

{{< /tab >}}
{{< tab name="CLI">}}

To share a container port, make sure your Docker context is set to `cloudengine` and then run: 
``` console
$ docker harmonia engine share create cloudengine 3000
```
This returns a publicly accessible URL for your React app hosted on port `3000`, that you can share with teammates.

To see a list of all your shared ports, run:

```console
$ docker harmonia engine share list 
```

{{< /tab >}}
{{< /tabs >}}

### Step six: Clean up 

{{< tabs group="method" >}}
{{< tab name="Docker Desktop">}}

To stop the running project:

```console
$ docker compose down
```

To remove a file sync session:
1. Navigate to your cloud engine in the **Project Harmonia** view.
2. Select the **Actions** menu and then **Manage file syncs**.
3. Select the **drop-down** icon on the file sync.
4. Select **Delete**.

To remove a cloud engine, navigate to the **Project Harmonia** view and then select the **delete** icon.

{{< /tab >}}
{{< tab name="CLI">}}

To stop the running project:

```console
$ docker compose down
```

To remove the file sync session, run:

```console
$ docker harmonia file-sync delete --engine cloudengine $PWD
```

To remove a cloud engine, run:

```console
$ docker harmonia engine delete <name-of-engine>
```

{{< /tab >}}
{{< /tabs >}}

## Troubleshoot

Run `docker harmonia doctor` to print helpful troubleshooting information. 

## Known issues

- KinD does not run on Project Harmonia due to some hard-coded assumptions to ensure it's running in a privileged container. K3d is a good alternative.
- Containers cannot access host through DNS `host.docker.internal`.
- File binds (non-directory binds) are currently static, meaning changes will not be reflected until the container is restarted. This also affects Compose configs and secrets directives.
- Bind _mounts_, such as `-v /localpath:/incontainer` in the `docker run` command, require creating a file-sync.
- Creating a [synchronized file share](/manuals/desktop/features/synchronized-file-sharing.md) for a directory with a large amount of may take extra time to sync and become ready for use in a container.
- Bind _volumes_, such as those created with `docker volume create --driver local --opt type=none --opt o=bind --opt device=/some/host/path myvolname` or via the compose equivalent, are not supported.
- Port-forwarding for UDP is not supported. 
- Docker Compose projects relying on `watch` in `sync` mode are not working with the `tar` synchronizer. Configure it to use `docker cp` instead, disable tar sync by setting `COMPOSE_EXPERIMENTAL_WATCH_TAR=0` in your environment.
- Some Docker Engine features that let you access the underlying host, such as `--pid=host`, `--network=host`, and `--ipc=host`, are currently disabled.
