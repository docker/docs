---
title: Docker Scout quickstart
keywords: scout, supply chain, vulnerabilities, packages, cves, scan, analysis, analyze
description: Learn how to get started with Docker Scout to analyze images and fix vulnerabilities
---

Docker Scout analyzes image contents and generates a detailed report of packages
and vulnerabilities that it detects. It can provide you with
suggestions for how to remediate issues discovered by image analysis.

This guide takes a vulnerable container image and shows you how to use Docker
Scout to identify and fix the vulnerabilities, compare image versions over time,
and share the results with your team.

## Step 1: Setup

[This example project](https://github.com/docker/scout-demo-service) contains
a vulnerable Node.js application that you can use to follow along.

1. Clone its repository:

   ```console
   $ git clone https://github.com/docker/scout-demo-service.git
   ```

2. Move into the directory:

   ```console
   $ cd scout-demo-service
   ```

3. Make sure you're signed in to your Docker account,
   either by running the `docker login` command or by signing in with Docker Desktop.

4. Build the image and push it to a `<ORG_NAME>/scout-demo:v1`,
   where `<ORG_NAME>` is the Docker Hub namespace you push to.

   ```console
   $ docker build --push -t <ORG_NAME>/scout-demo:v1 .
   ```

## Step 2: Enable Docker Scout

Docker Scout analyzes all local images by default. To analyze images in
remote repositories, you need to enable it first.
You can do this from Docker Hub, the Docker Scout Dashboard, and CLI.
[Find out how in the overview guide](/scout).

1. Sign in to your Docker account with the `docker login` command or use the
   **Sign in** button in Docker Desktop.

2. Next, enroll your organization with Docker Scout, using the `docker scout enroll` command.

   ```console
   $ docker scout enroll <ORG_NAME>
       ✓ Successfully enrolled organization <ORG_NAME> with Docker Scout Free
   ```

3. Enable Docker Scout for your image repository with the `docker scout repo enable` command.

   ```console
   $ docker scout repo enable --org <ORG_NAME> <ORG_NAME>/scout-demo
   ```

## Step 3: Analyze image vulnerabilities

After building, use the `docker scout` CLI command to see vulnerabilities
detected by Docker Scout.

The example application for this guide uses a vulnerable version of Express.
The following command shows all CVEs affecting Express in the image you just
built:

```console
$ docker scout cves --only-package express
```

Docker Scout analyzes the image you built most recently by default,
so there's no need to specify the name of the image in this case.

Learn more about the `docker scout cves` command in the
[`CLI reference documentation`](/reference/cli/docker/scout/cves).

## Step 4: Fix application vulnerabilities

The fix suggested by Docker Scout is to update
the underlying vulnerable express version to 4.17.3 or later.

1. Update the `package.json` file with the new package version.

   ```diff
      "dependencies": {
   -    "express": "4.17.1"
   +    "express": "4.17.3"
      }
   ```

2. Rebuild the image with a new tag and push it to your Docker Hub repository:

   ```console
   $ docker build --push -t <ORG_NAME>/scout-demo:v2 .
   ```

Now, viewing the latest tag of the image in Docker Desktop, the Docker Scout
Dashboard, or CLI, you can see that you have fixed the vulnerability.

```console
$ docker scout cves --only-package express
    ✓ Provenance obtained from attestation
    ✓ Image stored for indexing
    ✓ Indexed 79 packages
    ✓ No vulnerable package detected


  ## Overview

                      │                  Analyzed Image                   
  ────────────────────┼───────────────────────────────────────────────────
    Target            │  mobywhale/scout-demo:v2                   
      digest          │  ef68417b2866                                     
      platform        │ linux/arm64                                       
      provenance      │ https://github.com/docker/scout-demo-service.git  
                      │  7c3a06793fc8f97961b4a40c73e0f7ed85501857         
      vulnerabilities │    0C     0H     0M     0L                        
      size            │ 19 MB                                             
      packages        │ 1                                                 


  ## Packages and Vulnerabilities

  No vulnerable packages detected

```

## Step 5: Evaluate policy compliance

While inspecting vulnerabilities based on specific packages can be useful,
it isn't the most effective way to improve your supply chain conduct.

Docker Scout also supports policy evaluation,
a higher-level concept for detecting and fixing issues in your images.
Policies are a set of customizable rules that let organizations track whether
images are compliant with their supply chain requirements.

Because policy rules are specific to each organization,
you must specify which organization's policy you're evaluating against.
Use the `docker scout config` command to configure your Docker organization.

```console
$ docker scout config organization <ORG_NAME>
    ✓ Successfully set organization to <ORG_NAME>
```

Now you can run the `quickview` command to get an overview
of the compliance status for the image you just built.
The image is evaluated against the default, out-of-the-box policies.

```console
$ docker scout quickview

...
Policy status  FAILED  (2/6 policies met, 2 missing data)

  Status │                  Policy                      │           Results
─────────┼──────────────────────────────────────────────┼──────────────────────────────
  ✓      │ No copyleft licenses                         │    0 packages
  !      │ Default non-root user                        │
  !      │ No fixable critical or high vulnerabilities  │    2C    16H     0M     0L
  ✓      │ No high-profile vulnerabilities              │    0C     0H     0M     0L
  ?      │ No outdated base images                      │    No data
  ?      │ Supply chain attestations                    │    No data
```

Exclamation marks in the status column indicate a violated policy.
Question marks indicate that there isn't enough metadata to complete the evaluation.
A check mark indicates compliance.

## Step 6: Improve compliance

The output of the `quickview` command shows that there's room for improvement.
Some of the policies couldn't evaluate successfully (`No data`)
because the image lacks provenance and SBOM attestations.
The image also failed the check on a few of the evaluations.

Policy evaluation does more than just check for vulnerabilities.
Take the `Default non-root user` policy for example.
This policy helps improve runtime security by ensuring that
images aren't set to run as the `root` superuser by default.

To address this policy violation, edit the Dockerfile by adding a `USER`
instruction, specifying a non-root user:

```diff
  CMD ["node","/app/app.js"]
  EXPOSE 3000
+ USER appuser
```

Additionally, to get a more complete policy evaluation result,
your image should have SBOM and provenance attestations attached to it.
Docker Scout uses the provenance attestations to determine how the image was
built so that it can provide a better evaluation result.

Before you can build an image with attestations,
you must enable the [containerd image store](../desktop/containerd.md)
(or create a custom builder using the `docker-container` driver).
The default image store doesn't support manifest lists,
which is how the provenance attestations are attached to an image.

Open **Settings** in Docker Desktop. Under the **General** section,
check the **Use containerd for pulling and storing images** option.
Note that changing the image store hides existing images and containers until you switch back.

With the containerd image store enabled, rebuild the image with a new `v3` tag.
This time, add the `--provenance=true` and `--sbom=true` flags.

```console
$ docker build --provenance=true --sbom=true --push -t <ORG_NAME>/scout-demo:v3 .
```

## Step 7: View in Dashboard

After pushing the updated image with attestations, it's time to view the
results through a different lens: the Docker Scout Dashboard.

1. Open the [Docker Scout Dashboard](https://scout.docker.com/).
2. Sign in with your Docker account.
3. Select **Images** in the left-hand navigation.

The images page lists your Scout-enabled repositories.
Select the image in the list to open the **Image details** sidebar.
The sidebar shows a compliance overview for the last pushed tag of a repository.

> **Note**
>
> If policy results haven't appeared yet, try refreshing the page.
> It might take a few minutes before the results appear if this is your
> first time using the Docker Scout Dashboard.

Inspect the **No outdated base images** policy.
This policy checks whether base images you use are up-to-date.
It currently has a non-compliant status,
because the example image uses an old version `alpine` as a base image.

Select the **View fix** button next to the policy name for details about the violation,
and recommendations on how to address it.
In this case, the recommended action is to enable
[Docker Scout's GitHub integration](./integrations/source-code-management/github.md),
which helps keep your base images up-to-date automatically.

> **Tip**
>
> You can't enable this integration for the demo app used in this guide.
> Feel free to push the code to a GitHub repository that you own,
> and try out the integration there!
{ .tip }

## Summary

This quickstart guide has scratched the surface on some of the ways
Docker Scout can support software supply chain management:

- How to enable Docker Scout for your repositories
- Analyzing images for vulnerabilities
- Policy and compliance
- Fixing vulnerabilities and improving compliance

## What's next?

There's lots more to discover, from third-party integrations,
to policy customization, and runtime environment monitoring in real-time.

Check out the following sections:

- [Image analysis](/scout/explore/analysis.md)
- [Data sources](/scout/advisory-db-sources)
- [Docker Scout Dashboard](/scout/dashboard)
- [Integrations](./integrations/_index.md)
- [Policy evaluation](./policy/_index.md)
