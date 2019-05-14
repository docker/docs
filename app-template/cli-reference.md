---
title: Docker Template CLI reference
description: Docker Template CLI reference
keywords: Docker, application template, CLI, Application Designer,
---

This page provides information about the `docker template` command.

## Overview

Docker Template is a CLI plugin that introduces a top-level `docker template`command that allows users to create new Docker applications using a library of templates. With `docker template`, you can scaffold a full project structure for a chosen technical stack or a set of technical stacks using the best practices pre-configured in a generated Dockerfile and docker-compose file.

For more information about Docker Template, see [Working with Docker Template](/ee/docker-template/working-with-template).

## `docker template` commands

To view the commands and sub-commands available in `docker template`, run:

`docker template --help`

```
Usage:  docker template COMMAND

Use templates to quickly create new services

Commands:
  inspect     Inspect service templates or application templates
  list        List available templates with their information
  scaffold    Choose an application template or service template(s) and scaffold a new project
  version     Print version information

Run 'docker template COMMAND --help' for more information on a command.
```

### inspect

The `docker template inspect` command allows you to view the details of the template such as service parameters, default values, and for application templates, the list of services included in the application.

```
Usage:  docker template inspect <service or application>

Inspect service templates or application templates

Options:
      --format string   Configure the output format (pretty|json|yaml)
                        (default "pretty")
```

For example:

```
docker template inspect react-java-mysql
NAME: react-java-mysql
TITLE: React / Spring / MySQL application
DESCRIPTION: Sample React application with a Spring backend and a MySQL database
SERVICES:

 * PARAMETERS FOR SERVICE: front (react)
NAME           DESCRIPTION     TYPE       DEFAULT VALUE   VALUES
node           Node version    enum       9               10, 9, 8
externalPort   External port   hostPort   8080


 * PARAMETERS FOR SERVICE: back (spring)
NAME             DESCRIPTION               TYPE       DEFAULT VALUE           VALUES
java             Java version              enum       9                       10, 9, 8
groupId          Group Id                  string     com.company
artifactId       Artifact Id               string     project
appName          Application name          string     New App
appDescription   Application description   string     My new SpringBoot app
externalPort     External port             hostPort   8080


 * PARAMETERS FOR SERVICE: db (mysql)
NAME      DESCRIPTION   TYPE   DEFAULT VALUE   VALUES
version   Version       enum   5.7             5.7
```

### list

The `docker template list` command lists the available service and application templates.

```
Usage:  docker template list
List available templates with their information

Aliases:
  list, ls

Options:
      --format string   Configure the output format (pretty|json|yaml)
                        (default "pretty")
      --type string     Filter by type (application|service|all) (default
                        "all")
```

For example:

`docker template list`

```
NAME                    TYPE          DESCRIPTION
aspnet-mssql            application   Sample asp.net core application with mssql database
nginx-flask-mysql       application   Sample Python/Flask application with an Nginx proxy and a MySQL database
nginx-golang-mysql      application   Sample Golang application with an Nginx proxy and a MySQL database
nginx-golang-postgres   application   Sample Golang application with an Nginx proxy and a PostgreSQL database
react-java-mysql        application   Sample React application with an Spring backend and a MySQL database
react-express-mysql     application   Sample React application with a NodeJS backend and a MySQL database
sparkjava-mysql         application   Java application and a MySQL database
spring-postgres         application   Sample Java application with Spring framework and a Postgres database
angular                 service       Angular service
aspnetcore              service       A lean and composable framework for building web and cloud applications
consul                  service       A highly available and distributed service discovery and KV store
django                  service       A high-level Python Web framework
express                 service       NodeJS web application with Express server
flask                   service       A microframework for Python based on Werkzeug, Jinja 2 and good intentions
golang                  service       A powerful URL router and dispatcher for golang
gwt                     service       GWT (Google Web Toolkit) / Java service
jsf                     service       JavaServer Faces technology establishes the standard for building server-side user interfaces.
mssql                   service       Microsoft SQL Server for Docker Engine
mysql                   service       Official MySQL image
nginx                   service       An HTTP and reverse proxy server
postgres                service       Official PostgreSQL image
rails                   service       A web-application framework that includes everything needed to create database-backed web applications
react                   service       React/Redux service with Webpack hot reload
sparkjava               service       A micro framework for creating web applications in Java 8 with minimal effort
spring                  service       Customizable Java/Spring template
vuejs                   service       VueJS service
```

### scaffold

The `docker template scaffold` command allows you to generate a project structure for a template.

```
Usage:  docker template scaffold application [<alias=service>...] OR scaffold [alias=]service [<[alias=]service>...]

Choose an application template or service template(s) and scaffold a new project

Examples:
docker template scaffold react-java-mysql -s back.java=10 -s front.externalPort=80
docker template scaffold react-java-mysql java=back reactjs=front -s reactjs.externalPort=80
docker template scaffold back=spring front=react -s back.externalPort=9000
docker template scaffold react-java-mysql --server=myregistry:5000 --org=myorg

Options:
      --build             Run docker-compose build after deploy
      --name string       Application name
      --org string        Deploy to a specific organization / docker hub
                          user (if not specified, it will use your
                          current hub login)
      --path string       Deploy to a specific path
      --platform string   Target platform (linux|windows) (default "linux")
      --server string     Deploy to a specific registry server (host[:port])
  -s, --set stringArray   Override parameters values (service.name=value)
```

For example:

`docker template scaffold react-java-mysql`

If you want to change some of the parameter values (exposed port, specific version, etc.) you can pass additional parameters, and reference the service name it applies to with `--set` or `-s`. 

For example:

`docker template scaffold react-java-mysql -s back.java=10 -s front.externalPort=80`

By default, the `docker template scaffold` command generates the project structure in the current folder. However, you can specify another folder using the `--path` parameter. 

For example:

`docker template scaffold react-java-mysql --path /xxx`

You can also change service names by providing aliases when scaffolding either an application template or a list of service templates.

For example:

`docker template scaffold react-java-mysql java=back reactjs=front -s reactjs.externalPort=80`

### version

The `docker template version` command displays the Docker Template version number.

```
Usage:  docker template version

Print version information
```

For example:

`docker template version`

```
Version:      d6c11e577c592aad69d34db6d4dc740d65291e36
Git Commit:   96ea0063b0c9aaa0cc5b5ff811b51a6e2e752be9
```