# install Command

USAGE:
```
   docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        dockerorca/ucp \
        install [OPTIONS]
 ```
DESCRIPTION:

Install the UCP controller on a machine. You can only install on machines where
Docker Engine is already installed. If you intend to install a multi-node
cluster, you must open firewall ports between the engines for the following
ports:

* 443
* 12376
* 12379 through 12382
* 2376 (system default '--swarm-port')

You can optionally use an externally generated and signed certificate for the
UCP controller by specifying '--external-ucp-ca'.  Create a storage volume named
'ucp-server-certs' with ca.pem, cert.pem, and key.pem in the root directory
before running the install.


OPTIONS:
`--debug`, `-D	`			      
  Enable debug
`--jsonlog`				        
  Produce json formatted output for easier parsing
`--interactive`, `-i`			   
  Enable interactive mode. The system prompts you to enter all required information.
`--fresh-install`			     
   Destroy any existing state and start fresh
`--san` A subject alternative sames (SAN) for certs.  You can specify multiple SAN's, for example, `--san foo1.bar.com --san foo2.bar.com`.
`--host-address` 			
  Specify the visible IP/hostname for this node. (override automatic detection) [$UCP_HOST_ADDRESS]
`--old-kernel`				
    Install on older kernels (some features may not be supported)
`--image-version "latest"`		
    Select a specific UCP version
--swarm-port "2376"			
  Select what port to run the local Swarm manager on
--external-ucp-ca			
  Set up UCP with an external CA.
--preserve-certs			
  Don't (re)generate certs on the host if existing ones are found
--binpack				
  Set Swarm scheduler to binpack mode (default spread)
--random				
  Set Swarm scheduler to random mode (default spread)
--pull "missing"			
  Specify image pull behavior ('always', when 'missing', or 'never')
