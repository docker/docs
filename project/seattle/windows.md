# Windows Spec

We believe that windows agent node support will be largely "free"
in this release, given the chosen architecture.  Windows agent nodes
may be joined using the standard `docker cluster join ...` commands.
Once joined, the cluster will report those nodes.

We will ensure our node APIs (inherited from core) expose the node type,
and this is rendered in the UI.  Additional refinements can be taken on
an as needed basis as we approach DockerCon.
