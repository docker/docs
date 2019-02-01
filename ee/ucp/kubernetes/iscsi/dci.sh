#!/bin/sh

# Prior to executing this script, user should have
# setup `dci secrets` specific to their cloud provider.
# I used AWS and setup access keys prior.

export DCI_CLOUD=aws
export DCI_DEPLOYMENT=core-storage9
export DCI_REPOSITORY=dockereng

# To override the dci container image, use this option
# For most cases, setting this to ${DCI_CLOUD}-local is sufficient.
export DCI_TAG=aws-9857f6f

# Initialize cluster
dci cluster init

# Install a specific development version of UCP
# Perform a few other cluster specific configs.
dci cluster config set use_dev_version true
dci cluster config set docker_ucp_image_repository dockereng
dci cluster config set docker_ucp_version "3.2.0-latest"
dci cluster config set docker_ucp_storage_driver iscsi

# Preserve install containers. Ansible is case sensitive. 'False', not 'false'
dci cluster config set docker_remove_containers False

# set log level to debug.
# Use this to test --iscsiadm-path and --iscsidb-path
dci cluster config set docker_ucp_install_args "'--debug'"

dci cluster config set region us-west-2
dci cluster config set linux_ucp_worker_count 2
dci cluster config set linux_dtr_count 0

# Apply presets
dci cluster apply-preset "RHEL 7.5"

dci cluster provision

# Post provision installs.
dci cluster ssh "sudo yum install -y iscsi-initiator-utils"
dci cluster ssh "sudo modprobe iscsi_tcp"
# master doubles up as iscsi target. So let firewalld allow iscsi traffic.
# Restarting firewalld before installing dockerd is least intrusive;
# else iptable rules get messed up
dci cluster ssh "sudo yum install -y firewalld"
dci cluster ssh "sudo systemctl enable firewalld"
dci cluster ssh "sudo systemctl start firewalld"
dci cluster ssh "sudo firewall-cmd --add-service=iscsi-target --permanent"
dci cluster ssh "sudo firewall-cmd --add-port=18700/tcp --permanent"
dci cluster ssh "sudo firewall-cmd --reload"

# Install cluster. This is essential
# Apply = provision + install
dci cluster install --log-level debug
