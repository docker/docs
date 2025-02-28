---
description: Key benefits and use cases of Testcontainers OSS
keywords: documentation, docs, docker, testcontainers, containers, uses,  benefits
title: Why use Testcontainers?
weight: 10
linkTitle: Why use Testcontainers?
aliases: 
- /testcontainers/benefits/
---

### Benefits of using Testcontainers

* **On-demand isolated infrastructure provisioning:**
You don't need to have a pre-provisioned integration testing infrastructure. 
Testcontainers will provide the required services before running your tests. 
There will be no test data pollution, even when multiple build pipelines run in parallel 
because each pipeline runs with an isolated set of services.

* **Consistent experience on both local and CI environments:**
You can run your integration tests right from your IDE, just like you run unit tests. 
No need to push your changes and wait for the CI pipeline to complete.

* **Reliable test setup using wait strategies:**
Docker containers need to be started and fully initialized before using them in your tests. 
The Testcontainers library offers several out-of-the-box wait strategies implementations to make sure 
the containers (and the application within) are fully initialized. 
Testcontainers modules already implement the relevant wait strategies for a given technology, 
and you can always implement your own or create a composite strategy if needed.

* **Advanced networking capabilities:**
Testcontainers libraries map the container's ports onto random ports available on the host machine 
so that your tests connect reliably to those services. You can even create a (Docker) network and 
connect multiple containers together so that they talk to each other via static Docker network aliases.

* **Automatic clean up:**
The Testcontainers library takes care of removing any created resources (containers, volumes, networks etc.) 
automatically after the test execution is complete by using the Ryuk sidecar container. 
While starting the required containers, Testcontainers attaches a set of labels to the 
created resources (containers, volumes, networks etc) and Ryuk automatically performs 
resource clean up by matching those labels. 
This works reliably even when the test process exits abnormally (for example sending a SIGKILL).
