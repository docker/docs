# Orca

Docker Orchestration


# Getting started
Check out [docs/evaluation-install.md](docs/evaluation-install.md) for quickstart installation instructions.


Read on for more details on developing Orca.

# Typical developer flow

## Build

    ./script/run make build image

If your docker engine can volume mount your local directory, you can use the incremental run script to speed up builds:

    ./script/run_inc make all image

## Copy images

If you use machine (or your own hand-rolled VM) there are some utility scripts in `./script` to make your life easier

    ./script/copy_orca_images_machine <machine name>


## Install
Hint: most developers use `docker-machine` to install UCP on a secondary VM
so its easy to wipe the slate clean periodically.

    docker run --rm -it --name ucp -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:$(make print-TAG) \
        install --swarm-port 3376

## Unit tests

    ./script/run make test

## Integration Tests
These tests are also Go based, but typically rely on
`docker-machine`.  For more details on how to use these, check out
[project/integration_tests.md](project/integration_tests.md)


# Vendored Dependencies

Take a look at [project/godeps.md](project/godeps.md)
