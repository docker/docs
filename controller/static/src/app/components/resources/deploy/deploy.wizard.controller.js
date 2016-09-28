'use strict';

var angular = require('angular');
var $ = require('jquery');
var _ = require('lodash');

DeployWizardController.$inject = ['permissions', 'networks', 'volumes', 'AuthService', 'ServiceService', 'MessageService', '$scope', '$state', '$location', '$timeout', '$anchorScroll'];
function DeployWizardController(permissions, networks, volumes, AuthService, ServiceService, MessageService, $scope, $state, $location, $timeout, $anchorScroll) {
  var vm = this;
  vm.admin = AuthService.getCurrentUser().admin;
  vm.permissions = _.chain(permissions)
    .filter(function(p) {
      // show all labels for admins
      if(vm.admin) {
        return true;
      }

      // otherwise only show restricted/full control labels
      if(p.role < 2) {
        return false;
      } else {
        return true;
      }
    })
    .map(function(p) { return p.label; })
    .value();

  vm.swarmNetworks = networks.filter(defaultNetworkFilter);
  vm.swarmVolumes = volumes;
  vm.swarmNetworkMap = _.keyBy(vm.swarmNetworks, 'Id');

  vm.serviceMode = 'Replicated';
  vm.serviceScale = 1;

  vm.ports = [];
  vm.publishPort = null;
  vm.targetPort = null;
  vm.publishPortProtocol = 'TCP';

  vm.accessLabel = '';
  vm.labels = [];
  vm.labelName = '';
  vm.labelValue = '';

  vm.constraints = [];
  vm.currentConstraint = '';

  vm.envVars = [];
  vm.variableName = '';
  vm.variableValue = '';

  vm.volumes = [];
  vm.volumeType = 'volume';
  vm.volumeSource = '';
  vm.volumeTarget = '';
  vm.volumeWritable = 'Write';
  vm.updateConfigDelay = 0;
  vm.stopGracePeriod = 0;

  vm.networks = [];

  vm.service = {
    Name: '',

    Labels: {},

    TaskTemplate: {
      ContainerSpec: {
        Image: '',
        Command: [],
        Args: [],
        Env: [],
        Dir: '',
        User: '',
        Mounts: [],

        Resources: {
          Limits: {},
          Reservations: {}
        }
      },
      Placement: {
        Constraints: []
      }
    },

    Networks: [],

    EndpointSpec: {
      Ports: []
    },

    Mode: {},

    RestartPolicy: {},

    UpdateConfig: {
      Parallelism: 1,
      Delay: null
    }
  };

  vm.scrollTo = scrollTo;
  vm.deploy = deploy;
  vm.pushVolume = pushVolume;
  vm.removeVolume = removeVolume;
  vm.pushEnvVar = pushEnvVar;
  vm.removeEnvVar = removeEnvVar;
  vm.pushLabel = pushLabel;
  vm.removeLabel = removeLabel;
  vm.pushConstraint = pushConstraint;
  vm.removeConstraint = removeConstraint;
  vm.pushPort = pushPort;
  vm.removePort = removePort;
  vm.pushNetwork = pushNetwork;
  vm.removeNetwork = removeNetwork;

  function defaultNetworkFilter(network) {
    return !network.Name.endsWith('/none') &&
      !network.Name.endsWith('/bridge') &&
      !network.Name.endsWith('/host');
  }

  function pushNetwork() {
    vm.networks.push(vm.network);
    vm.network = '';
  }

  function removeNetwork(network) {
    var index = vm.networks.indexOf(network);
    vm.networks.splice(index, 1);
  }

  function pushVolume() {
    var volume = {
      Type: vm.volumeType,
      Target: vm.volumeTarget,
      Source: vm.volumeSource,
      Writable: vm.volumeWritable
    };
    vm.volumes.push(volume);
    vm.volumeSource = '';
    vm.volumeTarget = '';
    vm.volumeWritable = 'Write';

    // Reset dropdown
    $timeout(function() {
      if(vm.volumeType === 'volume') {
        $('.volume-name.dropdown').dropdown('clear');
      }
    });
  }

  function removeVolume(volume) {
    var index = vm.volumes.indexOf(volume);
    vm.volumes.splice(index, 1);
  }

  function pushPort() {
    var port = {
      Protocol: vm.publishPortProtocol,
      PublishedPort: vm.publishPort,
      TargetPort: vm.targetPort
    };
    vm.ports.push(port);
    vm.publishPortProtocol = 'TCP';
    vm.publishPort = null;
    vm.targetPort = null;
  }

  function removePort(port) {
    var index = vm.ports.indexOf(port);
    vm.ports.splice(index, 1);
  }

  function pushEnvVar() {
    var envVar = {
      name: vm.variableName,
      value: vm.variableValue
    };
    vm.envVars.push(envVar);
    vm.variableName = '';
    vm.variableValue = '';
  }

  function removeEnvVar(envVar) {
    var index = vm.envVars.indexOf(envVar);
    vm.envVars.splice(index, 1);
  }

  function pushLabel() {
    var label = {
      name: vm.labelName,
      value: vm.labelValue
    };
    vm.labels.push(label);
    vm.labelName = '';
    vm.labelValue = '';
  }

  function removeLabel(label) {
    var index = vm.labels.indexOf(label);
    vm.labels.splice(index, 1);
  }

  function pushConstraint() {
    vm.constraints.push(vm.currentConstraint);
    vm.currentConstraint = '';
  }

  function removeConstraint(c) {
    var index = vm.constraints.indexOf(c);
    vm.constraints.splice(index, 1);
  }

  function transformEnvVars() {
    var i;
    if (vm.variableName.length > 0) {
      vm.service.TaskTemplate.ContainerSpec.Env.push(vm.variableName + '=' + vm.variableValue);
    }
    for (i = 0; i < vm.envVars.length; i++) {
      vm.service.TaskTemplate.ContainerSpec.Env.push(vm.envVars[i].name + '=' + vm.envVars[i].value);
    }
  }

  function transformNetworks() {
    if (vm.network && vm.network !== '') {
      pushNetwork();
    }

    for (var i = 0; i < vm.networks.length; i++) {
      vm.service.Networks.push({
        Target: vm.networks[i]
      });
    }
  }

  function transformLabels() {
    if (vm.labelName.length > 0) {
      vm.service.Labels[vm.labelName] = vm.labelValue;
    }
    for (var i = 0; i < vm.labels.length; i++) {
      vm.service.Labels[vm.labels[i].name] = vm.labels[i].value;
    }
  }

  function transformConstraints() {
    for (var i = 0; i < vm.constraints.length; i++) {
      vm.service.TaskTemplate.Placement.Constraints = vm.constraints;
    }
  }

  function transformCommand() {
    if (vm.command && vm.command.length > 0) {
      vm.service.TaskTemplate.ContainerSpec.Command = vm.command.split(' ');
    }
    if (vm.args && vm.args.length > 0) {
      vm.service.TaskTemplate.ContainerSpec.Args = vm.args.split(' ');
    }
  }

  function transformVolumes() {
    if (vm.volumeType && vm.volumeSource !== '' && vm.volumeTarget !== '') {
      pushVolume();
    }
    for (var i = 0; i < vm.volumes.length; i++) {
      vm.service.TaskTemplate.ContainerSpec.Mounts.push(vm.volumes[i]);
    }
  }

  function transformPorts() {
    if (vm.publishPort) {
      // this is used in case there is just a single port and the
      // '+' has not been clicked to push the port onto the array
      pushPort();
    }

    for (var i = 0; i < vm.ports.length; i++) {
      var port = vm.ports[i];
      var portSpec = {
        Protocol: port.Protocol.toLowerCase(),
        TargetPort: port.TargetPort,
        PublishedPort: port.PublishedPort
      };

      vm.service.EndpointSpec.Ports.push(portSpec);
    }
  }

  function transformUpdateConfig() {
    // Convert seconds delay spec to nanos
    vm.service.UpdateConfig.Delay = vm.updateConfigDelay * 1000000000;
    vm.service.TaskTemplate.ContainerSpec.StopGracePeriod = vm.stopGracePeriod * 1000000000;
  }

  function isFormValid() {
    return $('#deploy .ui.form').form('validate form');
  }

  function scrollTo(id) {
    if (id === 'deploy-environment') {
      vm.environmentVarsVisible = true;
    } else if (id === 'deploy-network') {
      vm.networkConfigVisible = true;
    } else if (id === 'deploy-volumes') {
      vm.volumeConfigVisible = true;
    } else if (id === 'deploy-labels') {
      vm.labelsVisible = true;
    } else if (id === 'deploy-constraints') {
      vm.constraintsVisible = true;
    } else if (id === 'deploy-updateconfig') {
      vm.updateConfigVisible = true;
    }

    $location.hash(id);
    $anchorScroll();
  }

  function transformMode() {
      if(vm.serviceMode === 'Replicated') {
        vm.service.Mode = {
          Replicated: {
            Replicas: vm.serviceScale
          }
        };
      } else if(vm.serviceMode === 'Global') {
        vm.service.Mode = {
          Global: {}
        };
      }
  }

  function deploy() {
    transformMode();
    transformVolumes();
    transformUpdateConfig();
    transformEnvVars();
    transformCommand();
    transformPorts();
    transformLabels();
    transformConstraints();
    transformNetworks();

    ServiceService.create(vm.service)
      .then(function(data) {
        $state.go('dashboard.resources.services.inspect', {
          id: data.ID
        });
        MessageService.addSuccessMessage('Successfully created service ' + vm.service.Name);
      }, function(error) {
        MessageService.addErrorMessage('Error creating service', error.data.message);
      });
  }

}

module.exports = DeployWizardController;
