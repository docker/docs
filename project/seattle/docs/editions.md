# Docker AWS/Azure Editions

Notes:

* Azure editions not yet available (unknown if it will be in time for beta)
* Current setup uses a distinct cloud formation template.  By GA there will be only one.


## User story

As an ops person, I want to deploy UCP on AWS or Azure using docker
editions.

## Launch DDC on AWS Editions

* Click on [<img src=https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png>](https://console.aws.amazon.com/cloudformation/home?#/stacks/new?stackName=Docker&templateURL=https://s3.amazonaws.com/docker-for-aws/aws/alpha/aws-v1.12.0-rc3-beta1-ddc.json)
* Login to your AWS account
* On the "Select Template" page click "next"
* On the "Specify Details" pgae, fill out the number of managers, and worker nodes you want.  We recommend m3.medium or larger manager instance types.  Don't forget to select an SSH key, and select "Install Docker Data Center"
* No options are required
* Review the settings, and confirm to create

Note: The creation takes a while.  The cloudformation template will
claim it is done before UCP has finished installing (known issue today,
hopefully will be fixed before GA)

You can click on the "Stack" in the list of stacks, and expand the
"Outputs" section.  Once the deployment has finished, you'll be able to
see how to SSH in as well as login to the UCP web console.
