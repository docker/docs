'use strict';

var angular = require('angular');

var nodeService = require('./node.service');
var nodesController = require('./nodes.controller');
var inspectNodeController = require('./inspect.node.controller');

require('./nodes.css');

module.exports = angular.module('ducp.nodes', [])
  .factory('NodeService', nodeService)
  .controller('NodesController', nodesController)
  .controller('InspectNodeController', inspectNodeController);
