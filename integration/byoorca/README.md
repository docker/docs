# Bring-your-own Orca

The tests in this directory are set up to run against an already deployed
UCP instance.  In general, interesting test biz logic should be written
in ../utils/ so it can be re-used in different test scenarios

To run these integration tests, download a bundle from UCP, source
the env and run the script.  The tests will verify that `$DOCKER_HOST`
points to a UCP instance and skip if it doesn't.
