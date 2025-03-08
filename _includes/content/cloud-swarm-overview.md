You can now create _new_ Docker Swarms from within Docker Cloud as well as
register existing swarms.

When you create a swarm, Docker Cloud connects to the Cloud provider on your
behalf, and uses the provider's APIs and a provider-specific template to launch
Docker instances. The instances are then joined to a swarm and the swarm is
configured using your input. When you access the swarm from Docker Cloud, the
system forwards your commands directly to the Docker instances running in the
swarm.
