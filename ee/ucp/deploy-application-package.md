---
title: Deploy an application package
description: Learn how to deploy an application package in UCP
keywords: ucp, swarm, kubernetes, application, app package
---

## Application packages

Docker Enterprise Edition 2.1 introduces application packages in Docker. With application packages, you can add metadata and settings to an existing Compose file. This gives operators more context about applications they deploy and manage.

An application package can have one of these formats:

- **Directory format**: Defined by metadata.yml, a docker-compose.yml, and a settings.yml files inside a `my-app.dockerapp` folder. This is also called the folder format.
- **Single-file format**: Defined by metadata.yml, docker-compose.yml, and settings.yml concatenated in that order and separated by `---\n` in a single file named named `my-app.dockerapp`.

Once an application package has been deployed, you manipulate and manage it as you would any stack.

## Creating a stack in the UCP web interface

To create a stack in the UCP web interface, follow these steps:

1. Go to the UCP web interface.
2. In the lefthand menu, first select **Shared Resources**, then **Stacks**.

    ![Create stacks in UCP](/ee/ucp/images/ucp-create-stack.png)

3. Select **Create Stack** to display **1. Configure Application** in the stack creation dialog.

    ![Configure stacks in UCP](/ee/ucp/images/ucp-config-stack.png)

4. Enter a name for the stack in the **Name** field.
5. Select either **Swarm Services** or **Kubernetes Workloads** for the orchestrator mode. If you select Kubernetes, also select a namespace in the **Namespace** drop-down list.

    ![Specify namespace for a stack in UCP](/ee/ucp/images/ucp-stack-namespace.png)

6. Select either **Compose File** or **App Package** for the **Application File Mode**.
7. Select **Next**.
8. If you selected Compose file, enter or upload your `docker-compose.yml` in **2. Add Application File**.

    ![Provide docker-compose.yml in UCP](/ee/ucp/images/ucp-stack-compose.png)

   or if you selected **App Package**, enter or upload the application package in the single-file format.

    ![Provide application package in UCP](/ee/ucp/images/ucp-stack-package.png)

9. Select **Create**.

## Single-file format example

Here is an example of a single-file application package:

```
version: 0.1.0
name: hello-world
description: "Hello, World!"
namespace: myHubUsername
maintainers:
  - name: user
    email: "user@email.com"
---
version: "3.6"
services:
  hello:
    image: hashicorp/http-echo
    command: ["-text", "${text}"]
    ports:
      - ${port}:5678

---
port: 8080
text: Hello, World!
```
