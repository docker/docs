+++
title = "Remove an image and garbage collection"
description = "Remove an image and garbage collection"
keywords = ["docker, documentation, about, technology, hub, registry, soft deletion, cron job, garbage collection, enterprise"]
[menu.main]
parent="smn_dhe"
weight=12
+++


# Overview

This document describes the two-step process of removing an image from the
Trusted Registry. This process is first performed by users wanting to remove
their images then, an administrator removes those images mainly through a
recurring job.

## Images and manifests

Manifests describe images. They reference a list of layers and hold metadata
about the image. Manifests can share layers. Tags point to manifests. You can
reference an image by tag or directly by the hash of the manifest. If you
purposefully delete one of those manifests and the image layers referenced by
that manifest become orphaned, then they can be removed during the garbage
collection job. This occurs if they are not referenced by other manifests.

**Note**: Since many developers may use a base image for future images, it is possible that there will be image layers that may never be deleted. There might be other manifests that point to layers of the base image which could still be used by others.

## Prerequisites
You need an image to remove.

## Soft delete an image

Developers may want to remove an image for several reasons.
Examples include:

* the image is outdated and there is a new version
* the image is no longer to be used

If a developer removes the manifest of an image, then it is called a
soft deletion. Access to the image has been removed although physically, it is
still in the repository.

Soft deletion occurs through the following API endpoint `DELETE`. From the command line type:

`curl -u <username>:<password> -X DELETE https://<DTR HOST>/api/v0/repositories/<namespace>/<reponame>/manifests/<reference>`

You can only delete one image at a time and you must also be authenticated as a
user who has "write" level access to the repository.

## Garbage collection

Administrators may want to optimize registry storage and free up space. They can
set up a recurring job where the system searches for any layers that are not
referenced by any manifests and removes them. It can be either performed immediately, or more realistically, on a periodic automatic schedule. System
administrators can perform garbage collection of unreferenced layers tags in two ways:

* Immediately, by typing the following:

`curl -X POST -u <username>:<password> -H "Content-Type: application/json" "https://api/v0/admin/jobs" -d "{ \"job\" : \"registryGC\"}"`

* Routinely, through a cron job, as you can see in the following example:

`curl -u <username>:<password> -H 'Content-Type: application/json' -X POST https://<DTR
HOST>/api/v0/admin/settings/registry/garbageCollection/schedule -d '{"schedule":
"<schedule>"}'`

See the [cron package](https://godoc.org/github.com/robfig/cron) for other
examples depending on your needs. An example schedule: `0 0 3 * * SAT` (every
Saturday at 3 AM).

You can also set the cron job through the Trusted Registry UI. Note that it is
not as granular as setting a cron job from the command line. Navigate to
Settings > Garbage collection and schedule your job.

Your downtime depends on the number of images and/or layers that are to be
deleted. Docker recommends performing garbage collection weekly during off time.
While garbage collection is occurring, anyone who tries to push an image will
get an error message.

## See the results of the garbage collection

View the results by running the following example in a command line:

`curl -u <username>:<password> https://<DTR
HOST>/api/v0/admin/settings/registry/garbageCollection/lastSavings`

### See also

* [**Administrator Guide**](adminguide.md) Go here if you are an administrator
responsible for running and maintaining Docker Trusted Registry.

* [**Configuration**](configuration.md) Go here to find out details about
setting up and configuring Docker Trusted Registry for your particular
environment.
