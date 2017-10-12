---
description: Introduction and Overview of Compose
keywords: documentation, docs, docker, compose, orchestration, containers
title: Overview of Docker Compose
---

>**Looking for Compose file reference?** [Find the latest version here](/compose/compose-file/index.md).

Compose is a tool for defining and running multi-container Docker applications.
With Compose, you use a YAML file to configure your application's services.
Then, with a single command, you create and start all the services
from your configuration. To learn more about all the features of Compose
see [the list of features](overview.md#features).

Compose works in all environments: production, staging, development, testing, as
well as CI workflows. You can learn more about each case in [Common Use
Cases](overview.md#common-use-cases).

Using Compose is basically a three-step process.

1. Define your app's environment with a `Dockerfile` so it can be reproduced
anywhere.

2. Define the services that make up your app in `docker-compose.yml`
so they can be run together in an isolated environment.

3. Lastly, run
`docker-compose up` and Compose will start and run your entire app.

A `docker-compose.yml` looks like this:

    version: '3'
    services:
      web:
        build: .
        ports:
        - "5000:5000"
        volumes:
        - .:/code
        - logvolume01:/var/log
        links:
        - redis
      redis:
        image: redis
    volumes:
      logvolume01: {}

For more information about the Compose file, see the
[Compose file reference](compose-file/index.md).

Compose has commands for managing the whole lifecycle of your application:

 * Start, stop and rebuild services
 * View the status of running services
 * Stream the log output of running services
 * Run a one-off command on a service

## Compose documentation

- [Installing Compose](install.md)
- [Getting Started](gettingstarted.md)
- [Get started with Django](django.md)
- [Get started with Rails](rails.md)
- [Get started with WordPress](wordpress.md)
- [Frequently asked questions](faq.md)
- [Command line reference](./reference/index.md)
- [Compose file reference](compose-file/index.md)

## Features

The features of Compose that make it effective are:

* [Multiple isolated environments on a single host](overview.md#Multiple-isolated-environments-on-a-single-host)
* [Preserve volume data when containers are created](overview.md#preserve-volume-data-when-containers-are-created)
* [Only recreate containers that have changed](overview.md#only-recreate-containers-that-have-changed)
* [Variables and moving a composition between environments](overview.md#variables-and-moving-a-composition-between-environments)

### Multiple isolated environments on a single host

Compose uses a project name to isolate environments from each other. You can make use of this project name in several different contexts:

* on a dev host, to create multiple copies of a single environment (e.g., you want to run a stable copy for each feature branch of a project)
* on a CI server, to keep builds from interfering with each other, you can set
  the project name to a unique build number
* on a shared host or dev host, to prevent different projects, which may use the
  same service names, from interfering with each other

The default project name is the basename of the project directory. You can set
a custom project name by using the
[`-p` command line option](./reference/overview.md) or the
[`COMPOSE_PROJECT_NAME` environment variable](./reference/envvars.md#compose-project-name).

### Preserve volume data when containers are created

Compose preserves all volumes used by your services. When `docker-compose up`
runs, if it finds any containers from previous runs, it copies the volumes from
the old container to the new container. This process ensures that any data
you've created in volumes isn't lost.

If you use `docker-compose` on a Windows machine, see
[Environment variables](reference/envvars.md) and adjust the necessary environment
variables for your specific needs.


### Only recreate containers that have changed

Compose caches the configuration used to create a container. When you
restart a service that has not changed, Compose re-uses the existing
containers. Re-using containers means that you can make changes to your
environment very quickly.


### Variables and moving a composition between environments

Compose supports variables in the Compose file. You can use these variables
to customize your composition for different environments, or different users.
See [Variable substitution](compose-file.md#variable-substitution) for more
details.

You can extend a Compose file using the `extends` field or by creating multiple
Compose files. See [extends](extends.md) for more details.


## Common Use Cases

Compose can be used in many different ways. Some common use cases are outlined
below.

### Development environments

When you're developing software, the ability to run an application in an
isolated environment and interact with it is crucial. The Compose command
line tool can be used to create the environment and interact with it.

The [Compose file](compose-file.md) provides a way to document and configure
all of the application's service dependencies (databases, queues, caches,
web service APIs, etc). Using the Compose command line tool you can create
and start one or more containers for each dependency with a single command
(`docker-compose up`).

Together, these features provide a convenient way for developers to get
started on a project. Compose can reduce a multi-page "developer getting
started guide" to a single machine readable Compose file and a few commands.

### Automated testing environments

An important part of any Continuous Deployment or Continuous Integration process
is the automated test suite. Automated end-to-end testing requires an
environment in which to run tests. Compose provides a convenient way to create
and destroy isolated testing environments for your test suite. By defining the full environment in a [Compose file](compose-file.md) you can create and destroy these environments in just a few commands:

    $ docker-compose up -d
    $ ./run_tests
    $ docker-compose down

### Single host deployments

Compose has traditionally been focused on development and testing workflows,
but with each release we're making progress on more production-oriented features. You can use Compose to deploy to a remote Docker Engine. The Docker Engine may be a single instance provisioned with
[Docker Machine](/machine/overview.md) or an entire
[Docker Swarm](/engine/swarm/index.md) cluster.

For details on using production-oriented features, see
[compose in production](production.md) in this documentation.


## Release Notes

To see a detailed list of changes for past and current releases of Docker
Compose, please refer to the
[CHANGELOG](https://github.com/docker/compose/blob/master/CHANGELOG.md).

## Getting help

Docker Compose is under active development. If you need help, would like to
contribute, or simply want to talk about the project with like-minded
individuals, we have a number of open channels for communication.

* To report bugs or file feature requests: please use the [issue tracker on Github](https://github.com/docker/compose/issues).

* To talk about the project with people in real time: please join the
  `#docker-compose` channel on freenode IRC.

* To contribute code or documentation changes: please submit a [pull request on Github](https://github.com/docker/compose/pulls).

For more information and resources, please visit the [Getting Help project page](/opensource/get-help/).

#pt-BR
---
descrição: Introdução e Visão Geral do Compose
palavras-chave: documentação, docs, docker, compose, orchestration, containers
título: Visão geral do Docker Compose
---

>**Procurando por Arquivos de Referência Compose?** [Encontre a Última versão aqui](/compose/compose-file/index.md).

Compose é uma ferramenta para definição e execução de aplicações Docker multi-container.
Com Compose, você usa um arquivo YAML para configurar seus serviçoes de aplicações.
Então, com um único comando, você creia e inicia todos os serviços
de sua configuração. Para aprender mais sobre características do Compose
veja [a lista de características](overview.md#features).

Compose serve para vários ambientes: produção, estadiamento, desenvolvimento, teste, assim
como CI workflows. Você pode aprender mais sobre cada caso em [Common Use
Cases](overview.md#common-use-cases).

Usaar Compose é basicamente um processo de três passos.

1. Defina seu ambiente de desenvolvimento do app com um `Dockerfile` então ele pode ser reproduzido
em qualquer lugar.

2. Defina os serviços que terão seu app em `docker-compose.yml`
assim eles podem rodar juntos em um ambiente isolado.

3. Por fim, execute
`docker-compose up` e Compose irá iniciar e executar seu app.

Um `docker-compose.yml` parece com isso:

    version: '3'
    services:
      web:
        build: .
        ports:
        - "5000:5000"
        volumes:
        - .:/code
        - logvolume01:/var/log
        links:
        - redis
      redis:
        image: redis
    volumes:
      logvolume01: {}

Para mais informações sobre o arquivo Compose, veja o
[Compose file reference](compose-file/index.md).

Compose tem comandos para gerenciamente de todo ciclo de vida da sua aplicação:

 * Iniciar, parar e reconstruir serviços
 * Ver o status dos serviços em execução
 * Transmitir o log de saída dos serviços em execução.
 * Executar um comando liga-desliga em um serviço.

## Documentação do Compose

- [Instalando Compose](install.md)
- [Iniciando](gettingstarted.md)
- [Iniciando com Django](django.md)
- [Iniciando com Rails](rails.md)
- [Iniciando com WordPress](wordpress.md)
- [Perguntas frequentes](faq.md)
- [Linha de comando de Referência](./reference/index.md)
- [Arquivo de referência Compose](compose-file/index.md)

## Características

As características do Compose que o fazem efetivo são:

* [Múltiplos ambientes isolados em um único host](overview.md#Multiple-isolated-environments-on-a-single-host)
* [Preserva o volume de dados quando os containers são criados](overview.md#preserve-volume-data-when-containers-are-created)
* [Recria apenas os containers que tiveram mudança](overview.md#only-recreate-containers-that-have-changed)
* [Variáveis e mudanças na composição entre ambientes](overview.md#variables-and-moving-a-composition-between-environments)

### Múltiplos ambientes isolados em um único host

Compose usa um nome de projeto para isolar ambientes um do outro. Você pode fazer uso do nome desse projeto em diferentes contextos:

* em um dev host, para criar várias cópias de um único ambiente (ex: você pode querer executar uma cópia estável para cara branck de um projeto)
* em um CI server, para manter as contruções sem interferências, você pode atribuir o nome do projeto a um único número.
* em um host compartilhado ou dev host, para prever diferentes projetos de usarem os mesmos nomes de serviço, para que não interfiram um no outro.

O nome do projeto por padrão é o nome base do diretório do projeto. Você pode atribuir
um outro nome usando a
[`-p` linha de comando](./reference/overview.md) or the
[`COMPOSE_PROJECT_NAME` environment variable](./reference/envvars.md#compose-project-name).

### Preserva o volume de dados quando containers são criados

Compose preserva todo volume usado pelo seu serviço. Quando `docker-compose up`
é executado, se for ecnontrado algum container de execuções anteriores, ele copia o volume do
container antigo para o novo container. Esse processo assegura que nenhum dado
que você criou no volume será perdido.

Se você usa `docker-compose`em uma máquina Windows, veja
If you use `docker-compose` on a Windows machine, see
[Variáveis de ambiente](reference/envvars.md) e ajuste as variáveis de ambiente
necessárias para suas necessidades específicas.


### Recrie apenas containers que foram mudados

Compose grava a configuração usada para criar um container. Quando você
reinicia um serviço que não foi mudado, Compose reusará os containers
existentes. Reusar containers seignifica que você pode fazer mudanças em seu
ambiente muito rapidamente.


### Variáveis e mobilidade de uma composição entre ambientes

Compose suporta variáveis no arquivo Compose. Você pode usar essas variáveis
para customizar sua composição para diferentes ambientes, ou diferentes usuários.
Veja [Substituição de Variáveis](compose-file.md#variable-substitution) para mais
detalhes.

você pode extender um arquivo Compose usando o campo `extends` or criando
vários arquivos Compose. Veja [extends](extends.md) para mais detalhes.


## Casos de uso comuns

Compode pode ser usado de diferentes formas. Alguns casos de uso comuns estão
relacionados abaixo.


### Ambiente de desenvolvimento

Quando você está desenvolvendo software, a habilidade de executar uma aplicação em um
ambiente isolado e a interação com este é crucial. A ferramenta linha de comando Compose
pode ser usada para criar o ambiente e interagir com ele.

O [Arquivo Compose](compose-file.md) fornece um meio de documentar e configurar
tudo do serviço de dependências da aplicação (databases, queues, caches,
web service APIs, etc). Usando a ferramenta linha de comando do Compose você pode criar
e inicializar um ou amis containers para cada dependência com um único comando
(`docker-compose up`).

Juntas, essas características fornecem uma forma conveniente para desenvolvedores
iniciarem seu projeto. Compose pode reduzir uma multi-page "Guia iniciante do
desenvolvedor" para uma única máquina capaz de ler o arquivo Compose e alguns comandos.

### Testes em ambientes automatizados

Uma parte importante de qualquer aplicação ou processo de integração contínua
é o sistema automatizado de teste. Testes automatizados fim-a-fim requerem um
ambiente no qual possam rodar os testes. Compode fornece uma forma conveniente de criar
e destruir ambientes de testes isolados para sua aplicação de testes. Definindo todo ambiente em um [Arquivo Compose](compose-file.md) você
pode criar e destruir esses ambientes com apneas poucos comandos:

    $ docker-compose up -d
    $ ./run_tests
    $ docker-compose down

### Implementações em hosts únicos

Compose tem focado tradicionalmente em workfloes de desenvolvimento e teste,
mas com cada release nós estamos fazendo progresso características mais orientadas a produção. Você pode usar Compose para implementar em um docker remoto. O mecanismos Docker pode ser uma instância única provisionada com
Compose has traditionally been focused on development and testing workflows,
[Docker Machine](/machine/overview.md) ou um inteiro
[Docker Swarm](/engine/swarm/index.md) cluster.

Para detalhes no uso de características orientadas a produção, veja
[compose in production](production.md) nesta documentação.


## Notas de Release

Para ver uma lista detalhada de mudanças de releases passados e atuais do Docker
Compose, por favor acesse
[CHANGELOG](https://github.com/docker/compose/blob/master/CHANGELOG.md).

## Conseguindo ajuda

Docker Compose está em desenvolvimento. Se você precisa de ajuda, gostaria de
contribuir, ou simplesmente que falar sobre o projeto com ideias individuais
semelhantes, nós temos um número de canais abertos para comunicação.

* Para reportar bugs ou fazer requisições em arquivos: Por favor use o [issue tracker on Github](https://github.com/docker/compose/issues).

* Para discutir sobre o projeto com pessoas em tempo real: Por favor junte-se ao
  canal `#docker-compose` em freenode IRC.

* Para contribuir com mudanças de código ou documentação: por favor submeta um [pull request on Github](https://github.com/docker/compose/pulls).

Para mais informações e recursos, por favor visite o [Getting Help project page](/opensource/get-help/).
