#!/bin/bash

set -u
set -e
set -x

echo 'Upgrading system packages...'
sudo apt-get upgrade -y

echo 'Installing linux kernel extras for AUFS...'
sudo apt-get install linux-image-extra-$(uname -r) -qq

echo 'Installing CS Engine...'
curl -s 'https://pgp.mit.edu/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import
echo "deb https://packages.docker.com/1.9/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
sudo apt-get update
sudo apt-get install -y docker-engine

if [ $IMAGE_TYPE != "cs-engine" ]; then
    if [ $IMAGE_TYPE == "hourly" ]; then
        # This needs to happen before the install because the installer starts DTR
        echo 'Setting to hourly aws billing'
        sudo mkdir -p /usr/local/etc/dtr/
        echo HourlyAWS | sudo tee /usr/local/etc/dtr/.hourly > /dev/null
    fi

    if [ ${DTR_VERSION} == "tar" ]; then
        sudo docker load -i /tmp/dtr.tar
        rm /tmp/dtr.tar
        if [ ${DTR_CHANNEL} == "dev" ]; then
            mkdir -p /home/ubuntu/.docker
            cat <<EOF > /home/ubuntu/.docker/config.json
{
	"https://index.docker.io/v1/": {
		"auth": "${DOCKER_CREDS}",
		"email": "a@a.a"
	}
}
EOF
            cp /home/ubuntu/.docker/config.json /home/ubuntu/.dockercfg
            sudo bash -c "$(sudo docker run dockerhubenterprise/trusted-registry-dev install)"
            rm -rf /home/ubuntu/.docker
            rm -f /home/ubuntu/.dockercfg
        else
            sudo bash -c "$(sudo docker run docker/trusted-registry install)"
        fi
    else
        echo 'Installing DTR version '${DTR_VERSION:-latest}'...'
        sudo bash -c "$(sudo docker run docker/trusted-registry:${DTR_VERSION:-latest} install)"
    fi
fi

echo 'Installing ec2-ami-tools...'
sudo sed -i.dist 's/universe$/universe multiverse/' /etc/apt/sources.list
sudo apt-get update
sudo apt-get install -y ec2-ami-tools

echo 'Update AMI Tools on startup...'
sudo tee -a /etc/rc.local <<-EOM
# Update the Amazon EC2 AMI tools
echo " + Updating EC2 AMI tools"
apt-get upgrade -y ec2-ami-tools
echo " + Updated EC2 AMI tools"
EOM

echo 'Disabling root login...'
sudo sed -i 's/^.*PermitRootLogin\(.*\)$/#PermitRootLogin\1/' /etc/ssh/sshd_config
echo 'PermitRootLogin without-password' | sudo tee -a /etc/ssh/sshd_config

echo 'Disabling root password...'
sudo passwd -l root

echo 'Removing SSH Host Key Pairs...'
sudo shred -u /etc/ssh/*_key /etc/ssh/*_key.pub

echo 'Removing SSH Authorized Keys...'
shred -u ~/.ssh/authorized_keys
sudo shred -u /root/.ssh/authorized_keys || true

echo 'Setting default username to ec2-user...'
sudo sed -i 's/^\(\s*\)name:.*$/\1name: ec2-user/' /etc/cloud/cloud.cfg

if [ $IMAGE_TYPE != "cs-engine" ]; then
    # hacky script that set the admin password when an ami is provisioned
    echo 'runcmd:' | sudo tee --append /etc/cloud/cloud.cfg
    echo ' - [ /usr/local/etc/dtr/awsinit.sh ]' | sudo tee --append /etc/cloud/cloud.cfg
    sudo cp /tmp/awsinit.sh /usr/local/etc/dtr/awsinit.sh
    sudo chmod +x /usr/local/etc/dtr/awsinit.sh
    rm /tmp/awsinit.sh
fi

# This is optional
echo 'Disabling sshd DNS Checks...'
sudo sed -i 's/^.*UseDNS\(.*\)$/#UseDNS\1/' /etc/ssh/sshd_config
echo 'UseDNS no' | sudo tee -a /etc/ssh/sshd_config

# Very important!
echo 'Wiping history...'
sudo rm /var/log/auth.log
history -w
shred -u ~/.*history
history -c
