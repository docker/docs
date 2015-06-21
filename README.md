# Notary is still a work in progress and we invite contributions and reviews from the security community. It will need to go through a formal security review process before it should be used in production.

# Notary

## Overview

The Notary project comprises a [server](cmd/notary-server) and a [client](cmd/notary) for running and interacting
with trusted collections. Please see the READMEs for the individual tools for
more information on their specifics.

## Goals

Notary aims to make the internet more secure by making it easy for people to
publish and verify content. We often rely on TLS to secure our communications
with a web server which is inherently flawed, as any compromise of the server
enables malicious content to be substituted for the legitimate content.

With Notary, publishers can sign their content offline using keys kept highly
secure. Once the publisher is ready to make the content available, they can
push their signed trusted collection to a Notary Server.

Consumers, having acquired the publisher's public key through a secure channel,
can then communicate with any notary server or (insecure) mirror, relying
only on the publisher's key to determine the validity and integrity of the
received content.
