---
title: Integrate Docker Scout with Sysdig
linkTitle: Sysdig
description: Integrate your runtime environments with Docker Scout using Sysdig
keywords: scout, sysdig, integration, image analysis, environments, supply chain
---

{{% include "scout-early-access.md" %}}

The Sysdig integration enables Docker Scout to automatically detect the images
you're using for your running workloads. Activating this integration gives you
real-time insights about your security posture, and lets you compare your
builds with what's running in production.

## How it works

The Sysdig Agent captures the images of your container workloads. Docker Scout
integrates with the Sysdig API to discover the images in your cluster. This
integration uses Sysdig's Risk Spotlight feature. For more information, see
[Risk Spotlight Integrations (Sysdig docs)](https://docs.sysdig.com/en/docs/sysdig-secure/integrations-for-sysdig-secure/risk-spotlight-integrations/).

> [!TIP]
>
> Sysdig offers a free trial for Docker users to try out the new Docker Scout integration.
>
> {{< button url=`https://sysdig.com/free-trial-for-docker-customers/` text="Sign up" >}}

Each Sysdig integration maps to an environment. When you enable a Sysdig
integration, you specify the environment name for that cluster, such as
`production` or `staging`. Docker Scout assigns the images in the cluster to
the corresponding environment. This lets you use the environment filters to see
vulnerability status and policy compliance for an environment.

Only images analyzed by Docker Scout can be assigned to an environment. The
Sysdig runtime integration doesn't trigger image analysis by itself. To analyze
images automatically, enable a [registry integration](../_index.md#container-registries).

Image analysis must not necessarily precede the runtime integration, but the
environment assignment only takes place once Docker Scout has analyzed the
image.

## Prerequisites

- Install the Sysdig Agent in the cluster that you want to integrate, see [Install Sysdig Agent (Sysdig docs)](https://docs.sysdig.com/en/docs/installation/sysdig-monitor/install-sysdig-agent/).
- Enable profiling for Risk Spotlight Integrations in Sysdig, see [Profiling (Sysdig docs)](https://docs.sysdig.com/en/docs/sysdig-secure/policies/profiling/#enablement).
- You must be an organization owner to enable the integration in the Docker Scout Dashboard.

## Integrate an environment

1. Go to the [Sysdig integration page](https://scout.docker.com/settings/integrations/sysdig/)
   on the Docker Scout Dashboard.
2. In the **How to integrate** section, enter a configuration name for this
   integration. Docker Scout uses this label as a display name for the
   integration.

3. Select **Next**.

4. Enter a Risk Spotlight API token and select the region in the drop-down list.

   The Risk Spotlight API token is the Sysdig token that Docker Scout needs to
   integrate with Sysdig. For more instructions on how to generate a Risk
   Spotlight token, See [Risk Spotlight Integrations (Sysdig docs)](https://docs.sysdig.com/en/docs/sysdig-secure/integrations-for-sysdig-secure/risk-spotlight-integrations/docker-scout/#generate-a-token-for-the-integration).

   The region corresponds to the `global.sysdig.region` configuration parameter
   set when deploying the Sysdig Agent.

5. Select **Next**.

   After selecting **Next**, Docker Scout connects to Sysdig and retrieves the
   cluster names for your Sysdig account. Cluster names correspond to the
   `global.clusterConfig.name` configuration parameter set when deploying
   Sysdig Agents.

   An error displays if Docker Scout fails to connect to Sysdig using the
   provided token. If there's an error, you won't be able to continue the
   integration. Go back and verify that the configuration details are correct.

6. Select a cluster name in the drop-down list.

7. Select **Next**.

8. Assign an environment name for this cluster.

    You can reuse an existing environment or create a new one.

9. Select **Enable integration**.

After enabling the integration, Docker Scout automatically detects images
running in the cluster, and assigns those images to the environment associated
with the cluster. For more information about environments, see [Environment
monitoring](./_index.md).

> [!NOTE]
>
> Docker Scout only detects images that have been analyzed. To trigger an image
> analysis, enable a [registry integration](../_index.md#container-registries)
> and push an image to your registry.
>
> If you created a new environment for this integration, the environment
> appears in Docker Scout when at least one image has been analyzed.

To integrate more clusters, go to the [Sysdig integrations page](https://scout.docker.com/settings/integrations/sysdig/)
and select the **Add** button.
