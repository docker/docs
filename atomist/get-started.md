---
title: Getting started
description: Getting started with Atomist
keywords: atomist, software supply chain, vulnerability scanning, tutorial
toc_max: 2
---

{% include atomist/disclaimer.md %}

To get started with Atomist, you'll need to:

- Connect Atomist with your container registry
- Link your container images with their Git source

Before you can begin the setup, you’ll need a Docker ID. If you do not already
have one, you can [sign up here](https://hub.docker.com/signup){: target="blank"
rel="noopener" class=""}.

## Connect container registry

After completing this setup, Atomist will have read-only access to your
registry, and gets notified about pushed or deleted images.

Follow the applicable instructions depending on the type of container registry
you use.

<ul class="nav nav-tabs">
  <li><a data-toggle="tab" data-target="#tab-hub">Docker Hub</a></li>
  <li><a data-toggle="tab" data-target="#tab-ecr">Amazon ECR</a></li>
  <li><a data-toggle="tab" data-target="#tab-google">GCR/GAR</a></li>
  <li><a data-toggle="tab" data-target="#tab-ghcr">GitHub Container Registry</a></li>
  <li><a data-toggle="tab" data-target="#tab-jfrog">JFrog Artifactory</a></li>
</ul>
<div class="tab-content"><br>
<div id="tab-hub" class="tab-pane fade in active">
  <p>If you are using Docker Hub as your container registry, you can skip this step
  and go straight to <a href="#link-images-to-git-repository">linking images to Git source</a>. Atomist
  integrates seamlessly with your Docker Hub organizations.</p>
</div>
<div id="tab-ecr" class="tab-pane fade" markdown="1">
<!-- ECR -->

When setting up an Amazon Elastic Container Registry (ECR) integration with
Atomist, we need to create the following resources on the AWS side:

- Read-only IAM role, for Atomist to be able to access the container registry
- Amazon EventBridge, to notify Atomist of pushed and deleted images

To help you get started, we have created a public CloudFormation template. This
template will create an IAM role and Amazon EventBridge.

Our CloudFormation templates protects you from
[confused deputy attacks](https://docs.aws.amazon.com/IAM/latest/UserGuide/confused-deputy.html){:
target="blank" rel="noopener" class=""} by ensuring a unique `ExternalId`, along
with the appropriate condition on the IAM role statement.

1. Go to <https://dso.docker.com> and sign in using your Docker ID credentials.
2. Navigate to the **Integrations** tab and click **Configure** next to the
   **Elastic Container Registry** integration.
3. Fill out all the fields, except **Trusted Role ARN**. You will only know the
   role ARN after applying the CloudFormation template.

   Choose basic auth credentials to protect the endpoint that AWS will use to
   notify Atomist. The URL and the basic auth credentials will be used as
   parameters in the CloudFormation template.

4. Now we'll create the CloudFormation stack. Before creating the stack, AWS
   will ask you to enter three parameters.

   - `Url`: the API endpoint copied from Atomist
   - `Username`, `Password`: basic authentication credentials for the endpoint.
     Must match what you entered in the Atomist workspace.

   Click on one of the **Launch Stack** buttons below to start reviewing the
   details in your AWS account.

   > Note
   >
   > Before creating the stack, AWS will ask for acknowledgement that creating
   > this stack requires a capability. This stack creates a role that will grant
   > Atomist read-only access to ECR resources.
   >
   > ![confirm](./images/ecr/capability.png)

   <div style="text-align: center">
     <table>
       <tr>
         <th>Region</th>
         <th>ecr-integration.template</th>
       </tr>
       <tr>
         <th>us-east-1</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-us-east-1/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>us-east-2</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=us-east-2#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-us-east-2/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>us-west-1</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=us-west-1#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-us-west-1/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>us-west-2</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-us-west-2/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>eu-west-1</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=eu-west-1#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-eu-west-1/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>eu-west-2</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=eu-west-2#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-eu-west-2/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>eu-west-3</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=eu-west-3#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-eu-west-3/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>eu-central-1</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=eu-central-1#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-eu-central-1/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>ca-central-1</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=ca-central-1#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-ca-central-1/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
       <tr>
         <th>ap-southeast-2</th>
         <td>
           <a href="https://console.aws.amazon.com/cloudformation/home?region=ap-southeast-2#/stacks/new?stackName=atomist-public-templates-ecr-integration&templateURL=https://s3.amazonaws.com/atomist-ap-southeast-2/atomist-public-templates/latest/ecr-integration.template">
             <img alt="Launch Stack" src="https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png" />
           </a>
         </td>
       </tr>
     </table>
   </div>

5. Once the stack has been created, copy the **Value** for the **AssumeRoleArn**
   key.

   ![AWS stack creation output](./images/ecr/stackoutput.png){: width="700px"}

6. Paste the copied **AssumeRoleArn** value into the **Trusted Role ARN** field
   on the Atomist configuration page.

7. Click **Save Configuration**. Atomist will now test the connection with your
   ECR registry. You'll see a green check mark beside the integration if a
   successful connection was made.

   ![integration list showing a successful ECR integration](./images/ecr/ConnectionSuccessful.png){:
   width="700px"}

</div>
<div id="tab-google" class="tab-pane fade" markdown="1">
<!-- GCR/GAR -->

Setting up an Atomist integration with Google Container Registry (GCR) and
Google Artifact Registry (GAR) involves:

- Creating a service account and grant it the required role.
- Creating a PubSub subscription on the `gcr` topic to watch for activity in the
  registry.

To complete the following procedure requires administrator's permissions in the
project.

1. Set the following variables in your shell session. They will be used in
   subsequent steps when configuring the Google Cloud resources, using the
   `gcloud` CLI. The `SERVICE_ACCOUNT_ID` can be set to whatever you'd like.

   ```bash
   export SERVICE_ACCOUNT_ID="atomist-gar-integration"
   export PROJECT_ID="YOUR_GCP_PROJECT_ID"
   ```

2. Create the service account.

   ```bash
   gcloud iam service-accounts create ${SERVICE_ACCOUNT_ID} \
       --project ${PROJECT_ID} \
       --description="Atomist GAR Integration Service Account" \
       --display-name="Atomist GAR Integration"
   ```

3. Grant the service account read-only access to the artifact registry.

   ```bash
   gcloud projects add-iam-policy-binding ${PROJECT_ID} \
       --project ${PROJECT_ID} \
       --member="serviceAccount:${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com" \
       --role="roles/artifactregistry.reader"
   ```

4. Grant service account access to Atomist.

   ```bash
   gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com" \
       --project ${PROJECT_ID} \
       --member="serviceAccount:atomist-bot@atomist.iam.gserviceaccount.com" \
       --role="roles/iam.serviceAccountTokenCreator"
   ```

5. Go to <https://dso.docker.com> and sign in using your Docker ID credentials.
6. Navigate to the **Integrations** tab and click **Configure** next to the
   **Google Artifact Registry** integration.
7. Fill out the following fields:

   - **Project ID** is the `PROJECT_ID` used in previous steps.
   - **Service Account**: The email address of the service account created
     step 2.

8. Click **Save Configuration**. Atomist will test the connection. You'll see
   some green check marks when a connection has been established.

   ![GAR configuration successful](./images/gar/config_success.png){:
   width="700px"}

   Next, we will create a new PubSub subscription on the `gcr` topic in Google
   Artifact Registry. This subscription will notify Atomist when new images are
   pushed to the registry.

9. Copy the URL in the **GCR Events Webhook** field to your clipboard. This will
   be the `PUSH_ENDPOINT_URI` for the PubSub subscription.

10. Define the following three variable values, in addition to the `PROJECT_ID`
    and `SERVICE_ACCOUNT_ID` from earlier:

    - `PUSH_ENDPOINT_URL`: the webhook URL copied from the Atomist workspace.
    - `SERVICE_ACCOUNT_EMAIL`: the service account address; a combination of the
      service account ID and project ID.
    - `SUBSCRIPTION`: the name of the PubSub to be created (can be anything).

    ```bash
    PUSH_ENDPOINT_URI={COPY_THIS_FROM_ATOMIST}
    SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com"
    SUBSCRIPTION="atomist-gar-integration-subscription"
    ```

11. Create the PubSub for the `gcr` topic.

    ```bash
    gcloud pubsub subscriptions create ${SUBSCRIPTION} \
      --topic='gcr' \
      --push-auth-token-audience='atomist' \
      --push-auth-service-account="${SERVICE_ACCOUNT_EMAIL}" \
      --push-endpoint="${PUSH_ENDPOINT_URI}"
    ```

When the first image push is successfully detected, a green check mark on the
integration page will indicate that the webhook event was received and that the
integration works.

</div>
<div id="tab-ghcr" class="tab-pane fade" markdown="1">
<!-- GitHub Container Registry -->

To integrate Atomist with GitHub Container Registry, connect your GitHub
account, and enter a personal access token for Atomist to use when pulling
container images.

1. Go to <https://dso.docker.com> and sign in using your Docker ID credentials.
2. Connect your GitHub account as instructed in the
   [GitHub app page](./integrate/github.md#connect-to-github). Install the app
   into the GitHub account that contains your GitHub Container Registry.
3. Open the [**Integrations**](https://dso.docker.com/r/auth/integrations){:
   target="blank" rel="noopener" class=""} tab, and click the **Configure** link
   next to the **GitHub Container Registry** in the list of integrations.
4. Fill out the fields and click **Save Configuration**.

   The **Personal access token** is required for scanning images in private
   repositories. The token must have the
   [`read:packages` scope](https://docs.github.com/en/packages/learn-github-packages/about-permissions-for-github-packages).

   You can leave the **Personal access token** field blank if you only want to
   scan images in public repositories.

</div>
<div id="tab-jfrog" class="tab-pane fade" markdown="1">
<!-- JFrog Artifactory -->

Atomist can index images in a JFrog Artifactory repository by means of a
monitoring agent.

The agent scans configured repositories at regular intervals, and send newly
discovered images' metadata to the Atomist data plane.

In the following example, `https://hal9000.atomist.com` is a private registry
only visible on an internal network.

```
docker run -ti atomist/docker-registry-broker:latest\
  index-image remote \
  --workspace AQ1K5FIKA \
  --api-key team::6016307E4DF885EAE0579AACC71D3507BB38E1855903850CF5D0D91C5C8C6DC0 \
  --artifactory-url https://hal9000.docker.com \
  --artifactory-repository atomist-docker-local \
  --container-registry-host atomist-docker-local.hal9000.docker.com
  --username admin \
  --password password
```

| Parameter                 | Description                                                                                                     |
| ------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `workspace`               | ID of your Atomist workspace.                                                                                   |
| `api-key`                 | Atomist API key.                                                                                                |
| `artifactory-url`         | Base URL of the Artifactory instance. Must not contain trailing slashes.                                        |
| `artifactory-repository`  | The name of the container registry to watch.                                                                    |
| `container-registry-host` | The hostname associated with the Artifactory repository containing images, if different from `artifactory-url`. |
| `username`                | Username for HTTP basic authentication with Artifactory.                                                        |
| `password`                | Password for HTTP basic authentication with Artifactory.                                                        |

</div>
<hr>
</div>

## Link images to Git repository

Knowing the source repository of an image is a prerequisite for Atomist to
interact with the Git repository. For Atomist to be able to link scanned images
back to a Git repository repository, you must annotate the image at build time.

The image labels that Atomist requires are:

| Label                                | Value                                             |
| ------------------------------------ | ------------------------------------------------- |
| `org.opencontainers.image.revision`  | The commit revision that the image is built for.  |
| `com.docker.image.source.entrypoint` | Path to the Dockerfile, relative to project root. |

For more information about pre-defined OCI annotations, see the
[specification document on GitHub](https://github.com/opencontainers/image-spec/blob/main/annotations.md#pre-defined-annotation-keys).

There are different ways of adding these labels to an image.

### Add labels using Docker Buildx

> Beta
>
> Git provenance labels via Buildx is a beta feature.

To add the image labels using Docker Buildx, set the environment variable
`BUILDX_GIT_LABELS=1`:

```bash
export BUILDX_GIT_LABELS=1
docker buildx build . -f docker/Dockerfile
```

### Add labels using the label CLI argument

You can create labels using the `--label` argument for `docker build`.

```bash
docker build . -f docker/Dockerfile -t $IMAGE_NAME \
    --label "org.opencontainers.image.revision=10ac8f8bdaa343677f2f394f9615e521188d736a" \
    --label "com.docker.image.source.entrypoint=docker/Dockerfile"
```

Images built in a CI/CD environment can leverage the built-in environment
variables. For example, to set the `org.opencontainers.image.revision` in GitHub
Actions, you can use {% raw %}`${{ github.sha }}`{% endraw %}. Consult the
documentation for your CI/CD platform to learn which variables to use.

### Add labels in the Dockerfile

You can specify the labels directly in the Dockerfile using the `LABEL` command,
should you want to.

```dockerfile
LABEL org.opencontainers.image.revision="10ac8f8bdaa343677f2f394f9615e521188d736a"
LABEL com.docker.image.source.entrypoint="docker/Dockerfile"
```

## Where to go next

Atomist is now tracking bill of materials, packages, and vulnerabilities for
your images! You can view your image scan results on the
[images overview page](https://dso.docker.com/r/auth/overview/images).

Teams use Atomist to protect downstream workloads from new vulnerabilities. It's
also used to help teams track and remediate new vulnerabilities that impact
existing workloads. The following sections describe integrate and configure
Atomist further. For example, to gain visibility into container workload systems
like Kubernetes.

- Connect Atomist with your GitHub repositories by
  [installing the Atomist app](./integrate/github.md) for your GitHub
  organization.
- Manage which Atomist features you use in [settings](./configure/settings.md).
- Learn about [deployment tracking](integrate/deploys.md) and how Atomist can
  help monitor your deployed containers.
- Atomist watches for new advisories from public sources, but you can also
  [add your own internal advisories](reference/advisories.md) for more
  information.
