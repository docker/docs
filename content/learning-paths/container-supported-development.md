---
title: "Faster development and testing with container-supported development"
description: |
  Use containers in your local development loop to develop and test faster… even if your main app isn't running in containers.
summary: |
  Containers don't have to be just for your app. Learn how to run your app's dependent services and other debugging tools to enhance your development environment.
params:
  image: images/learning-paths/container-supported-development.png
  skill: Beginner
  time: TBD
  prereq: None
weight: 50
---

{{< columns >}}

Containers provide a standardized ability to build, share, and run applications. While containers are typically used to containerize your application, they also make it incredibly easy to run essential services needed for development. Instead of installing or connecting to a remote database, you can easily launch your own database. But the possibilities don't stop there.

With container-supported development, you use containers to enhance your development environment by emulating or running your own instances of the services your app needs. This provides faster feedback loops, less coupling with remote services, and a greater ability to test error states.

And best of all, you can have these benefits regardless of whether the main app under development is running in containers.

## Who's this for?

- Teams that want to reduce the coupling they have on shared or deployed infrastructure or remote API endpoints
- Teams that want to reduce the complexity and costs associated with using cloud services directly during development
- Developers that want to make it easier to visualize what's going on in their databases, queues, etc.
- Teams that want to reduce the complexity of setting up their development environment without impacting the development of the app itself


<!-- break -->

## What you'll learn

- The meaning of container-supported development
- How to connect non-containerized applications to containerized services
- Several examples of using containers to emulate or run local instances of services
- How to use containers to add additional troubleshooting and debugging tools to your development environment


## Tools integration

Works well with Docker Compose and Testcontainers

{{< /columns >}}

## Modules

{{< accordion large=true title=`What is container-supported development?` icon=`play_circle` >}}

Container-supported development is the idea of using containers to enhance your development environment by running local instances or emulators of the services your application relies on. Once you're using containers, it's easy to add additional services to visualize or troubleshoot what's going on in your services.

**Duration**: TBD

{{< youtube-embed "8AqKhEO2PQA" >}}

{{< /accordion >}}

{{< accordion large=true title=`Demo: running databases locally` icon=`play_circle` >}}

With container-supported development, it's easy to run databases locally. In this demo, you'll see how to do so, as well as how to connect a non-containerized application to the database.

**Duration**: TBD

{{< youtube-embed "oPGq2AP5OtQ" >}}

{{< /accordion >}}

{{< accordion large=true title=`Demo: mocking remote API endpoints` icon=`play_circle` >}}

Many APIs require data from other data endpoints. In development, this adds complexities such as the sharing of credentials, uptime/availability, and rate limiting. Instead of relying on those services directly, your application can interact with a mock API server.

This demo will demonstrate how using Wiremock can make it easy to develop and test an application, including the APIs various error states.


**Duration**: TBD

{{< youtube-embed "wvLdInoVBGg" >}}

{{< /accordion >}}

{{< accordion large=true title=`Demo: developing the cloud locally` icon=`play_circle` >}}

When developing apps, it's often easier to outsource aspects of the application to cloud services, such as Amazon S3. However, connecting to those services in local development introduces IAM policies, networking constraints, and provisioning complications. While these requirements are important in a production setting, they complicate development environments significantly. 

With container-supported development, you can run local instances of these services during development and testing, removing the need for complex setups. In this demo, you'll see how LocalStack makes it easy to develop and test applications entirely from the developer's workstation.

**Duration**: TBD

{{< youtube-embed "wvLdInoVBGg" >}}

{{< /accordion >}}

{{< accordion large=true title=`Demo: adding additional debug and troubleshooting tools` icon=`play_circle` >}}

Once you start using containers in your development environment, it becomes much easier to add additional containers to visualize the contents of the databases or message queues, seed document stores, or event publishers. In this demo, you'll see a few of these examples, as well as how you can connect multiple containers together to make testing even easier.

**Duration**: TBD

{{< youtube-embed "wvLdInoVBGg" >}}

{{< /accordion >}}

{{< accordion large=true title=`Common challenges and questions` icon=`quiz` >}}

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

{{< /accordion >}}

{{< accordion large=true title=`Resources` icon=`link` >}}

- TBD links to use case guides

{{< /accordion >}}

<div id="lp-survey-anchor"></div>