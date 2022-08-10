---
description: Docker documentation redirect URLs
title: Docker documentation redirect URLs
published: false
---

This page is used to reserve the "/go/" URL/path for redirecting purposes. This
page is not published, and should not be visible on the documentation website.

Use the `/go/` URL prefix to provide "vanity" URLs linking to the documentation.
These URLs can be used in (for example) the `docker` CLI to provide a short,
permanent link to refer to a topic in the documentation.

Having a fixed `/go/` URL allows you to redirect the URL to the most current
location in the documentation for a given topic, which may evolve over time.

For example, the `docker` CLI can output `https://docs.docker.com/go/some-topic`
in its help output, which can be redirected to elsewhere in the documentation,
using the `redirect_from` metadata field in the target page.
