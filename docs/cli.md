<!--[metadata]>
+++
title = "Notary CLI"
description = "Description of the Notary CLI"
keywords = ["docker, notary, trust, image, signing, repository, cli"]
[menu.main]
parent="mn_notary"
weight=4
+++
<![end-metadata]-->

# Notary CLI

## Notary Server

The default notary server URL is [https://notary-server:4443/]. This default value can overridden (by priority order):

- by specifying the option `--server/-s` on commands requiring call to the notary server.
- by setting the `NOTARY_SERVER_URL` environment variable.
