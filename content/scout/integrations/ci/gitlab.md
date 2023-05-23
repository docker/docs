---
description: How to integrate Docker Scout with GitLab CI
keywords: supply chain, security, ci, continuous integration, gitlab
title: Integrate Docker Scout with GitLab CI
---

{% include scout-early-access.md %}

The following examples runs in GitLab CI in a repository containing a Docker
image's definition and contents. Triggered by a commit, the pipeline builds the
image. If the commit was to the default branch, it uses Docker Scout to get a
CVE report. If the commit was to a different branch, it uses Docker Scout to
compare the new version to the current published version.

First, set up the rest of the workflow. There's a lot that's not specific to
Docker Scout but needed to create the images to compare.

Add the following to a _.gitlab-ci.yml_ file at the root of your repository.

```yaml
docker-build:
  image: docker:latest
  stage: build
  services:
    - docker:dind
  before_script:
    - docker login -u "$CI_REGISTRY_USER" -p "$CI_REGISTRY_PASSWORD" $CI_REGISTRY

    # Install curl and the Docker Scout CLI
    - |
      apk add --update curl
      curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s -- 
      apk del curl 
      rm -rf /var/cache/apk/*
    # Login to Docker Hub required for Docker Scout CLI
    - docker login -u "$DOCKER_HUB_USER" -p "$DOCKER_HUB_PAT"
```

This sets up the workflow to build Docker images with GitLab's
"Docker-in-Docker" mode to run Docker inside a container.

It then downloads curl and the Docker CLI and logs into the GitLab CI registry
and the Docker registry using environment variables defined in your repository's
settings.

Add the following to the YAML file:

{% raw %}
```yaml
script:
  - |
    if [[ "$CI_COMMIT_BRANCH" == "$CI_DEFAULT_BRANCH" ]]; then
      tag=""
      echo "Running on default branch '$CI_DEFAULT_BRANCH': tag = 'latest'"
    else
      tag=":$CI_COMMIT_REF_SLUG"
      echo "Running on branch '$CI_COMMIT_BRANCH': tag = $tag"
    fi
  - docker build --pull -t "$CI_REGISTRY_IMAGE${tag}" .
  - |
    if [[ "$CI_COMMIT_BRANCH" == "$CI_DEFAULT_BRANCH" ]]; then
      # Get a CVE report for the built image and fail the pipeline when critical or high CVEs are detected
      docker scout cves "$CI_REGISTRY_IMAGE${tag}" --exit-code --only-severity critical,high    
    else
      # Compare image from branch with latest image from the default branch and fail if new critical or high CVEs are detected
      docker scout compare "$CI_REGISTRY_IMAGE${tag}" --to "$CI_REGISTRY_IMAGE:latest" --exit-code --only-severity critical,high --ignore-unchanged
    fi

  - docker push "$CI_REGISTRY_IMAGE${tag}"
```
{% endraw %}

This creates the flow mentioned previously. If the commit was to the default
branch, Docker Scout generates a CVE report. If the commit was to a different
branch, Docker Scout compares the new version to the current published version.
It only shows critical or high-severity vulnerabilities and ignores
vulnerabilities that haven't changed since the last analysis.

Add the following to the YAML file:

```yaml
rules:
  - if: $CI_COMMIT_BRANCH
    exists:
      - Dockerfile
```

These final lines ensure that the pipeline only runs if the commit contains a
Dockerfile and if the commit was to the CI branch.

_The following is a video walkthrough of the process of setting up the workflow with GitLab._

<div style="position: relative; padding-bottom: 64.86486486486486%; height: 0;"><iframe src="https://www.loom.com/embed/451336c4508c42189532108fc37b2560?sid=f912524b-276d-417d-b44a-c2d39719aa1a" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"></iframe></div>
