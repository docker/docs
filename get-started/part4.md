---
title: "Getting Started, Part 4: Scaling Your App on a Cluster"
---

In [Getting Started, Part 3: Stateful, Multi-container Applications](part3.md),
we figured out how to relate containers to each other. We organized an
application into two simple services -- a frontend and a backend -- and defined
how they are linked together.

In Part 4, we are going to take this application, which all ran on one **host**
(a virtual or physical machine), and deploy it onto a cluster of hosts, just
like we would in production.

## Understanding Swarm clusters

Up until now you have been using Docker in a single-host mode on your local
machine, which allows the client, which is the command-line interface (CLI), to
make assumptions about how to operate. Namely, the client assumes that the
Docker Daemon is running on the same host as the client. Single-host operations
could also be done on remote machines with your client.

But Docker also can be switched into "swarm mode." A swarm is a group of hosts
that are running Docker and have been joined into a cluster. After that has
happened, you continue to run the Docker commands you're used to, but now the
concept of a "host" changes from a single virtual or physical machine, to a
swarm. And, "a single virtual or physical machine" is not referred to as a host,
it's called a node -- or, a computing resource inside your cluster.

{% capture local-instructions %}
We now have two VMs created, named `myvm1` and `myvm2`. These
machines are running boot2docker, a light-weight Linux distribution, as an
operating system.

Finally, you'll assemble these machines into a swarm. The first one will act as
the manager, which executes `docker` commands and authenticates workers to join
the swarm, and the second VM will act as a worker.

You can send a single command to one of your VMs using `docker-machine ssh`. In
this case, we're telling `myvm1` to become a swarm manager with
`docker swarm init`:

```node
$ docker-machine ssh myvm1 "docker swarm init"
Swarm initialized: current node <node ID> is now a manager.

To add a worker to this swarm, run the following command:

  docker swarm join \
  --token <token> \
  <ip>:<port>
```

As you can see, the response to `docker swarm init` contains a pre-configured
`docker swarm join` command for you to run on any nodes you want to add. This
command will join `myvm2` to your new swarm as a worker. Copy this command,
and send it to `myvm2`:

```
$ docker-machine ssh myvm2 "docker swarm join \
--token <token> \
<ip>:<port>"

This node joined a swarm as a worker.
```

> **Note**: You can also run `docker-machine ssh myvm2` with no command attached
to open a terminal session on that VM. Type `exit` when you're ready to return
to the host shell prompt.

Congratulations, you created your first swarm!

Here are some commands you might like to run to interact with your swarm a bit:

```
docker-machine ssh myvm1 "docker node ls" # list the nodes in your swarm
docker-machine ssh myvm1 "docker node inspect <node ID>" # inspect a node
docker-machine ssh myvm1 "docker swarm join-token -q worker" # view join token
docker-machine ssh myvm2 "docker swarm leave" # make the worker leave the swarm
docker-machine stop $(docker-machine ls -q) # stop all running VMs
docker-machine rm $(docker-machine ls -q) # destroy all VMs
```
{% endcapture %}

{% capture aws %}
#### Amazon Web Services with Docker Cloud

The easiest way to demonstrate all this is to use Docker Cloud, which manages
clusters that you run on popular cloud providers, like Heroku, Amazon Web
Services (AWS), and so on. Because AWS has a free tier of service, which lets
you provision low-resource virtual machines for free, we're going to use that
to learn these concepts. We're also not going to be using any of Docker Cloud's
paid features, so let's dive in and deploy something!

### Sign up for AWS, and configure it

All we have to do to let Docker Cloud manage nodes for us on free-tier AWS is
create a service policy that grants certain permissions, and apply that to an
identity called a "role," using AWS's Identity and Access Management (IAM) tool.

-   Go to [aws.amazon.com](https://aws.amazon.com) and sign up for an account. It's free.
-   Go to [the IAM panel](https://console.aws.amazon.com/iam/home#policies)
-   Click **Create Policy**, then **Create Your Own Policy**.
-   Name the policy `dockercloud-policy` and paste the following text in the
    space provided for **Policy Document**, then click **Create Policy**.

    ```json
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Action": [
            "ec2:*",
            "iam:ListInstanceProfiles"
          ],
          "Effect": "Allow",
          "Resource": "*"
        }
      ]
    }
    ```
-   Now [create a role](https://console.aws.amazon.com/iam/home#roles) with a name
    of your choice.
-   Select **Role for Cross-Account Access**, and in the submenu that opens select **Allows IAM users from a 3rd party AWS account to access this account**.
-   In the **Account ID** field, enter the ID for the Docker Cloud service: `689684103426`.
-   In the **External ID** field, enter your Docker Cloud username.
-   On the next screen, select the `dockercloud-policy` you created to attach to the role.
-   On next page review your entries and copy the full **Role ARN** string. The
    ARN string should look something like `arn:aws:iam::123456789123:role/dockercloud-role`. You'll use the ARN in the next step.
-   Finally, click **Create Role**.

And you've done it! Your AWS account will allow Docker Cloud to control
virtual machines, if we configure Docker Cloud to use the role you've created.
So, let's do that now.

> Note: If you had any trouble along the way, there are more detailed
  [instructions in the Docker Cloud docs](/docker-cloud/infrastructure/link-aws.md).
  If you'd like to use a cloud provider besides AWS, check out
  [the list](/docker-cloud/infrastructure/index.md). We're just using AWS here
  because you don't have to pay.

### Configure Docker Cloud to manage to your AWS instances

- Go to [cloud.docker.com](http://cloud.docker.com) and sign in with the
  same Docker ID you used in [part 2](/getting-started/part2.md).
- Click **Settings**, and in the Cloud Providers section, click the plug icon.
- Enter the Role ARN string you copied earlier, e.g. `arn:aws:iam::123456789123:role/dockercloud-role`.
- Click **Save**.

And now, Docker Cloud can create and manage instances for you, and turn them
into a swarm.

## Creating your first Swarm cluster

1.  Go back to Docker Cloud by visiting [cloud.docker.com](https://cloud.docker.com).
2.  Click **Node Clusters** in the left navigation, then click the **Create** button.
    This pulls up a form where you can create our cluster.
3.  Leave everything default, except:
    - Name: Give your cluster a name
    - Region: Select a region that's close to you
    - Provider: Set to "Amazon Web Services"
    - Type/Size: Select the `t2.nano` option as that is free-tier
4.  Launch the cluster by clicking **Launch node cluster**; this will spin
    up a free-tier Amazon instance.
5.  Now, click **Services** in the left navigation, then the **Create** button,
    then the **globe icon**.
6.  Search Docker Hub for the image you uploaded
{% endcapture %}

{% capture local %}
#### VMs running on your local machine (Mac, Linux, Windows 7 and 8)

(Virtualbox stuff)

{{ local-instructions | markdownify }}
{% endcapture %}

{% capture localwin %}
#### VMs running on your local machine (Windows)

Set up a virtual switch for your VMs to use, so they will be able to connect
to each other.

1. Launch Hyper-V Manager
2. Click **Virtual Switch Manager** in the right-hand menu
3. Click **Create Virtual Switch** of type **External**
4. Give it the name `myswitch`, and check the box to share your host machine's
   active network adapter

Now, create a couple of virtual machines using our node management tool,
`docker-machine`:

```none
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm1
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm2
```

{{ local-instructions | markdownify}}
{% endcapture %}

## Set up your swarm

A swarm is made up of multiple nodes, which can be either physical or virtual
machines. The basic concept is simple enough: run `docker swarm init` to make
your current machine a manager node, and run `docker swarm join` on other
machines to have them join the swarm as a worker. Choose a tab below to see how
this plays out in various contexts.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#aws">Amazon Web Services</a></li>
  <li><a data-toggle="tab" href="#local">Local VMs (Mac, Linux, Windows 7 and 8)</a></li>
  <li><a data-toggle="tab" href="#localwin">Local VMs (Windows 10/Hyper-V)</a></li>
</ul>
<div class="tab-content">
  <div id="aws" class="tab-pane fade active in">{{ aws | markdownify }}</div>
  <div id="local" class="tab-pane fade">{{ local | markdownify }}</div>
  <div id="localwin" class="tab-pane fade">{{ localwin | markdownify }}</div>
</div>

[On to next >>](part5.md){: class="button darkblue-btn"}
