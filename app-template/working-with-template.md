---
title: Working with Docker Template
description: Working with Docker Application Template
keywords: Docker, application template, Application Designer,
---

## Overview

Docker Template is a CLI plugin that introduces a top-level `docker template` command that allows users to create new Docker applications by using a library of templates. There are two types of templates — service templates and application templates.

- A _service template_ is a container that holds the assets such as code, Dockerfile, docker-compose.yaml to generate an application.
- An _application template_ is a collection of one or more service templates.

A Docker template contains a predefined set of service and application templates. However, you can create **custom templates** by completing the following steps:

1. Create a service template
2. Create the service template definition
3. Create the application template definition
4. Add the template to Docker Template settings

## Create a service template

A service template is a container that generates all the assets such as code, Dockerfile, docker-compose.yaml for a given service.

A service template provides the description required by Docker Template to scaffold a project. A service template runs inside a container with two bind mounts:

1. `/run/configuration`, a JSON file which contains all settings such as parameters, image name, etc.

```
  {
  "parameters": {
    "externalPort": "80",
    "artifactId": "com.company.app"
  },
  ...
}
```

2. `/project`, where the container will write assets to

Docker Template converts the `docker-compose.yaml` file into a template using a Go template and aggregates the services to create the final application.

### Create a basic service template

To create a basic service template, you need to create two files — a dockerfile and a docker compose file in a new folder. For example, to create a new MySQL service template, create the following files in a folder called  `my-service`:

`docker-compose.yaml`

```
version: "3.6"
services:
 {{ .Name }}:
   image: mysql
```

`Dockerfile`

```
FROM alpine
COPY docker-compose.yaml .
CMD cp docker-compose.yaml /project/
```

This adds a MySQL service to your application.

### Create a service with code

Services that generate a template using code must contain the following files that are valid:

- A *Dockerfile* located at the root of the `my-service` folder. This is the Dockerfile that is used for the service when running the application.

- A *docker-compose.yaml* file  located at the root of the `my-service` folder. The `docker-compose.yaml` file must contain the service declaration and any  optional volumes / secrets.

Here’s an example of a simple NodeJS service:

```
my-service
├── Dockerfile    # The Dockerfile of the service template
└── assets
    ├── Dockerfile           # The Dockerfile of the generated service
    └── docker-compose.yaml  # The service declaration
```

The NodeJS service contains the following contents:

`my-service/Dockerfile`

```
FROM alpine
COPY assets /assets
CMD ["cp", "/assets", "/project"]
```

`my-service/assets/docker-compose.yaml`

```
version: "3.6"
services:
  {{ .Name }}:
    build: {{ .Name }}
    ports:
      - {{ .Parameters.externalPort }}:3000
```

`my-service/assets/Dockerfile`

```
FROM NODE:9
WORKDIR /app
COPY package.json .
RUN yarn install
COPY . .
CMD ["yarn", "run", "start"]
```

> **Note:** After scaffolding the template, you can add the default files your template contains to the `assets` folder.

#### Build and push the template image

The next step is to build and push the service template image to a remote repository by running the following command:

```
cd [...]/my-service
docker build -t org/my-service .
docker push org/my-service
```

To build and push the image to an instance of Docker Trusted Registry(DTR), specify a repository name:

```
cd [...]/my-service
docker build -t myrepo:5000/my-service .
docker push myrepo:5000/my-service
```

## Create the service template definition

The service definition contains metadata that describes a service template. It contains the name of the service, description, and available parameters such as ports, volumes, etc.
After creating the service definition, you can proceed to [Add templates to Docker Template](#add-templates-to-docker-template) to add the service definition to the Docker Template repository.

Of all the available service and application definitions, Docker Template has access to only one catalog, referred to as the ‘repository’. It uses the catalog content to display service and application templates to the end user.

Here is an example of the  Express service definition:

```
- apiVersion: v1alpha1 # constant
  kind: ServiceTemplate  # constant
  metadata:
    name: Express         # the name of the service
  spec:
    title: Express    # The title/label of the service
    icon: https://docker-application-template.s3.amazonaws.com/assets/express.png # url for an icon
    description: NodeJS web application with Express server
    source:
      image: org/my-service:latest
```

The most important section here is `image: org/my-service:latest`. This is the image associated with this service template. You can use this line to point to any image. For example, you can use an Express image directly from the hub `docker.io/dockertemplate/express:latest` or from the DTR private repository `myrepo:5000/my-service:latest`. The other properties in the service definition are mostly metadata for display and indexation purposes.

### Adding parameters to the service

Now that you have created a simple express service, you can customize it based on your requirements. For example, you can choose the version of NodeJS to use when running the service.

To customize a service, you need to complete the following tasks:

1. Declare the parameters in the service definition. This allows [Application Designer](/ee/desktop/app-designer) to be aware of the new options.

2. Use the parameters during service construction.


#### Declare the parameters

Add the parameters available to the application. The following example adds the NodeJS version and the external port:

```
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

When you run the service template container,  a volume is mounted making the service parameters available at `/run/configuration`.

The file matches the following go struct:

```
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

```
type ConfiguredService struct {
	ID      string            `json:"serviceId,omitempty"`
	Name    string            `json:"name,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}
```

You can then use the file to obtain values for the parameters and use this information based on your requirements. However, in most cases, the file is used to interpolate the variables. Therefore, we provide a utility called `interpolator` that expands variables in templates.

Basically, `interpolator` is an image that:

- takes a folder (assets folder) and a service parameter map as input,
- replaces variables in the input folder using the parameters specified by the user (for example, the service name, external port, etc), and
- writes the interpolated parameters to the application destination.

To use the `interpolator` image, update `my-service/Dockerfile` to use the following Dockerfile:

```
FROM dockertemplate/interpolator:v0.0.3-beta1
COPY assets .
```

> **Note:** The interpolator tag must match the version used in docker template. Verify this using the `docker template version` command .

This places the  interpolator image in the `/assets` folder and copies the folder to the target `/project` folder. If you prefer to do this manually, use a Dockerfile instead:

```
WORKDIR /assets
CMD ["/interpolator", "-config", "/run/configuration", "-source", "/assets", "-destination", "/project"]
```

When this is complete, use the newly added node option in `my-service/assets/Dockerfile`, by replacing the line:

`FROM node:9`

with

`FROM node:{{ .Parameters.node }}`

Now, build and push the image to your repository.

For instructions on building and pushing an image to your repository, see [Build and push the template image](#build-and-push-the-template-image).

## Create the application definition

An application template definition contains metadata that describes an application template. It contains information such as the name and description of the template, the services it contains, and  the parameters for each of the services.

Before you create an application template definition, you must create a repository that contains the services you are planning to include in the template. For more information, see [Create the repository file](#create-the-repository-file).

For example, to create an Express and MySQL application, the application definition must be similar to the following yaml file:

```
apiVersion: v1alpha1  #constant
kind: ApplicationTemplate  #constant
metadata:
  name: express-mysql #the name of the application
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

## Add templates to Docker Template

You must add the service you have created to a repository file in order to see the service when running the `docker template ls` command, or to make the service available in Application Designer. The following sections contain instructions on creating a repository file and adding the repository to the Docker Template settings.

### Create the repository file

Create a local repository file called `library.yaml` anywhere on your local drive and add the newly created service definitions and application definitions to it.

`library.yaml`

```
apiVersion: v1alpha1
generated: "2018-06-13T09:24:07.392654524Z"
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
  [...]
```

### Add the local repository to Docker Template settings

> **Note:** You can also use the instructions in this section to add templates to the Application Designer. For more information, see [Application Designer](/ee/desktop/app-designer).

Now that you have created a local repository and added service and application definitions to it, you must make Docker Template aware of these. To do this:

1. Edit `~/.docker/dockertemplate/preferences.yaml` as follows:

```
apiVersion: v1alpha1
channel: master
kind: Preferences
repositories:
- name: library-master
  url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
```

2. Add your local repository:

```
apiVersion: v1alpha1
channel: master
kind: Preferences
repositories:
- name: custom-services                # here
  url: file://path/to/my/library.yaml
- name: library-master
  url: https://docker-application-template.s3.amazonaws.com/master/library.yaml
```

After updating the `preferences.yaml` file, run `docker template ls` or restart the Application Designer and select Custom application. The new service should now be visible in the list of available services.

## Sharing custom service templates

To share a service template, you must complete the following steps:

1. Push the image to an available end point (you can use Docker hub)

2. Share the service definition (for example, GitHub)

3. Ensure the receiver has modified their `preferences.yaml` file to point to the service definition that you have shared on GitHub.