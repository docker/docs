{% if include.product == "admin" %}
  {% assign product_link="[Docker Admin](https://admin.docker.com)" %}
  {% assign iam_navigation="Select your organization in the left navigation drop-down menu, and then select **Image Access**." %}
{% else %}
  {% assign product_link="[Docker Hub](https://hub.docker.com)" %}
  {% assign iam_navigation="Select **Organizations**, your organization, **Settings**, and then select **Image Access**." %}
{% endif %}

>Note
>
>Image Access Management is available to [Docker Business](/subscription/details/) customers only.

Image Access Management gives administrators control over which types of images, such as Docker Official Images, Docker Verified Publisher Images, or community images, their developers can pull from Docker Hub.

For example, a developer, who is part of an organization, building a new containerized application could accidentally use an untrusted, community image as a component of their application. This image could be malicious and pose a security risk to the company. Using Image Access Management, the organization owner can ensure that the developer can only access trusted content like Docker Official Images, Docker Verified Publisher Images, or the organization’s own images, preventing such a risk.

## Prerequisites

You need to [configure a registry.json to enforce sign-in](/docker-hub/configure-sign-in/). For Image Access Management to take effect, Docker Desktop users must authenticate to your organization.

## Configure Image Access Management permissions

1. Sign in to {{ product_link }}{: target="_blank" rel="noopener" class="_"}.
2. {{ iam_navigation }}
3. Enable Image Access Management to set the permissions for the following categories of images you can manage:
- **Organization Images**: When Image Access Management is enabled, images from your organization are always allowed. These images can be public or private created by members within your organization.
- **Docker Official Images**: A curated set of Docker repositories hosted on Hub. They provide OS repositories, best practices for Dockerfiles, drop-in solutions, and applies security updates on time.
- **Docker Verified Publisher Images**: published by Docker partners that are part of the Verified Publisher program and are qualified to be included in the developer secure supply chain. You can set permissions to **Allowed** or **Restricted**.
- **Community Images**: Images are always disabled when Image Access Management is enabled. These images are not trusted because various Docker Hub users contribute them and pose security risks.

    > **Note**
    >
    > Image Access Management is turned off by default. However, owners in your organization have access to all images regardless of the settings.

4. Select the category restrictions for your images by selecting **Allowed**.
     Once the restrictions are applied, your members can view the organization permissions page in a read-only format.

## Verify the restrictions

The new Image Access Management policy takes effect after the developer successfully authenticates to Docker Desktop using their organization credentials. If a developer attempts to pull a disallowed image type using Docker, they receive an error message.