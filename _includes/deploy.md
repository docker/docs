<!-- This text will be included as part of the Deploy your app section in the Language-specific guides -->

Now, that we have configured a CI/CD pipleline, let's look at how we can deploy the application. Docker supports deploying containers on Azure ACI and AWS ECS. You can also deploy your application to Kubernetes if you have enabled Kubernetes in Docker Desktop.

## Docker and Azure ACI

The Docker Azure Integration enables developers to use native Docker commands to run applications in Azure Container Instances (ACI) when building cloud-native applications. The new experience provides a tight integration between Docker Desktop and Microsoft Azure allowing developers to quickly run applications using the Docker CLI or VS Code extension, to switch seamlessly from local development to cloud deployment.

For detailed instructions, see [Deploying Docker containers on Azure](/cloud/aci-integration/).

## Docker and AWS ECS

The Docker ECS Integration enables developers to use native Docker commands in Docker Compose CLI to run applications in Amazon EC2 Container Service (ECS) when building cloud-native applications.

The integration between Docker and Amazon ECS allows developers to use the Docker Compose CLI to set up an AWS context in one Docker command, allowing you to switch from a local context to a cloud context and run applications quickly and easily simplify multi-container application development on Amazon ECS using Compose files.

For detailed instructions, see [Deploying Docker containers on ECS](/cloud/ecs-integration/).

## Kubernetes

Docker Desktop includes a standalone Kubernetes server and client, as well as Docker CLI integration that runs on your machine. Enabling Kubernetes allows you to test your workloads on Kubernetes.

To enable Kubernetes:

1. From the Docker menu, select **Preferences** (**Settings** on Windows).
2. Select **Kubernetes** and then click **Enable Kubernetes**.

    This starts a Kubernetes single-node cluster when Docker Desktop starts.

For detailed information, see [Deploy on Kubernetes](/desktop/kubernetes/) and [Describing apps using Kubernetes YAML](/get-started/kube-deploy/#describing-apps-using-kubernetes-yaml).
