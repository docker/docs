# UCP specific swarm

We now build swarm as part of the UCP CI system, from our own private forks of the project.

CI Jobs:

* https://jenkins-orca.dockerproject.com/job/swarm-classic/

Github trees:

* https://github.com/docker/ucp-swarm -- Private fork of docker/swarm
* https://github.com/docker/ucp-swarm-library-image -- The "build" process that produces minimal swarm images


## Process when tracking upstream exactly

Our strategy is to keep master of the `ucp-swarm` tree identical to
`docker/swarm`, and inherit the same tags.  If we need to deviate, we'll
do so on a branch, and use UCP specific version tags.  Deviations should
be kept to a minimum.

1. Create a local checkout of `docker/swarm` and add the `ucp-swarm` tree above as a remote
2. `git fetch --all`
3. Push `docker/swarm`'s master to `ucp-swarm`
4. Push any updated tags
5. Kick off a `swarm-classic` build on the (new) tags
6. Update the `Dockerfile` in this directory to use that new version tag (without the preceeding "v")


## Process when deviating

Try not to do this!

Details TBD
