### Typical Docker workflow

<span style="font-size: 122%; color: #5D6D7E; font-weight: 150;">&#9312; Get your code and its dependencies into Docker [containers](engine/getstarted/step_two.md).</span>

<span style="font-size: 110%; color: #5D6D7E; font-weight: 150; display: block; padding-left: 2em;">- [Write a Dockerfile](engine/getstarted/step_four.md) that
defines the execution environment and pulls in your code.</span>

<span style="font-size: 110%; color: #5D6D7E; font-weight: 150; display: block; padding-left: 2em;">- If your app depends on external services (such as Redis or MySQL), [find them on a registry like Docker Hub](docker-hub/repos.md), and refer to them in [a Docker Compose file](compose/overview.md), along with a call to your app, so they'll run simultaneously.</span>

<span style="font-size: 110%; color: #5D6D7E; font-weight: 150; display: block; padding-left: 2em;">-  Software providers also distribute paid software on the [Docker Store](https://store.docker.com).</span>

<span style="font-size: 110%; color: #5D6D7E; font-weight: 150; display: block; padding-left: 2em;">-  Build, then run your containers on a virtual host with [Docker Machine](machine/overview.md) as you develop.</span>
<p />
<span style="font-size: 122%; color: #5D6D7E; font-weight: 150;">&#9313; Configure [networking](engine/tutorials/networkingcontainers.md) and [storage](engine/tutorials/dockervolumes.md) for your solution, if needed.</span>

<span style="font-size: 122%; color: #5D6D7E; font-weight: 150;">&#9314; Upload builds to a registry ([ours](/engine/getstarted/step_six.md) or [yours](/datacenter/dtr/2.0/index.md)) or your cloud providers to collaborate with your team.</span>

<span style="font-size: 122%; color: #5D6D7E; font-weight: 150;">&#9315;
To run your app as a set of services across multiple hosts, [set up a Swarm cluster](/engine/swarm/index.md) and
[scale it to meet demand](/engine/swarm/swarm-tutorial/scale-service.md). Use [Universal Control Plane](/datacenter/ucp/1.1/overview.md) to manage your swarm in a friendly UI!</span>

<span style="font-size: 122%; color: #5D6D7E; font-weight: 150;">&#9316;
Deploy to your preferred cloud providers with [Docker Cloud](/docker-cloud/index.md), or use [Docker Datacenter](https://www.docker.com/products/docker-datacenter) to deploy to your own on-premise hardware.</span>
