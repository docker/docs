---
description: Discover how to browse and search Docker Hub's extensive resources.
keywords: Docker Hub, Hub, explore, search, image library
title: Docker Hub search
linkTitle: Search
weight: 10
---

The [Docker Hub search interface](https://hub.docker.com/search) lets you
explore millions of resources. To help you find exactly what you need, it offers
a variety of filters that let you narrow your results or discover different
types of content.

## Filters

The search functionality includes filters to narrow down
results based on your requirements, such as products, categories, trusted
content, and publishers. This ensures that you can quickly find and access the resources best
suited to your project.

### Products

Docker Hub's content library features various products, each designed to meet
specific needs of developers and organizations. These products include:

- Images
- Extensions
- Helm charts
- Compose
- AI models
- Plugins

#### Images

Docker Hub hosts millions of container images, making it the go-to repository
for containerized applications and solutions. These images include:

- Operating system images: Foundational images for Linux distributions like
  Ubuntu, Debian, and Alpine, or Windows Server images.
- Database and storage images: Pre-configured databases such as MySQL,
  PostgreSQL, and MongoDB to simplify application development.
- Languages and frameworks-based images: Popular images for Java, Python,
  Node.js, Ruby, .NET, and more, offering pre-built environments for faster
  development.

Images in Docker Hub simplify the development process by providing pre-built,
reusable building blocks, reducing the need to start from scratch. Whether
you're a beginner building your first container or an enterprise managing
complex architectures, Docker Hub images provide a reliable foundation.

#### Extensions

Docker Hub offers extensions for Docker Desktop, which enhance its core
functionality. These extensions are purpose-built to streamline the software
development lifecycle. Extensions provide tools for:

- System optimization and monitoring: Manage resources and optimize Docker
  Desktop’s performance.
- Container management: Simplify container deployment and monitoring.
- Database management: Facilitate efficient database operations within
  containers.
- Kubernetes and cloud integration: Bridge local environments with cloud-native
  and Kubernetes workflows.
- Visualization tools: Gain insights into container resource usage through
  graphical representations.

Extensions help developers and teams create a more efficient and unified
workflow by reducing context switching and bringing essential tools into Docker
Desktop's interface.

To learn more about extensions, see [Docker
Extensions](/manuals/extensions/_index.md).

#### Helm charts

Helm charts in Docker Hub provide a streamlined way to package, configure, and
deploy Kubernetes applications. Helm is the package manager for Kubernetes, and
charts are pre-configured templates that define the resources needed to run an
application in a Kubernetes cluster. Docker Hub hosts a variety of Helm charts
that provide:

- Application packaging: Bundle Kubernetes manifests, configurations, and
  dependencies into a single, reusable chart.
- Version management: Track and manage different versions of your application
  deployments.
- Configuration templating: Customize deployments with values files, making it
  easy to deploy the same application across different environments.
- Dependency management: Automatically handle chart dependencies, ensuring all
  required components are deployed together.

Helm charts reduce the complexity of Kubernetes deployments, making it easier
for teams to deploy, upgrade, and manage applications in production
environments.

#### Compose

Docker Compose files in Docker Hub enable multi-container application
orchestration through simple YAML configuration files. Compose is a tool for
defining and running multi-container Docker applications, and Docker Hub hosts
Compose files that help you:

- Multi-container orchestration: Define and run applications consisting of
  multiple interconnected containers with a single command.
- Service configuration: Specify container images, environment variables,
  networks, volumes, and dependencies in a declarative format.
- Development and testing: Quickly spin up complete application stacks for local
  development, testing, or demonstration purposes.
- Environment consistency: Ensure consistent application behavior across
  development, staging, and production environments.

Compose files simplify the process of managing complex applications by
providing a clear, version-controlled definition of your entire application
stack.

#### AI models

Docker Hub hosts AI and machine learning models in containerized formats,
making it easier to deploy, share, and run AI applications across different
environments. These containerized AI models provide:

- Pre-trained models: Access ready-to-use machine learning models for common
  tasks such as image recognition, natural language processing, and predictive
  analytics.
- Model serving: Deploy models as containerized services that can be easily
  integrated into applications and scaled as needed.
- Reproducible environments: Package models with their dependencies, ensuring
  consistent behavior across development and production environments.
- Framework support: Find models built with popular frameworks like TensorFlow,
  PyTorch, scikit-learn, and others.

Containerized AI models remove the need to manage infrastructure dependencies,
let you deploy models across different environments, and enable scaling as
application demands change.

#### Plugins

Plugins in Docker Hub let you extend and customize Docker Engine to suit
specialized requirements. Plugins integrate directly with the Docker Engine and
provide capabilities such as:

- Network plugins: Enhance networking functionality, enabling integration with
  complex network infrastructures.
- Volume plugins: Provide advanced storage options, supporting persistent and
  distributed storage across various backends.
- Authorization plugins: Offer fine-grained access control to secure Docker
  environments.

By leveraging Docker plugins, teams can tailor Docker Engine to meet their
specific operational needs, ensuring compatibility with existing infrastructures
and workflows.

To learn more about plugins, see [Docker Engine managed plugin
system](/manuals/engine/extend/_index.md).


### Trusted content

Docker Hub's trusted content provides a curated selection of high-quality,
secure images designed to give developers confidence in the reliability and
security of the resources they use. These images are stable, regularly updated,
and adhere to industry best practices, making them a strong foundation for
building and deploying applications. Docker Hub's trusted content includes,
Docker Hardened Images, Docker Official Images, Verified Publisher images, and
Docker-Sponsored Open Source Software images.

For more details, see [Trusted content](./trusted-content.md).

### Publishers

The **Publishers** filter lets you narrow image results by the organization
that published the image.

Publishers of trusted content appear first. Only a select number of publishers
are shown.

### Categories

Docker Hub makes it easy to find and explore container images with categories.
Categories group images based on their primary use case, helping you quickly
locate the tools and resources you need to build, deploy, and run your
applications.

{{% include "hub-categories.md" %}}

### Operating systems

The **Operating systems** filter lets you narrow your search to container
images compatible with specific host operating systems. This filter ensures that
the images you use align with your target environment, whether you're developing
for Linux-based systems, Windows, or both.

- **Linux**: Access a wide range of images tailored for Linux environments.
  These images provide foundational environments for building and running
  Linux-based applications in containers.
- **Windows**: Explore Windows container images.

> [!NOTE]
>
> The **Operating systems** filter is only available for images. If you select
> the **Extensions** or **Plugins** filter, then the **Operating systems**
> filter isn't available.

### Architectures

The **Architectures** filter lets you find images built to support specific CPU
architectures. This ensures compatibility with your hardware environment, from
development machines to production servers.

- **ARM**: Select images compatible with ARM processors, commonly used in IoT
  devices and embedded systems.
- **ARM 64**: Locate 64-bit ARM-compatible images for modern ARM processors,
  such as those in AWS Graviton or Apple Silicon.
- **IBM POWER**: Find images optimized for IBM Power Systems, offering
  performance and reliability for enterprise workloads.
- **PowerPC 64 LE**: Access images designed for the little-endian PowerPC 64-bit
  architecture.
- **IBM Z**: Discover images tailored for IBM Z mainframes, ensuring
  compatibility with enterprise-grade hardware.
- **x86**: Choose images compatible with 32-bit x86 architectures, suitable for
  older systems or lightweight environments.
- **x86-64**: Filter images for modern 64-bit x86 systems, widely used in
  desktops, servers, and cloud infrastructures.

> [!NOTE]
>
> The **Architectures** filter is only available for images. If you select the
> **Extensions** or **Plugins** filter, then the **Architectures** filter isn't
> available.

### Reviewed by Docker

The **Reviewed by Docker** filter provides an extra layer of assurance when
selecting extensions. This filter helps you identify whether a Docker Desktop
extension has been reviewed by Docker for quality and reliability.

- **Reviewed**: Extensions that have undergone Docker's review process, ensuring
  they meet high standards.
- **Not Reviewed**: Extensions that have not been reviewed by Docker.

> [!NOTE]
>
> The **Reviewed by Docker** filter is only available for extensions. To make
> the filter available, you must select only the **Extensions** filter in **Products**.
