#!/usr/bin/env python2.7

import os.path
import subprocess
import sys
import time
import argparse
import base64

import boto.ec2


BASE_AMI_ID = "ami-d05e75b8"
DTR_AMI_SECURITY_GROUP = "sg-d4c6bcb3"
SSH_KEY_PATH = "~/.ssh/id_rsa"
CONFIG_SCRIPT = os.path.join(os.path.dirname(sys.argv[0]), "configure.sh")

# usage examples:
# ./generate_ami.py cs-engine
# ./generate_ami.py hourly --dtr_version=tar --src=./dtr.tar --image_name=dtr_1.3.2_hourly
# ./generate_ami.py byol --dtr_version=tar --src=./dtr.tar --image_name=dtr_1.3.2_byol
# ./generate_ami.py byol --dtr_version=1.3.2
# ./generate_ami.py byol --dtr_version=tar --src=/home/v/dev/dtr/go/src/github.com/docker/dhe-deploy/dtr1.3.2.tar --image_name=dtr_1.3.2_byol2

def main():
    parser = argparse.ArgumentParser(description='Build an amazon AMI for the marketplace.')
    parser.add_argument('image_type', type=str, choices=['byol', 'cs-engine', 'hourly'])
    parser.add_argument('--image_name', type=str)
    parser.add_argument('--channel', type=str)
    parser.add_argument('--dtr_version', type=str)  # 'tar' or version number
    parser.add_argument('--awsinit', type=str)
    parser.add_argument('--ssh_key_name', type=str)
    parser.add_argument('--hub_password', type=str)
    parser.add_argument('--src', type=str)  # used with --dtr_version=tar
    parser.set_defaults(
        dtr_version='1.3.1',
        image_name='',
        awsinit='./awsinit.sh',
        channel='production',
        ssh_key_name='viktor'
    )
    args = parser.parse_args()

    if args.image_name == '':
        if args.image_type == 'cs-engine':
            args.image_name = 'docker_cs_engine'
        else:
            args.image_name = ('dtr_%s_%s' % (args.dtr_version, args.image_type))
    if args.image_name != 'cs-engine' and args.awsinit == '':
        print 'Path awsinit script required.'
        sys.exit(1)

    docker_creds = ''
    if args.channel == 'dev':
        print "DO NOT RELEASE THIS AMI!!! THIS AMI CONTAINS PRIVATE HUB CREDENTIALS. THIS SHOULD BE USED ONLY FOR UPGRADE TESTING"
        username = 'dockerautomation'
        password = args.hub_password
        docker_creds = base64.encodestring('{}:{}'.format(username, password)).strip()

    conn = boto.ec2.connect_to_region("us-east-1")
    print "Creating instance..."
    reservation = conn.run_instances(
        BASE_AMI_ID,
        key_name=args.ssh_key_name,
        instance_type="m3.medium",
        security_group_ids=[DTR_AMI_SECURITY_GROUP],
    )
    instance = reservation.instances[0]
    print "Waiting for instance creation..."
    while instance.ip_address is None or len(instance.ip_address) == 0:
        time.sleep(5)
        instances = conn.get_only_instances(instance_ids=[instance.id])
        instance = instances[0]

    print "Waiting for instance to come up at %s..." % (instance.ip_address)
    # wait for the server to come up
    while True:
        p = subprocess.Popen("ssh -oStrictHostKeyChecking=no -oConnectTimeout=5 -i %s ubuntu@%s /bin/true" % (SSH_KEY_PATH, instance.ip_address), shell=True, stdout=sys.stdout, stderr=sys.stderr)
        if p.wait() == 0:
            break
        else:
            time.sleep(5)

    # copy tar over if necessary
    if args.dtr_version == 'tar':
        print "Copying tar..."
        p = subprocess.Popen("scp -oStrictHostKeyChecking=no -i %s %s ubuntu@%s:/tmp/dtr.tar"
            % (SSH_KEY_PATH, args.src, instance.ip_address), shell=True, stdout=sys.stdout, stderr=sys.stderr)
        ret = p.wait()
        if ret != 0:
            print "Failed to copy tar to instance"
            sys.exit(ret)

    # copy awsinit over
    if args.image_type != 'cs-engine':
        print "Copying awsinit..."
        p = subprocess.Popen("scp -oStrictHostKeyChecking=no -i %s %s ubuntu@%s:/tmp/awsinit.sh"
            % (SSH_KEY_PATH, args.awsinit, instance.ip_address), shell=True, stdout=sys.stdout, stderr=sys.stderr)
        ret = p.wait()
        if ret != 0:
            print "Failed to copy awsinit.sh to instance"
            sys.exit(ret)

    p = subprocess.Popen("ssh -oStrictHostKeyChecking=no -i %s ubuntu@%s 'DTR_VERSION=%s DTR_CHANNEL=%s IMAGE_TYPE=%s DOCKER_CREDS=%s bash -s' < %s"
        % (SSH_KEY_PATH, instance.ip_address, args.dtr_version, args.channel, args.image_type, docker_creds, CONFIG_SCRIPT), shell=True, stdout=sys.stdout, stderr=sys.stderr)
    ret = p.wait()
    if ret != 0:
        print "Configuration script failed with return code: %d" % ret
        sys.exit(ret)

    print "Creating image..."
    ami_id = instance.create_image(args.image_name)
    images = conn.get_all_images(image_ids=[ami_id])
    ami = images[0]

    print "Waiting for image..."
    while ami.state != "available":
        print ".",
        sys.stdout.flush()
        time.sleep(5)
        images = conn.get_all_images(image_ids=[ami_id])
        ami = images[0]
    instance.terminate()
    print "New AMI saved as %s" % ami.id

if __name__ == "__main__":
    main()
