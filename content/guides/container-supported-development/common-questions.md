---
title: Common challenges and questions
description: Explore common challenges and questions related to Docker Compose.
weight: 60
---

<!-- vale Docker.HeadingLength = NO -->

### Can I use containers to run services locally even if my main application isn't containerized?

Yes! But there's one extra step required - exposing the required ports. Once the required ports are exposed on the host, the non-containerized application can connect to the service running in the container.

### What benefits do I gain by running services in containers over other methods (natively installed, VMs, etc.)?

Containers provide faster startup times and lower resource consumption compared to VMs. By using containers, you have the ability to package and distribute the services and environment your app needs to run, ensuring everyone on your team can get up and going quickly without spending time installing and troubleshooting native installation issues.

### Should I worry about emulating services instead of using the real thing during development?

For many services, such as databases and caches, your application is connecting to the same thing. As an example, your app will use a MySQL database running in a container the same as a managed MySQL database. Therefore, there's no need to worry here.

When using emulated services, there is always a risk that behavioral differences might occur. Due diligence is required any time you use an emulating service to ensure it has proper tests and is being kept up-to-date with the original service's behaviors.

### How does running services in containers speed up development and testing?

Containers allow developers to quickly spin up and tear down services, enabling rapid testing and iteration without the overhead of managing dependencies or system configurations manually. By running services locally, you are able to decouple your application from remote and potentially shared resources.

### How can I run a completely customized service in my development environment?

Containers allow you to run any service you'd like. To run customized services, you can build and publish your own images and then use them in your development environment.

<div id="container-supported-development-lp-survey-anchor"></div>
