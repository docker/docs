/* eslint-disable max-len */

/* Evaluate DDC on */

export const evaluationInstructions = {
  toInstallDDC: {
    1: 'Get Docker Toolbox',
    2: 'Download license into your project directory',
    3: 'Run the following command in the same directory',
  },
  pullCommand: '$ docker run alexmavr/ddc-in-a-box | bash',
  scriptTiming: 'The script will take approximately 15-20 mins depending on your internet connectivity.',
};

// Right side Additional Information
export const additionalInformation = [
  {
    title: 'System Requirements',
    bulletPoints: [
      '5GB hard drive space',
      '2GB of RAM',
      'Internet connectivity (downloads ~4GB)',
      'docker-machine & bash (included in Docker Toolbox)',
      'Your license file downloaded to a directory on your machine',
    ],
  }, {
    title: 'What this container does',
    bulletPoints: [
      'Creates a VM using docker-machine',
      'Installs Universal Control Plane 1.1',
      'Installs Trusted Registry 2.0',
      'Configures authentication in both systems with a default user',
      'Configures licensing in both products',
      'Configures UCP to trust the registry service in DTR',
    ],
  }, {
    title: 'Advanced Options',
    bulletPoints: [
      'You can set 3 environment variables for this script, related to the use of docker-machine:',
    ],
    environmentVariables: [
      {
        variable: 'MACHINE_DRIVER',
        description: 'The default is \'virtualbox\'',
      }, {
        variable: 'MACHINE_DRIVER_FLAGS',
        description: 'The default is "--virtualbox-memory 2048 --virtualbox-disk-size 16000"',
      }, {
        variable: 'MACHINE_NAME',
        description: 'The default is = \'ddc-eval\'',
      },
    ],
  },
];

/* Deploy on */

export const AWS = {
  label: 'AWS using Quickstart',
  url: 'https://console.aws.amazon.com/cloudformation/home?#/stacks/new?stackName=DockerDatacenter&templateURL=https://s3-us-west-2.amazonaws.com/ddc-on-aws-public/ddc_on_aws.json',
};

export const Azure = {
  label: 'Microsoft Azure from the Marketplace ',
  url: 'https://azure.microsoft.com/en-us/marketplace/partners/docker/dockerdatacenterdocker-datacenter/',
};

export const Linux = {
  label: 'Linux',
  url: 'https://docs.docker.com/docker-trusted-registry/cs-engine/install/',
};

// Links in the right column

export const UCPGuide = {
  label: 'Guide for Universal Control Plane',
  url: 'https://docs.docker.com/ucp/install-sandbox/',
};

export const DTRGuide = {
  label: 'Guide for Docker Trusted Registry',
  url: 'https://docs.docker.com/docker-trusted-registry/install/install-dtr/',
};

export const prodInstall = {
  label: 'Planning a Production Ready Install',
  url: 'https://docs.docker.com/ucp/installation/plan-production-install/',
};
