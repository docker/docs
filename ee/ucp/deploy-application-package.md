---
title: Deploy an application package
description: Learn how to deploy an application package in UCP
keywords: ucp, swarm, kubernetes, application, app package
---

>{% include enterprise_label_shortform.md %}

## Application packages

Docker Enterprise Edition 2.1 introduces application packages in Docker. With application packages, you can add metadata and settings to an existing Compose file. This gives operators more context about applications they deploy and manage.

Application packages can present in one of two different formats, **Directory** or **Single-file**:

- **Directory**: Defined by metadata.yml, a docker-compose.yml, and a settings.yml files inside a `my-app.dockerapp` folder. This is also called the folder format.
- **Single-file**: Defined by metadata.yml, docker-compose.yml, and settings.yml concatenated in that order and separated by `---\n` in a single file named named `my-app.dockerapp`.

Once an application package has been deployed, you manipulate and manage it as you would any stack.

## Creating a stack in the UCP web interface

1. Access the UCP web interface.
2. From the left menu, select **Shared Resources** > **Stacks**.

    ![Create stacks in UCP](/ee/ucp/images/v32stacks.png)

3. Click **Create Stack** to open the **Create Application** window. The **1. Configure Application** section will become active.

    ![Configure stacks in UCP](/ee/ucp/images/ucp-config-stack.png)

4. Enter a name for the stack in the **Name** field.
5. Click to indicate the **Orchestrator Mode**, either **Swarm Services** or **Kubernetes Workloads**. Note that if you select Kubernetes Workloads, the **Namespace** drop-down list will display, from which you must select one of the namespaces offered. 

    ![Specify namespace for a stack in UCP](/ee/ucp/images/ucp-stack-namespace.png)

6. Click to indicate the **Application File Mode**, either **Compose File** or **App Package**.
7. Click **Next** to open the **2. Add Application File** section.
8. Add the application file, according to the Application File Mode selected in section 1. 
    - **Compose File:** Enter or upload the docker-compose.yml file.
    - **App Package:** Enter or upload the application package in the single-file format. 

    ![Provide docker-compose.yml in UCP](/ee/ucp/images/ucp-stack-compose.png)

    ![Provide application package in UCP](/ee/ucp/images/ucp-stack-package.png)

9. Select **Create**.

## Single-file format example

The following file is an example of a single-file application package.

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
