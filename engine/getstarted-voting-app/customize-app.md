---
description: Customize the voting app
keywords: multi-container, services, swarm mode, cluster, voting app, docker-stack.yml, docker stack deploy
title: Customize the app and redeploy
---


In this step, we'll make a simple change to the application and redeploy it.
We'll change the focus of the poll from Cats or Dogs to .NET or Java.

In the real world, you might do this by making code changes and rebuilding to
create new images, or new versions of the same images.

For our example, we've done this for you. We already have a different version of
the application available as built images. So, all you need to do is reconfigure
`docker-stack.yml` to point to the new images, then redeploy.

## Update docker-stack.yml to reference new images

Go back to `docker-stack.yml` and replace the `before` tags on both the `vote` and `result` images to reference `after` tags.

![before tags for vote and result in yml](images/customize-before.png)

![after tags for vote and result in yml](images/customize-after.png)

## Redeploy

Run the same deploy command again to run the app with the new configuration.

```
docker stack deploy --compose-file docker-stack.yml vote
```

The output will look similar to this:

```
docker@manager:~$ docker stack deploy --compose-file docker-stack.yml vote
Updating service vote_redis (id: md73fohylg8q85aryz07852o0)
Updating service vote_db (id: gny9ieqxancnufrg1oeazz9gq)
Updating service vote_vote (id: 0ig0s4gb10q8auek9tris5i8z)
Updating service vote_result (id: lqwxjmmdhmegs2aw0a6ehipsp)
Updating service vote_worker (id: 8u4cfu60dtliz77x1o74kiwpr)
Updating service vote_visualizer (id: ya2vt9z2b4to248tccjjeqitw)
```

## Try it out

Take the app for another test drive.

You'll see the new voting choices for Java and .NET at `<MANAGER-IP:>5000`.

![New voting page](images/vote-2.png)

And the related results at `<MANAGER-IP:>5001`.

![New results page](images/vote-results-2.png)

The visualizer at  `<MANAGER-IP:>8080` will show some differences, such as
updates to the containers, and some services might have moved between the
manager and the worker.

However, the PostgreSQL container (`vote_db`) and the
visualizer (`vote_visualizer`) will still be running on the manager because of
the `[node.role == manager]` constraints on the those services, which we
did not change.

![New visualizer web page](images/visualizer-2.png)

## Resources

The source code for the voting app and a deeper dive example walkthrough will be made available at [Docker Labs](https://github.com/docker/labs).

The lab walkthrough provides more technical detail on deployment configuration, Compose file keys, the application and stack and services concepts, and networking.

The images you used are available on Docker Hub. When you ran `docker-stack.yml` you pulled those images from Docker Hub.

For information on how to build your own images, see:

*   [Build your own image](/engine/getstarted/step_four.md) in the [Get Started with Docker tutorial](/engine/getstarted/index.md)

*   [Build your own images](/engine/tutorials/dockerimages.md) in the [Learn by example](/engine/tutorials/index.md) tutorials

*   [Docker Labs](https://github.com/docker/labs)

For more about `docker-stack.yml` and the `docker stack deploy` command, see [deploy](/compose/compose-file.md#deploy) in the [Compose file reference](/compose/compose-file.md).

To learn more about swarm mode, start with the [Swarm mode overview](/engine/swarm/overview.md).
