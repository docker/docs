---
description: Frequently asked questions
keywords: ibm faqs
title: Docker for IBM Cloud frequently asked questions (FAQs)
---

## How do I sign up?
Docker EE for IBM Cloud is an unmanaged, native Docker environment within IBM Cloud that runs Docker Enterprise Edition software. Docker EE for IBM Cloud is available on **December 20th 2017 as a closed Beta**.

[Request access to the beta](mailto:sealbou@us.ibm.com). Once you do, we'll be in touch shortly!

## What IBM Cloud infrastructure permissions do I need?

  To provision the resources that make up a Docker swarm, the account administrator needs to enable certain permissions for users in the [IBM Cloud infrastructure customer portal](https://control.softlayer.com/).

  You can navigate to user permissions by going to **Account > Users > User name > Permissions**.

  Make sure that you enable the permissions in the following table.

  * The **View Only** user role does not have any of these enabled by default.
  * The **Basic User** role has some of these enabled by default.
  * The **Super User** role has most of these enabled by default.

  > Save your setting changes!
  >
  > Don't forget to click **Set Permissions** as you go through the tabs of each permission set so that you don't lose your settings.

  <table summary="The minimum user permissions that are required to provision and manage a Docker EE swarm mode cluster for IBM Cloud.">
  <caption>Table 1. The minimum user permissions that are required to provision and manage a Docker EE swarm mode cluster for IBM Cloud.
  </caption>
  <thead>
  <th colspan="1">Permissions set</th>
  <th colspan="1">Description</th>
  <th colspan="1">Required permissions</th>
  </thead>
  <tbody>
  <tr>
  <td>Devices</td>
  <td>Connect to and configure your VSI, load balancers, and firewalls.</td>
  <td>
  <ul>
  <li>View hardware detail</li>
  <li>View virtual server details</li>
  <li>Hardware firewall</li>
  <li>Software firewall manage</li>
  <li>Manage load balancers</li>
  <li>Manage device monitoring</li>
  <li>Reboot server and view IPMI system information</li>
  <li>Issue OS Reloads and initial rescue kernel<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Manage port control</li>
  </ul>
  </td>
  </tr>
  <tr>
  <td>Network</td>
  <td>Provision, connect, and expose IP addresses.</td>
  <td>
  <ul>
  <li>Add compute with public network port<a href="#edge-footnote2"><sup>2</sup></a></li>
  <li>View bandwidth statistics</li>
  <li>Add IP addresses</li>
  <li>Manage email delivery service</li>
  <li>Manage network VLAN spanning<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Manage security groups<a href="#edge-footnote2"><sup>2</sup></a></li>
  </ul></td>
  </tr>
  <tr>
  <td>Services</td>
  <td>Provision and manage services such as CDN, DNS records, SSH keys, NFS storage volumes.</td>
  <td>
  <ul>
  <li>View CDN bandwidth statistics</li>
  <li>Vulnerability scanning</li>
  <li>Manage CDN account<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Manage CDN file transfers<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>View licenses</li>
  <li>Manage DNS, reverse DNS, and WHOIS</li>
  <li>Antivirus/spyware</li>
  <li>Host IDS</li>
  <li>Manage SSH keys<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Manage storage<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>View Certificates (SSL)<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Manage Certificates (SSL)<a href="#edge-footnote1"><sup>1</sup></a></li>
  </ul>
  </td>
  </tr>
  <tr>
  <td>Account</td>
  <td>General settings to provision or remove services and instances.</td>
  <td>
  <ul>
  <li>View account summary</li>
  <li>Manage notification subscribers</li>
  <li>Add/upgrade cloud instances<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Cancel server<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Cancel services<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Add server<a href="#edge-footnote1"><sup>1</sup></a></li>
  <li>Add/upgrade services<a href="#edge-footnote1"><sup>1</sup></a></li>
  </ul>
  </td>
  </tr></tbody></table>

`1`: A **Basic User** needs these permissions added to the account.
{: id="edge-footnote1" }

`2`: Both **Basic** and **Super** users need these permissions added to the account.
{: id="edge-footnote2" }

## Which IBM Cloud region and locations (data centers) will this work with?

Docker EE for IBM Cloud is available in the following IBM Cloud regions and locations (data centers).

| Region | Region Prefix | Cities | Available locations |
| --- | --- | --- | --- |
| Frankfurt region | `eu-de`| Frankfurt, Paris | `fra02`, `par01` |
| United Kingdom | `eu-gb` | London | `lon04` |
| Sydney | `au-syd` | Hong Kong, Sydney | `hkg02`, `syd01`, `syd04` |
| US South | `ng` | Dallas, Toronto, Washington DC | `dal12`, `dal13`, `tor01`, `wdc06`, `wdc07`|

> Default location
>
> By default, clusters are created in US South, `wdc07`.

## Where are my container logs and metrics?

You must enable logging. See [Enabling logging and metric data for your swarm](logging.html) for more information.

## Why don't `bx d4ic` commands work?

The Docker EE for IBM Cloud CLI plug-in simplifies your interaction with IBM Cloud infrastructure resources. As such, many `bx d4ic` commands require you to provide your infrastructure account user name and API key credentials as options during the command (`--sl-user <user.name.1234567> --sl-api-key <api-key>`).

Instead of including these in each command, you can [set your environment variables](/docker-for-ibm-cloud/index.md#set-infrastructure-environment-variables).

## Why can't I target an organization or space in IBM Cloud?

Before you can target an organization or space in IBM Cloud, the account owner or administrator must set up the organization or space. See [Creating organizations and spaces](https://console.bluemix.net/docs/admin/orgs_spaces.html#orgsspacesusers) for more information.

## Can I manually change the load balancer configuration?

No. If you make any manual changes to the load balancer, they are removed the next time that the load balancer is updated or swarm changes are made. This is because the swarm service configuration is the source of record for service ports. If you add listeners to the load balancer manually, they could conflict with what is in cluster, and cause issues.

## How do I run administrative commands?

SSH into a manager node. Manager nodes are accessed on port 56422.

**Tip**: Because this port differs from the default (-p 22), you can add an alias to your `.profile` to make the SSH process simpler:

```none
alias ssh-docker='function __t() { ssh-keygen -R [$1]:56422 > /dev/null 2>&1; ssh -A -p 56422 -o StrictHostKeyChecking=no docker@$1; unset -f __t; }; __t'
```

## Are there any known issues?

Yes. News, updates, and known issues are recorded by version on the [Release notes](release-notes.md) page.

## Where do I report problems or bugs?

Contact us through email at docker-for-ibmcloud-beta@docker.com.

If your stack is misbehaving, run the following diagnostic tool from one of the managers to collect your docker logs and send them to Docker:

```bash
$ docker-diagnose
OK hostname=manager1
OK hostname=worker1
OK hostname=worker2
Done requesting diagnostics.
Your diagnostics session ID is 1234567890-xxxxxxxxxxxxxx
Please provide this session ID to the maintainer debugging your issue.
```

> **Note**: Your output may be slightly different from the above, depending on your swarm configuration.
