This tutorial walks you through the process of setting up and using Docker GitHub
Actions for building Docker images, and pushing images to Docker Hub. You will
complete the following steps:

1. Create a new repository on GitHub.
2. Define the GitHub Actions workflow.
3. Run the workflow.

To follow this tutorial, you need a Docker ID and a GitHub account.

### Step one: Create the repository

Create a GitHub repository and configure the Docker Hub secrets.

1. Create a new GitHub repository using
   [this template repository](https://github.com/dvdksn/clockbox/generate).

   The repository contains a simple Dockerfile, and nothing else. Feel free to
   use another repository containing a working Dockerfile if you prefer.

2. Open the repository **Settings**, and go to **Secrets and variables** > **Actions**.

3. Create a new secret named `DOCKERHUB_USERNAME` and your Docker ID as value.

4. Create a new
   [Personal Access Token (PAT)](/security/for-developers/access-tokens/#create-an-access-token)
   for Docker Hub. You can name this token `clockboxci`.

5. Add the PAT as a second secret in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

With your repository created, and secrets configured, you're now ready for
action!

### Step two: Set up the workflow

Set up your GitHub Actions workflow for building and pushing the image to Docker
Hub.

1. Go to your repository on GitHub and then select the **Actions** tab.
2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub actions workflow file in
   your repository, under `.github/workflows/main.yml` by default.

3. In the editor window, copy and paste the following YAML configuration.

   ```yaml
   name: ci

   on:
     push:
       branches:
         - "main"

   jobs:
     build:
       runs-on: ubuntu-latest
   ```

   - `name`: the name of this workflow.
   - `on.push.branches`: specifies that this workflow should run on every push
     event for the branches in the list.
   - `jobs`: creates a job ID (`build`) and declares the type of machine that
     the job should run on.

For more information about the YAML syntax used here, see
[Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions).

### Step three: Define the workflow steps

Now the essentials: what steps to run, and in what order to run them.


```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v4
      -
        name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          file: ./Dockerfile
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/clockbox:latest
```


The previous YAML snippet contains a sequence of steps that:

1. Checks out the repository on the build machine.
2. Signs in to Docker Hub, using the
   [Docker Login](https://github.com/marketplace/actions/docker-login) action and your Docker Hub credentials.
3. Creates a BuildKit builder instance using the
   [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx) action.
4. Builds the container image and pushes it to the Docker Hub repository, using
   [Build and push Docker images](https://github.com/marketplace/actions/build-and-push-docker-images).

   The `with` key lists a number of input parameters that configures the step:

   - `context`: the [build context](/build/building/context/).
   - `file`: filepath to the Dockerfile.
   - `push`: tells the action to upload the image to a registry after building
     it.
   - `tags`: tags that specify where to push the image.

Add these steps to your workflow file. The full workflow configuration should
look as follows:


```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v4
      -
        name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          file: ./Dockerfile
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/clockbox:latest
```


### Run the workflow

Save the workflow file and run the job.

1. Select **Commit changes...** and push the changes to the `main` branch.

   After pushing the commit, the workflow starts automatically.

2. Go to the **Actions** tab. It displays the workflow.

   Selecting the workflow shows you the breakdown of all the steps.

3. When the workflow is complete, go to your
   [repositories on Docker Hub](https://hub.docker.com/repositories).

   If you see the new repository in that list, it means the GitHub Actions
   successfully pushed the image to Docker Hub!