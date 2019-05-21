---
title: Docker Template API reference
description: Docker Template API reference
keywords: application, template, API, definition
---

This page contains information about the Docker Template API reference.

## Service template definition

The following section provides information about the valid parameters that you can use when you create a service template definition.

```
apiVersion: v1alpha1
kind: ServiceTemplate
metadata:
 name: angular
 platforms:
   - linux
spec:
 title: Angular
 description: Angular service
 icon: https://cdn.worldvectorlogo.com/logos/angular-icon-1.svg
 source:
   image: docker.io/myorg/myservice:version
 parameters:
 - name: node
   description: Node version
   type: enum
   defaultValue: "9"
   values:
   - value: "10"
     description: "10"
   - value: "9"
     description: "9"
   - value: "8"
     description: "8"
 - name: externalPort
   description: External port
   defaultValue: "8080"
   type: hostPort
```

### root

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| apiVersion             |yes                 |  The api format version. Current latest is v1alpha1|
|kind| yes|The kind of object. Must be `ServiceTemplate` For services templates.|

### metadata

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
|name |yes                 | The identifier for this service. Must be unique within a given library. |
|platform| yes|A list of allowed target platforms. Possible options are `windows` and `linux`|

### spec

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| title     |yes  |The label for this service, as displayed when listed in `docker template` commands or in the `application-designer`|
|description| no|A short description for this service|
|icon|no|An icon representing the service. Only used in the Application Designer|

### spec/source

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| image |yes| The name of the image associated with this service template. Must be in full `repo/org/service:version` format|

### spec/parameters

The parameters section allows to specify the input parameters that are going to be used by the service.

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
|name            |yes| The identifier for this parameter. Must be unique within the service parameters. |
|description| no|A short description of the parameter. Will be used as label in the Application Designer|
|type| yes|The type of the parameter. Possible options are: `string` - The default type, with no validation or specific features. `enum` - Allow the user to choose a value included in a specific list of options. Must specify the values parameter. `hostPort` - Specify that this parameter is a port that is going to be exposed. Use port format regexp validation, and avoid duplicate ports within an application.|
|defaultValue| yes|The default value for this parameter. For enum type, must be a valid value from the values list.|
|values| no|For enum type, specify a list of value with a value/description tuple.|

## Application template definition

The following section provides information about the valid parameters that you can use when you create a application template definition.

```
apiVersion: v1alpha1
kind: ApplicationTemplate
metadata:
 name: nginx-flask-mysql
 platforms:
   - linux
spec:
 title: Flask / NGINX / MySQL application
 description: Sample Python/Flask application with an Nginx proxy and a MySQL database
 services:
 - name: back
   serviceId: flask
   parameters:
     externalPort: "80"
 - name: db
   serviceId: mysql
 - name: proxy
   serviceId: nginx
```

### root

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| apiVersion             |yes                 |  The api format version. Current latest is v1alpha1|
|kind| yes|The kind of object. Must be `ApplicationTemplate` For application templates.|

### metadata

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
|name |yes                 | The identifier for this application template. Must be unique within a given library.|
|platform| yes|A list of allowed target platforms. Possible options are `windows` and `linux`|

### spec

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| title     |yes  |The label for this application template, as displayed when listed in `docker template` commands or in `application-designer` |
|description| no|A short description for this service|

### spec/services

This section lists the service templates used in the application.

| Parameter              |Required?              | Description |
| :----------------------|:----------------------|:----------------------------------------|
| name |yes|The name of the service. It will be used for image name and for subfolder within the application structure. |
|serviceId  |yes|The id of the service to use (equivalent to the metadata/name field of the service) |
| parameters |no|A map (string to string) that can be used to override the default values of the service parameters.|

## Service configuration file

The file is mounted at `/run/configuration` in every service template container and contains the template context in a JSON format.

| Parameter             |Description |
| :----------------------|:----------------------|
|ServiceId |The service id|
| name |The name of the service as specified by the application template or overridden by the user|
|parameters  |A map (string to string) containing the service’s parameter values.|
| targetPath |The destination folder for the application on the host machine.|
|namespace  |The service image’s namespace (org and user)|
|services  |A list containing all the services of the application (see below)|

### Attributes

The items in the services list contains the following attributes:

| Parameter             |Description |
| :----------------------|:----------------------|
|serviceId  |The service id|
| name |The name of the service as specified by the application template or overridden by the user|
| parameters |A map (string to string) containing the service’s parameter values.|