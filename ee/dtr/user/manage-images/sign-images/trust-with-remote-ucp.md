---
title: Using Docker Content Trust with a Remote UCP Cluster
description: Learn how to use a single DTR's trust data with remote UCPs.
keywords: registry, sign, trust, notary
---

For more advanced deployments, you may want to share one Docker Trusted Registry
across multiple Universal Control Planes. However, customers wanting to adopt
this model, alongside the [Only Run Signed
Images](../.../../ucp/admin/configure/run-only-the-images-you-trust.md) feature
of UCP, run into problems as each Universal Control Plane operates an
independent set of users.

Docker Content Trust with a Remote UCP gets around this problem, as User's from
a remote UCP are able to sign images in the central DTR, and still apply runtime
enforcement.

In the following example we will connect DTR to 1 remote UCP cluster, sign the
image with a User from that remote UCP cluster, and provide runtime enforcement
within the remote UCP cluster. This process could be repeated over and over,
integrating DTR with multiple remote UCP clusters, signing the image with Users
from each environment, and then providing runtime enforcement in each remote UCP
cluster seperately.

![](../../../images/remoteucp-graphic.png)

> Before attempting this guide, please familiarize yourself with [Docker Content
> Trust](engine/security/trust/content_trust/#signing-images-with-docker-content-trust)
> and [Only Run Signed
> Images](../.../../ucp/admin/configure/run-only-the-images-you-trust.md) on a
> single UCP. A lot of the concepts within this guide may be new without that
> background. 

## Prerequisites

- Cluster 1, running UCP 3.0.x or higher, with a DTR 2.5.x or higher deployed
  within the cluster.
- Cluster 2, running UCP 3.0.x or higher, it is expected that
  there is no DTR installed on this environment.
- Nodes on Cluster 2 need to trust the Certificate Authority that signed DTR's
  TLS Certificate. This can be tested by logging on to a Cluster 2 virtual
  machine and running `curl https://dtr.example.com`.
- The DTR TLS Certificate needs be properly configured, ensuring that the
  **Loadbalancer/Public Address** field has been configured, with this address
  included [within the
  certificate](../../../admin/configure/use-your-own-tls-certificates/).
- A workstation with the [Docker Client](/ee/ucp/user-access/cli/) (CE 17.12 /
  EE 1803 or newer) installed, as this contains the relevant `$ docker trust`
  commands. 

## Registering DTR with a remote Universal Control Plane

As there is no registry running within Cluster 2, by default UCP will not know
where to check for trust data. Therefore, the first thing we need to do is
register DTR within the Universal Control Plane of Cluster 2. When you normally
install Docker Trusted Registry, this registration process happens by default to
a local UCP.

> The registration process allows the remote UCP to get signature data from DTR,
> however this will not provide Single Sign On, Users on Cluster 2 will not be
> synced with Cluster 1's Universal Control Plane or Docker Trusted Registry.
> Therefore when pulling images, if the repository is private, registry
> authentication will still need to be passed as part of the service definition
> (For example in
> [Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#create-a-secret-in-the-cluster-that-holds-your-authorization-token)
> or [Docker
> Swarm](https://docs.docker.com/engine/swarm/services/#create-a-service-using-an-image-on-a-private-registry).

To add a new registry, the first thing we need to retrieve is the Certificate
Authority used to sign the DTR TLS Certificate. This can be done through DTR's
`/ca` endpoint.

```bash
$ curl -ks https://dtr.example.com/ca > dtr.crt
```

Next we need to convert this DTR certificate into a JSON configuration file,
this can then be used to register DTR within the 2nd Clusters Universal Control
Plane. 

A template of the json file called `dtr-bundle.json` is found below. Please
replace the host address with the relevant URL, and enter the contents of the
DTR CA certificate between the new line commands `\n and \n`.

> Note within the json file, ensure there are no line breaks between each line
> of the DTR CA certificate. 

```bash
$ cat dtr-bundle.json
{
  "hostAddress": "dtr.example.com",
  "caBundle": "-----BEGIN CERTIFICATE-----\n<contents of cert>\n-----END CERTIFICATE-----"
}
```

Now we will upload this configuration file to Cluster 2's Universal Control
Plane using the UCP API endpoint `/api/config/trustedregistry_`. To authenticate
against the API of Cluster 2's UCP, we have downloaded a [UCP client
bundle](/ee/ucp/user-access/cli/#download-client-certificates/), extracted it in
the current directory, and will reference the keys for authentication. 

```bash
$ curl --cacert ca.pem --cert cert.pem --key key.pem \
    -X POST \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d @dtr-bundle.json \
    https://cluster2.example.com/api/config/trustedregistry_
```

To check this has been imported successfully, as the UCP endpoint will not
output anything, we can check within Cluster 2's UCP UI. Select **Admin** in the
top left hand corner, select **Admin Settings** and the finally select **Docker
Trusted Registry**. If the registry has been added successfully we should see
the DTR listed. 

![](../../../images/remoteucp-addregistry.png){: .with-border}


You could also check the full [configuration
file](/ee/ucp/admin/configure/ucp-configuration-file/) within Cluster 2's UCP.
Once downloaded the `ucp-config.toml` file should now contain a section called
`[registries]`

```bash
$ curl --cacert ca.pem --cert cert.pem --key key.pem https://cluster2.example.com/api/ucp/config-toml > ucp-config.toml
```

If the new registry isn't shown in the list, please check the logs of the
`ucp-controller` container running on Cluster 2.

## Signing an image in DTR

We will now sign an image and push this to DTR, to sign images we need a User's
key pair from Cluster 2. Key pairs can be found in a client bundle, with the
`key.pem` being a private key, and `cert.pem` being the public key within a x509
certificate.

First we load the Private key into the local Docker trust store
`(~/.docker/trust)`. The name used here is purely metadata to help keep track of
which keys you have imported.

```
$ docker trust key load --name cluster2admin key.pem
Loading key from "key.pem"...
Enter passphrase for new cluster2admin key with ID a453196:
Repeat passphrase for new cluster2admin key with ID a453196:
Successfully imported key from key.pem
```

Next we will initiate the repository, and add the public key of Cluster 2's User
as a signer. You will be asked for a number of passphrases to protect the keys.
Please keep note of these passphrases, and to learn more about managing Keys
head to the Docker Content Trust documentation
[here](/engine/security/trust/trust_delegation/#managing-delegations-in-a-notary-server).


```
$ docker trust signer add --key cert.pem cluster2admin dtr.example.com/admin/trustdemo
Adding signer "cluster2admin" to dtr.example.com/admin/trustdemo...
Initializing signed repository for dtr.example.com/admin/trustdemo...
Enter passphrase for root key with ID 4a72d81:
Enter passphrase for new repository key with ID dd4460f:
Repeat passphrase for new repository key with ID dd4460f:
Successfully initialized "dtr.example.com/admin/trustdemo"
Successfully added signer: cluster2admin to dtr.example.com/admin/trustdemo
```

Finally we will sign an image tag. This pushes the image up to DTR, as well as
signing the tag with the User from Cluster 2's keys.

```
$ docker trust sign dtr.example.com/admin/trustdemo:1
Signing and pushing trust data for local image dtr.example.com/admin/trustdemo:1, may overwrite remote trust data
The push refers to repository [dtr.olly.dtcntr.net/admin/trustdemo]
27c0b07c1b33: Layer already exists
aa84c03b5202: Layer already exists
5f6acae4a5eb: Layer already exists
df64d3292fd6: Layer already exists
1: digest: sha256:37062e8984d3b8fde253eba1832bfb4367c51d9f05da8e581bd1296fc3fbf65f size: 1153
Signing and pushing trust metadata
Enter passphrase for cluster2admin key with ID a453196:
Successfully signed dtr.example.com/admin/trustdemo:1
```

Within the DTR UI, you should now be able to see a new tag has been pushed, as well as the **Signed** text next to the size. 

![](../../../images/remoteucp-signedimage.png){: .with-border}


We could sign this image multiple times if required, whether there were multiple
teams from the same Cluster that wanted to sign the image, or you wanted to
integrate DTR with more remote UCP's, and therefore a User from Cluster 1,
Cluster 2, Cluster 3...etc can all to sign the same image. 

## Enforce Signed Image Tags on the Remote UCP 

We can now enable **Only Run Signed Images** on the Remote UCP. To do this,
login to Cluster 2's UCP UI as an Admin user, select **Admin** in the top left
hand corner, select **Admin Settings** and then go down to **Docker Content
Trust**. 

For more information on **Only Run Signed Images** in UCP, refer to the [UCP
Documentation](/ee/ucp/admin/configure/run-only-the-images-you-trust/). 


![](../../../images/remoteucp-enablesigning.png){: .with-border}


Finally we are in a position to deploy a workload on Cluster 2, using a signed
image from a DTR running Cluster 1. This workload could be a simple `$ docker
run`, a Swarm Service or a Kubernetes workload. As a simple test, source a
client bundle, and try and run one of your signed images. 

```
$ source env.sh

$ docker service create dtr.example.com/admin/trustdemo:1
nqsph0n6lv9uzod4lapx0gwok
overall progress: 1 out of 1 tasks
1/1: running   [==================================================>]
verify: Service converged

$ docker service ls
ID                  NAME                    MODE                REPLICAS            IMAGE                                   PORTS
nqsph0n6lv9u        laughing_lamarr         replicated          1/1                 dtr.example.com/admin/trustdemo:1
```

## Troubleshooting

1) If the image is stored in a Private Repository within DTR, as there is no
   Single Sign On between Cluster 2 and DTR, you need to pass credentials to the
   Orchestrator. Please see the relevant
   [Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#create-a-secret-in-the-cluster-that-holds-your-authorization-token)
   or [Docker
   Swarm](https://docs.docker.com/engine/swarm/services/#create-a-service-using-an-image-on-a-private-registry)
   documentation.

2) If you see:

```
image or trust data does not exist for dtr.example.com/admin/trustdemo:1
```

This means something went wrong when initiating the repository or signing the
image, as the tag contains no signing data. 

3) If you see:

```
Error response from daemon: image did not meet required signing policy

dtr.example.com/admin/trustdemo:1: image did not meet required signing policy 
```

This means that the image was signed correctly, however the User that signed the
image does not mean the signing policy in Cluster 2. This could be because you
signed the Image with the wrong Users Keys.

## Where to go next

- [Learn more about Notary](/notary/advanced_usage.md)
- [Notary architecture](/notary/service_architecture.md)