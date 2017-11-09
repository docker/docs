---
description: Explains next steps after the tour
keywords: Kitematic, open source, contribute, contributor, tour, development, contribute
title: Where to learn more
---

You've just created your first pull request to Kitematic!

## Take the development challenge

Now that you’ve had some practice, try adding another feature to Kitematic on your own.

As you learned in the previous exercise, adding the container ID to the container **Settings** tab is fairly simple.

Let's provide another missing piece of information for Kitematic users:

"_When I look at an active container in Kitematic, I want to know what command the container is currently running_."

![An active container in Kitematic](images/kitematic_gui_container_id.png)

In a terminal window, users can get this by looking at running containers with: `docker ps`

As an exercise, implement the code changes needed to display the current container's running command. When you are ready to share the new mini feature, create a PR for it.

## Where to go next

- To learn more about contributing to open source projects, see
[Contribute to the Moby project](https://github.com/moby/moby/blob/master/CONTRIBUTING.md).

- To learn more about contributing to Docker product documentation, see the [README on docker/docker.github.io](https://github.com/docker/docker.github.io/blob/master/README.md)
