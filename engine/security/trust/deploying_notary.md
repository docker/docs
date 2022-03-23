---
description: Deploying Notary
keywords: trust, security, notary, deployment
title: Deploy Notary Server with Compose
---

The easiest way to deploy Notary Server is by using Docker Compose. To follow the procedure on this page, you must have already [installed Docker Compose](../../../compose/install.md).

1. Clone the Notary repository.

       git clone https://github.com/theupdateframework/notary.git

2. Build and start Notary Server with the sample certificates.

       docker-compose up -d


  For more detailed documentation about how to deploy Notary Server, see the [instructions to run a Notary service](https://github.com/theupdateframework/notary/blob/master/docs/running_a_service.md){:target="_blank" rel="noopener" class="_"} as well as [the Notary repository](https://github.com/theupdateframework/notary){:target="_blank" rel="noopener" class="_"} for more information.
3. Make sure that your Docker or Notary client trusts Notary Server's certificate before you try to interact with the Notary server.

See the instructions for [Docker](../../reference/commandline/cli.md#notary) or
for [Notary](https://github.com/docker/notary#using-notary) depending on which one you are using.

## If you want to use Notary in production

Check back here for instructions after Notary Server has an official
stable release. To get a head start on deploying Notary in production, see
[the Notary repository](https://github.com/theupdateframework/notary).
