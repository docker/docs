---
title: Docker Template
description: Working with Docker Template
keywords: Docker, application template, Application Designer,
---

>This is an experimental feature.
>
>{% include experimental.md %}

## Overview

Docker Template is a CLI plugin that introduces a top-level `docker template`
command that allows users to create new Docker applications by using a library
of templates. There are two types of templates — service templates and
application templates.

A _service template_ is a container image that generates code and contains the
metadata associated with the image.

- The container image takes `/run/configuration` mounted file as input to
  generate assets such as code, Dockerfile, and `docker-compose.yaml` for a
given service, and writes the output to the `/project` mounted folder.

- The metadata file that describes the service template is called the service
  definition.  It contains the name of the service, description, and available
parameters such as ports, volumes, etc. For a complete list of parameters that
are allowed, see [Docker Template API
reference](/app-template/api-reference/).

An _application template_ is a collection of one or more service templates. An
application template generates a Dockerfile per service and only one Compose
file for the entire application, aggregating all services.

## Create a custom service template

A Docker template contains a predefined set of service and application
templates. To create a custom template based on your requirements, you must
complete the following steps:

1. Create a service container image
2. Create the service template definition
3. Add the service template to the library
4. Share the service template

### Create a service container image

A service template provides the description required by Docker Template to
scaffold a project. A service template runs inside a container with two bind
mounts:

1. `/run/configuration`, a JSON file which contains all settings such as
   parameters, image name, etc. For example:

    ```json
      {
      "parameters": {
        "externalPort": "80",
        "artifactId": "com.company.app"
      },
      ...
    }
    ```

2. `/project`, the output folder to which the container image writes the generated assets.

#### Basic service template

Services that generate a template using code must contain the following files
that are valid:

- A *Dockerfile* located at the root of the `my-service` folder. This is the
  Dockerfile that is used for the service when running the application.

- A *docker-compose.yaml* file  located at the root of the `my-service` folder.
  The `docker-compose.yaml` file must contain the service declaration and any
optional volumes or secrets.

Here’s an example of a simple NodeJS service:

```bash
my-service
├── Dockerfile    # The Dockerfile of the service template
└── assets
    ├── Dockerfile           # The Dockerfile of the generated service
    └── docker-compose.yaml  # The service declaration
```

The NodeJS service contains the following files:

`my-service/Dockerfile`

```conf
FROM alpine
COPY assets /assets
CMD ["cp", "/assets", "/project"]
FROM dockertemplate/interpolator:v0.1.5 as interpolator
COPY assets /assets
```

`my-service/assets/docker-compose.yaml`

{% raw %}
```yaml
version: "3.6"
services:
  {{ .Name }}:
    build: {{ .Name }}
    ports:
      - {{ .Parameters.externalPort }}:3000
```
{% endraw %}

`my-service/assets/Dockerfile`

```conf
FROM NODE:9
WORKDIR /app
COPY package.json .
RUN yarn install
COPY . .
CMD ["yarn", "run", "start"]
```

> **Note:** After scaffolding the template, you can add the default files your
> template contains to the `assets` folder.

The next step is to build and push the service template image to a remote
repository by running the following command:

```bash
cd [...]/my-service
docker build -t org/my-service .
docker push org/my-service
```

### Create the service template definition

The service definition contains metadata that describes a service template. It
contains the name of the service, description, and available parameters such as
ports, volumes, etc.  After creating the service definition, you can proceed to
[Add templates to Docker Template](#add-templates-to-docker-template) to add
the service definition to the Docker Template repository.

Of all the available service and application definitions, Docker Template has
access to only one catalog, referred to as the ‘repository’. It uses the
catalog content to display service and application templates to the end user.

Here is an example of the Express service definition:

```yaml
- apiVersion: v1alpha1 # constant
  kind: ServiceTemplate  # constant
  metadata:
    name: Express # the name of the service
    platforms:
    - linux
  spec:
    title: Express    # The title/label of the service
    icon: https://docker-application-template.s3.amazonaws.com/assets/express.png # url for an icon
    description: NodeJS web application with Express server
    source:
      image: org/my-service:latest
```

The most important section here is `image: org/my-service:latest`. This is the
image associated with this service template. You can use this line to point to
any image. For example, you can use an Express image directly from the hub
`docker.io/dockertemplate/express:latest` or from the DTR private repository
`myrepo/my-service:latest`. The other properties in the service definition are
mostly metadata for display and indexation purposes.

#### Adding parameters to the service

Now that you have created a simple express service, you can customize it based
on your requirements. For example, you can choose the version of NodeJS to use
when running the service.

To customize a service, you need to complete the following tasks:

1. Declare the parameters in the service definition. This tells Docker Template
   whether or not the CLI can accept the parameters, and allows the
   [Application Designer](/ee/desktop/app-designer) to be aware of the new
   options.

2. Use the parameters during service construction.

#### Declare the parameters

Add the parameters available to the application. The following example adds the
NodeJS version and the external port:

```yaml
- [...]
  spec:
    [...]
    parameters:
    - name: node
      defaultValue: "9"
      description: Node version
      type: enum
      values:
      - value: "10"
        description: "10"
      - value: "9"
        description: "9"
      - value: "8"
        description: "8"
    - defaultValue: "3000"
      description: External port
      name: externalPort
      type: hostPort
    [...]
```

#### Use the parameters during service construction

When you run the service template container, a volume is mounted making the
service parameters available at `/run/configuration`.

The file matches the following go struct:

```golang
type TemplateContext struct {
   ServiceID string            `json:"serviceId,omitempty"`
   Name      string            `json:"name,omitempty"`
   Parameters   map[string]string `json:"parameters,omitempty"`

   TargetPath string `json:"targetPath,omitempty"`
   Namespace string `json:"namespace,omitempty"`

   Services []ConfiguredService `json:"services,omitempty"`
}
```

Where `ConfiguredService` is:

```go
type ConfiguredService struct {
	ID      string            `json:"serviceId,omitempty"`
	Name    string            `json:"name,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}
```

You can then use the file to obtain values for the parameters and use this
information based on your requirements. However, in most cases, the JSON file
is used to interpolate the variables. Therefore, we provide a utility called
`interpolator` that expands variables in templates. For more information, see
[Interpolator](#interpolator).

To use the `interpolator` image, update `my-service/Dockerfile` to use the
following Dockerfile:

```conf
FROM dockertemplate/interpolator:v0.1.5
COPY assets .
```

> **Note:** The interpolator tag must match the version used in Docker
> Template. Verify this using the `docker template version` command .

This places the  interpolator image in the `/assets` folder and copies the
folder to the target `/project` folder. If you prefer to do this manually, use
a Dockerfile instead:

```conf
WORKDIR /assets
CMD ["/interpolator", "-config", "/run/configuration", "-source", "/assets", "-destination", "/project"]
```

When this is complete, use the newly added node option in
`my-service/assets/Dockerfile`, by replacing the line:

`FROM node:9`

with

{% raw %}`FROM node:{{ .Parameters.node }}`{% endraw %}

Now, build and push the image to your repository.

### Add service template to the library

You must add the service to a repository file in order to see it when you run
the `docker template ls` command, or to make the service available in the
Application Designer.

#### Create the repository file

Create a local repository file called `library.yaml` anywhere on your local
drive and add the newly created service definitions and application definitions
to it.

`library.yaml`

```yaml
apiVersion: v1alpha1
kind: RepositoryContent
services: # List of service templates available
- apiVersion: v1alpha1 # here is the service definition for our service template.
  kind: ServiceTemplate
  name: express
  spec:
    title: Express
    [...]
```

#### Add the local repository to docker-template settings

> **Note:** You can also use the instructions in this section to add templates
> to the [Application Designer](/ee/desktop/app-designer).

Now that you have created a local repository and added service definitions to
it, you must make Docker Template aware of these. To do this:

1. Edit `~/.docker/application-template/preferences.yaml` as follows:

   ```yaml
   apiVersion: v1alpha1
   channel: master
   kind: Preferences
   repositories:
   - name: library-master
     url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
   ```

2. Add your local repository:

> **Note:** Do not remove or comment out the default library `library-master`.
> This library contain template plugins that are required to build all Docker
> Templates.

   ```yaml
   apiVersion: v1alpha1
   channel: master
   kind: Preferences
   repositories:
   - name: custom-services
     url: file:///path/to/my/library.yaml
   - name: library-master
     url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
   ```

When configuring a local repository on Windows, the `url` structure is slightly
different:

```yaml
- name: custom-services
  url: file://c:/path/to/my/library.yaml
```

After updating the `preferences.yaml` file, run `docker template ls` or restart
the Application Designer and select **Custom application**. The new service
should now be visible in the list of available services.

### Share custom service templates

To share a custom service template, you must complete the following steps:

1. Push the image to an available endpoint (for example, Docker Hub)

2. Share the service definition (for example, GitHub)

3. Ensure the receiver has modified their `preferences.yaml` file to point to
   the service definition that you have shared, and are permitted to accept
   remote images.

## Create a custom application template

An application template is a collection of one or more service templates. You
must complete the following steps to create a custom application template:

1. Create an application template definition

2. Add the application template to the library

3. Share your custom application template

### Create the application definition

An application template definition contains metadata that describes an
application template. It contains information such as the name and description
of the template, the services it contains, and  the parameters for each of the
services.

Before you create an application template definition, you must create a
repository that contains the services you are planning to include in the
template. For more information, see [Create the repository
file](#create-the-repository-file).

For example, to create an Express and MySQL application, the application
definition must be similar to the following yaml file:

```yaml
apiVersion: v1alpha1  #constant
kind: ApplicationTemplate  #constant
metadata:
  name: express-mysql #the name of the application
  platforms:
  - linux
spec:
  description: Sample application with a NodeJS backend and a MySQL database
  services: # list of the services
  - name: back
    serviceId: express # service name
    parameters:  # (optional) define the default application parameters
      externalPort: 9000
  - name: db
    serviceId: mysql
  title: Express / MySQL application
```

### Add the template to the library

Create a local repository file called `library.yaml` anywhere on your local
drive. If you have already created the `library.yaml` file, add the application
definitions to it.

`library.yaml`

```yaml
apiVersion: v1alpha1
kind: RepositoryContent
services: # List of service templates available
- apiVersion: v1alpha1 # here is the service definition for our service template.
  kind: ServiceTemplate
  name: express
  spec:
    title: Express
    [...]
templates: # List of application templates available
- apiVersion: v1alpha1  #constant
  kind: ApplicationTemplate # here is the application definition for our application template
  metadata:
    name: express-mysql
  spec:
```

### Add the local repository to `docker-template` settings

Now that you have created a local repository and added application definitions,
you must make Docker Template aware of these. To do this:

1. Edit `~/.docker/application-template/preferences.yaml` as follows:

   ```yaml
   apiVersion: v1alpha1
   channel: master
   kind: Preferences
   repositories:
   - name: library-master
     url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
   ```

2. Add your local repository:

> **Note:** Do not remove or comment out the default library `library-master`.
> This library contain template plugins that are required to build all Docker
> Templates.

   ```yaml
   apiVersion: v1alpha1
   channel: master
   kind: Preferences
   repositories:
   - name: custom-services
     url: file:///path/to/my/library.yaml
   - name: library-master
     url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
   ```

When configuring a local repository on Windows, the `url` structure is slightly
different:

```yaml
- name: custom-services
  url: file://c:/path/to/my/library.yaml
```

After updating the `preferences.yaml` file, run `docker template ls` or restart
the Application Designer and select **Custom application**. The new template
should now be visible in the list of available templates.

### Share the custom application template

To share a custom application template, you must complete the following steps:

1. Push the image to an available endpoint (for example, Docker Hub)

2. Share the application definition (for example, GitHub)

3. Ensure the receiver has modified their `preferences.yaml` file to point to
   the application definition that you have shared, and are permitted to accept
   remote images.

## Interpolator

The `interpolator` utility is basically an image containing a binary which:

- takes a folder (assets folder) and the service parameter file as input,
- replaces variables in the input folder using the parameters specified by the
  user (for example, the service name, external port, etc), and
- writes the interpolated files to the destination folder.

The interpolator implementation uses [Golang
template](https://golang.org/pkg/text/template/) to aggregate the services to
create the final application. If your service template uses the `interpolator`
image by default, it expects all the asset files to be located in the `/assets`
folder:

`/interpolator -source /assets -destination /project`

However, you can create your own scaffolding script that performs calls to the
`interpolator`.

> **Note:** It is not mandatory to use the `interpolator` utility. You can use
> a utility of your choice to handle parameter replacement and file copying to
> achieve the same result.

The following table lists the `interpolator` binary options:

 | Parameter        | Default value        | Description                                                   |
 | :----------------|:---------------------|:--------------------------------------------------------------|
 | `-source`        | none                 |  Source file or folder to interpolate from                    |
 | `-destination`   | none                 |  Destination file or folder to copy the interpolated files to |
 | `-config`        | `/run/configuration` |  The path to the json configuration file                      |
 | `-skip-template` | false                | If set to `true`, it copies assets without any transformation |
