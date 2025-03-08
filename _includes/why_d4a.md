The [Docker for {{cloudprovider}}](/docker-for-{{cloudprovider | downcase}}/)
project was created and is being actively developed to ensure that Docker users
can enjoy a fantastic out-of-the-box experience on {{cloudprovider}}. It is now
generally available and can now be used by everyone.

As an informed user, you might be curious to know what this project offers
you for running your development, staging, or production workloads.

## Native to Docker

Docker for {{cloudprovider}} provides a Docker-native solution that avoids
operational complexity and adding unneeded additional APIs to the Docker stack.

Docker for {{cloudprovider}} allows you to interact with Docker directly
(including native Docker orchestration), instead of distracting you with the
need to navigate extra layers on top of Docker. You can focus instead on the
thing that matters most: running your workloads. This helps you and your
team to deliver more value to the business faster, to speak one common
"language", and to have fewer details to keep in your head at once.

The skills that you and your team have already learned, and continue to
learn, using Docker on the desktop or elsewhere automatically carry over to
using Docker on {{cloudprovider}}. The added consistency across clouds also
helps to ensure that a migration or multi-cloud strategy is easier to accomplish
in the future if desired.

## Skip the boilerplate and maintenance work

Docker for {{cloudprovider}} bootstraps all of the recommended infrastructure to
start using Docker on {{cloudprovider}} automatically. You don't need to worry
about rolling your own instances, security groups, or load balancers when using
Docker for {{cloudprovider}}.

Likewise, setting up and using Docker swarm mode functionality for container
orchestration is managed across the cluster's lifecycle when you use Docker for
{{cloudprovider}}. Docker has already coordinated the various bits of automation
you would otherwise be gluing together on your own to bootstrap Docker swarm
mode on these platforms. When the cluster is finished booting, you can jump
right in and start running `docker service` commands.

We also provide a prescriptive upgrade path that helps users upgrade between
various versions of Docker in a smooth and automatic way. Instead of
experiencing "maintenance dread" as you ponder your future responsibilities
upgrading the software you are using, you can easily upgrade to new versions
when they are released.

## Minimal, Docker-focused base

The custom Linux distribution used by Docker for {{cloudprovider}} is carefully
developed and configured to run Docker well. Everything from the kernel
configuration to the networking stack is customized to make it a favorable place
to run Docker. For instance, we make sure that the kernel versions are
compatible with the latest and greatest in Docker functionality, such as the
`overlay2` storage driver.

Instead of facing the trade-offs of a general purpose operating system, Docker's
custom Linux distribution focuses on only one thing: providing the best _Docker_
experience for you and your team.

## Self-cleaning and self-healing

Even the most conscientious admin can be caught off guard by issues such as
unexpectedly aggressive logging or the Linux kernel killing memory-hungry
processes. In Docker for {{cloudprovider}}, your cluster is resilient to a
variety of such issues by default.

Log rotation native to the host is configured for you automatically, so chatty
logs don't use up all of your disk space. Likewise, the "system prune" option
allows you to ensure unused Docker resources such as old images are cleaned up
automatically. The lifecycle of nodes is managed using auto-scaling groups or
similar constructs, so that if a node enters an unhealthy state for unforeseen
reasons, the node is taken out of load balancer rotation and/or replaced
automatically and all of its container tasks are rescheduled.

These self-cleaning and self-healing properties are enabled by default and don't
need configuration, so you can breathe easier as the risk of downtime is
reduced.

## Logging native to the platforms

Centralized logging is a critical component of many modern infrastructure
stacks. To have these logs indexed and searchable proves invaluable for
debugging application and system issues as they come up. Out of the box, Docker
for {{cloudprovider}} forwards logs from containers to a native cloud provider
abstraction ({{cloudprovider_log_dest}}).

## Next-generation Docker bug reporting tools

One common pain point in open source issue reporting is effectively
communicating the current state of your infrastructure and the issues you are
seeing to the upstream. In Docker for {{cloudprovider}}, you receive new tools
to communicate any issues you experience quickly and securely to Docker
employees. The Docker for {{cloudprovider}} shell includes a `docker-diagnose`
script which, at your request, transmits detailed diagnostic information to
Docker support staff to reduce the traditional
"please-post-the-output-of-this-command" back and forth frequently encountered
in bug reports.

# Try it today

Ready to get started? [Try Docker for {{cloudprovider}} today](/docker-for-{{cloudprovider | downcase}}/).
Search for existing issues, or create a new one, within the
[for {{cloudprovider}}](https://github.com/docker/for-{{cloudprovider | downcase}}) repository.
