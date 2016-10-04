Docker provides a way to run applications securely isolated in a container, packaged 
with all its dependencies and libraries. Because your application can always be
run with the environment it expects right in the build image, testing
and deployment is simpler than ever, as your build will be fully portable and ready
to run as designed in any environment. And because containers are lightweight and run
without the extra load of a hypervisor, you can run many applications that all rely
on different libraries and environments on a single kernel, each one never interfering
with the other. This allows you to get more out of your hardware by shifting the "unit
of scale" for your application from a virtual or physical machine, to a container 
instance.
