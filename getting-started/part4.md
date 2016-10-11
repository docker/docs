---
title: "Getting Started, Part 4: Scaling Your App on a Cluster"
---

# Part 4: Scaling Your App on a Cluster

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

## Before we get started: signup and configuration

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
  same Docker ID you used in [step 2](/getting-started/step2.md).
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


[On to next >>](part5.md){: class="button darkblue-btn"}
