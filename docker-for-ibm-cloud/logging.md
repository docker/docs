---
description: Logging for Docker for IBM Cloud
keywords: ibm, ibm cloud, logging, iaas, tutorial
title: Send logging and metric cluster data to IBM Cloud
---

You can enable Docker Enterprise Edition for IBM Cloud to send logging and metric data about the nodes and containers in your Docker EE cluster to the IBM Cloud [Log Analysis](https://console.bluemix.net/docs/services/CloudLogAnalysis/log_analysis_ov.html#log_analysis_ov) and [Monitoring](https://console.bluemix.net/docs/services/cloud-monitoring/monitoring_ov.html#monitoring_ov) services.

> Logging on services other than IBM Cloud
>
> If you want to configure logging and metrics for Docker EE to a remote logging service that is not IBM Cloud, see [Configure UCP logging](/datacenter/ucp/2.0/guides/configuration/configure-logs/).

## Enable logging and metrics

By default, logging and metrics are disabled. After you enable logging and metrics, containers are deployed to your cluster and begin to transmit data to IBM Cloud [Log Analysis](https://console.bluemix.net/docs/services/CloudLogAnalysis/log_analysis_ov.html#log_analysis_ov) and [Monitoring](https://console.bluemix.net/docs/services/cloud-monitoring/monitoring_ov.html#monitoring_ov) services.

Before you begin, make sure that you [installed the IBM Cloud CLI and Docker for IBM Cloud plug-in](/docker-for-ibm-cloud/index.md).

To enable logging and metrics:

1. Log in to IBM Cloud. If you have a federated ID, use the `--sso` option.

   ```bash
   $ bx login [--sso]
   ```

2. After logging in to IBM Cloud, target the organization and space to which you want to send the logging and metric data:

    ```bash
    $ bx target --cf
    ```

3. Connect to your swarm by setting the environment variables from the [client certificate bundle that you downloaded](administering-swarms.md#download-client-certificates). For example:

   ```bash
   $ cd filepath/to/certificate/repo && source env.sh
   ```

4. Get the name of your cluster. If you did not [set your environment variables](/docker-for-ibm-cloud/index.md#set-infrastructure-environment-variables), include your IBM Cloud infrastructure credentials.

    ```bash
    $ bx d4ic list --sl-user user.name.1234567 --sl-api-key api_key
    ```

5. Enable logging and metrics. Replace the _my_swarm_ variable with the name of your cluster and include the path to the Docker EE client certificate bundle. Include your IBM Cloud infrastructure credentials if you have not set the environment variables.

    ```bash
    $ bx d4ic logmet --swarm-name my_swarm \
    --cert-path filepath/to/certificate/repo \
    --sl-user user.name.1234567 \
    --sl-api-key api_key \
    --enable
    ```

## Disable logging and metrics
You might want to disable logging and metrics for reasons such as sending data to a different server [specified in Docker Enterprise Edition UCP](/datacenter/ucp/2.0/guides/configuration/configure-logs/). After disabling, data is no longer transmitted to IBM Cloud Log Analysis and Monitoring services.

To disable logging and metrics, get the name of the swarm and run the disable command:

  ```bash
  $ bx d4ic logmet --swarm-name my_swarm \
  --cert-path filepath/to/certificate/repo \
  --sl-user user.name.1234567 \
  --sl-api-key api_key \
  --disable
  ```

## Review logging and metrics
Use the following links to access Kibana and Grafana for data transmitted to IBM Cloud. Select the IBM Cloud organization and space that your cluster is in to view its information.

View the [IBM Cloud Log Analysis](https://console.bluemix.net/docs/services/CloudLogAnalysis/log_analysis_ov.html#log_analysis_ov) and [IBM Cloud Monitoring](https://console.bluemix.net/docs/services/cloud-monitoring/monitoring_ov.html#monitoring_ov) documentation to learn more.

| Region | Logging | Metrics|
| --- | --- | --- |
| US South | [https://logging.ng.bluemix.net/](https://logging.ng.bluemix.net/) | [https://metrics.ng.bluemix.net/](https://metrics.ng.bluemix.net/) |
| United Kingdom | [https://logmet.eu-gb.bluemix.net/](https://logmet.eu-gb.bluemix.net/)| [https://metrics.eu-gb.bluemix.net/](https://metrics.eu-gb.bluemix.net/) |
