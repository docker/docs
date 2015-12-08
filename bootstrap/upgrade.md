# upgrade

docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      dockerorca/ucp \
      upgrade [OPTIONS]

DESCRIPTION:

Upgrades UCP on a node. When upgrading UCP, you must run
the 'upgrade' command against every server and node running Docker Engine in
your cluster. You should upgrade your controller and replica nodes first,
followed by your compute nodes.  

Depending on the upgrade path for your version of UCP, you may also need to upgrade Docker Engine. You can also choose to upgrade Docker Engine simply because a new version is available.  Always upgrade a nodes Docker Engine installation before upgrading its UCP installation.

After upgrading each node, confirm the node is present in the UCP console
before proceeding to the next node.

OPTIONS:
 --debug, -D			enable debug
 --jsonlog			Produce json formatted output for easier parsing
 --interactive, -i		Enable interactive mode.  You will be prompted to enter all required information.
 --image-version "latest"	Select a specific UCP version
 --id 			The ID of the UCP instance to upgrade
 --pull "always"		Specify image pull behavior ('always', when 'missing', or 'never')
