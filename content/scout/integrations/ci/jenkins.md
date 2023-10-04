---
description: How to integrate Docker Scout with Jenkins
keywords: supply chain, security, ci, continuous integration, jenkins
title: Integrate Docker Scout with Jenkins
---

You can add the following stage and steps definition to a `Jenkinsfile` to run
Docker Scout as part of a Jenkins pipeline. The pipeline needs two secrets
defined to authenticate with Docker Hub: `DOCKER_HUB_USER` and `DOCKER_HUB_PAT`
It also needs an environment variable defined for the image and tag.

```groovy
…
stage('Analyze image') {
    steps {
        // Install Docker Scout
        sh 'curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s -- -b /usr/local/bin'
        
        // Log into Docker Hub
        sh 'echo $DOCKER_HUB_PAT | docker login -u $DOCKER_HUB_USER --password-stdin'

        // Analyze and fail on critical or high vulnerabilities
        sh 'docker-scout cves $IMAGE_TAG --exit-code --only-severity critical,high'
    }
}
```

This installs Docker Scout, logs into Docker Hub, and then runs Docker Scout to
generate a CVE report for an image and tag. It only shows critical or
high-severity vulnerabilities.
