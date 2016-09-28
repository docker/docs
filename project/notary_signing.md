# Notes on notary signing

Some misc non-UCP-specific notes that might be helpful are captured in
https://github.com/docker/notary/issues/779

## How it was set up

```sh
TAG=1.1.1
ORG=docker
# ...and dockerorcadev
for R in ucp ucp-controller ucp-proxy ucp-cfssl ucp-etcd ucp-dsinfo ucp-swarm ucp-compose ucp-auth ucp-auth-store ; do
    docker pull ${ORG}/${R}:${TAG}
    DOCKER_CONTENT_TRUST=1 docker push ${ORG}/${R}:${TAG} || break
done
```

Then a self-signed cert was generated for jenkins, and the public cert portion of that was added as a delegation signer

```sh
for R in ucp ucp-controller ucp-proxy ucp-cfssl ucp-etcd ucp-dsinfo ucp-swarm ucp-compose ucp-auth ucp-auth-store ; do
    notary delegation add docker.io/dockerorcadev/${R} targets/releases ./delegation.pem --all-paths && notary publish docker.io/dockerorcadev/${R}
done
```

## Dirty secrets

dhiltgen has a yubikey in his locking filing cabinet.  This holds **the**
root key for all UCP related repos.  A copy of that key is encrypted
on his laptop.  Each repo contains a unique targets key signed by this
root key.  (encrypted copies also on his laptop)  The delegation cert
was then signed by each of these target kets (see above)

At present, notary doesn't support multiple target keys per repo.
Once this feature is added, then one (or more) other person on the team
(ehazlett) would use the yubikey on their laptop to push to all the repos
repos.  We would then be rid of the SPOF.

For now though, the SPOF is OK, because jenkins pushes using the
delegation key.
