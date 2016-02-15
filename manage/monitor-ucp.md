<!--[metadata]>
+++
title = "Monitor and Manage UCP"
description = "Monitor and Manage UCP"
keywords = ["tbd, tbd"]
[menu.main]
parent="mn_ucp"
weight=-80
+++
<![end-metadata]-->

# Monitor, manage, troubleshoot your UCP installation

Intro 1-2 paras, page purpose, intended user, list steps if page is tutorial.


## Understand UCP process
* controller
* nodes
* images
* Swarm
* Discovery

## Configure UCP logging


tbd

## Getting Status from the command line
Using Client Bundle, you can get access to UCP stats, but you can also bypass UCP to get stats directly from Swarm or specific nodes. Admins can use this to troubleshoot UCP if it is not working. Steps below:

1. Download client bundle on to your desktop and
2. run (e.g. “eval $(<env.sh)”)
3. Now, running “docker info” or “docker version” provides useful statistics on the UCP cluster and nodes.
However, you can also get access directly to Swarm or specific nodes, bypassing UCP.
4. To do this, type “docker –H tcp:<Swarm_IP or Node_IP> info”.
The “-H” flag allows you to access stats for a specific host
5. Can also run “docker –H tcp:<Swarm_IP or Node_IP> version” for version data

## Review UCP logs


tbd

## Requesting support
