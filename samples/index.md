---
title: Samples
---

## Docker Enterprise Edition walkthrough

Test drive a running instance of Docker EE, without installing anything; this
walkthrough instructs you on how to perform various common tasks with the Docker
EE dashboard. You'll probably never want to use Docker without a UI again.

[Launch Docker Enterprise Edition walkthrough](https://dockertrial.com){: class="button outline-btn" onclick="ga('send', 'event', 'EE Trial Referral', 'Samples Page', 'Click');"}

## Play with Docker

Learn how to use Docker through our guided labs that let you actually
run Docker in a sandboxed environment we host for you, all in your browser.
These labs cover general Docker use cases that are useful for both Docker
Community Edition and Enterprise Edition.

<script language="JavaScript" src="/js/jquery.js"></script>
<script language="JavaScript">
    var output = '';
		$.ajax({
         type: "get",
         url: "http://training.play-with-docker.com/feed.xml",
         dataType: "xml",
         success: function(data) {
             /* handle data here */
             console.log(data)
         },
         error: function(xhr, status) {
             /* handle error here */
             console.log(status)
         }
     });
</script>

<div id="PWDTable"></div>

{% assign labsbase = "https://github.com/docker/labs/tree/master" %}

## GitHub-hosted samples

Learn how to develop and ship containerized applications, by walking through a
GitHub-hosted sample that exhibits canonical practices, with README walkthroughs
explaining things along the way as a lab. These samples are from the
[Docker Labs repository]({{ labsbase }}).

| Sample | Description |
| ------ | ----------- |
| [Docker for Beginners]({{ labsbase }}/beginner/){: target="_blank"} | A good "Docker 101" course. |
| [Docker Swarm mode]({{ labsbase}}/swarm-mode){: target="_blank"} | Use Docker for natively managing a cluster of Docker Engines called a swarm. |
| [Configuring developer tools and programming languages]({{ labsbase }}/developer-tools/README.md){: target="_blank"} | How to set-up and use common developer tools and programming languages with Docker. |
| [Live Debugging Java with Docker]({{ labsbase }}/developer-tools/java-debugging){: target="_blank"} | Java developers can use Docker to build a development environment where they can run, test, and live debug code running within a container. |
| [Docker for Java Developers]({{ labsbase }}/developer-tools/java/){: target="_blank"} | Offers Java developers an intro-level and self-paced hands-on workshop with Docker. |
| [Live Debugging a Node.js application in Docker]({{ labsbase }}/developer-tools/nodejs-debugging){: target="_blank"} | Node developers can use Docker to build a development environment where they can run, test, and live debug code running within a container. |
| [Dockerizing a Node.js application]({{ labsbase }}/developer-tools/nodejs/porting/){: target="_blank"} | This tutorial starts with a simple Node.js application and details the steps needed to Dockerize it and ensure its scalability. |
| [Docker for ASP.NET and Windows containers]({{ labsbase }}/windows/readme.md){: target="_blank"} | Docker supports Windows containers, too! Learn how to run ASP.NET, SQL Server, and more in these tutorials. |
| [Docker Security]({{ labsbase }}/security/README.md){: target="_blank"} | How to take advantage of Docker security features. |
| [Building a 12-factor application with Docker]({{ labsbase}}/12factor){: target="_blank"} | Use Docker to create an app that conforms to Heroku's "12 factors for cloud-native applications." |

## Library references

These docs are imported from
[the official Docker Library docs](https://github.com/docker-library/docs/),
and help you use some of the most popular software that has been
"Dockerized" into Docker images.

| Image name | Description |
| ---------- | ----------- |
{% for page in site.samples %}| [{{ page.title }}]({{ page.url }}) | {{ page.description | strip }} |
{% endfor %}

## Sample applications

Run popular software using Docker.

| Sample | Description |
| ------ | ----------- |
| [apt-cacher-ng](/engine/examples/apt-cacher-ng) | Run a Dockerized apt-cacher-ng instance. |
| [ASP.NET Core + SQL Server on Linux](/compose/aspnet-mssql-compose) | Run a Dockerized ASP.NET Core + SQL Server environment. |
| [CouchDB](/engine/examples/couchdb_data_volumes) | Run a Dockerized CouchDB instance. |
| [Django + PostgreSQL](/compose/django/) | Run a Dockerized Django + PostgreSQL environment. |
| [PostgreSQL](/engine/examples/postgresql_service) | Run a Dockerized PosgreSQL instance. |
| [Rails + PostgreSQL](/compose/rails/) | Run a Dockerized Rails + PostgreSQL environment. |
| [Riak](/engine/examples/running_riak_service) | Run a Dockerized Riak instance. |
| [SSHd](/engine/examples/running_ssh_service) | Run a Dockerized SSHd instance. |
