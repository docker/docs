'use strict';

import React, { createClass, PropTypes } from 'react';
const { bool, shape, string } = PropTypes;
import _ from 'lodash';
import Button from '@dux/element-button';
import FA from '../common/FontAwesome.jsx';
import connectToStores from 'fluxible-addons-react/connectToStores';
import triggerBuild from '../../actions/triggerBuild.js';
import TriggerBuildStore from '../../stores/TriggerBuildStore.js';
import styles from './TriggerButton.css';
var debug = require('debug')('TriggerButton');

var TriggerButton = createClass({
  displayName: 'TriggerButton',
  propTypes: {
    isAutomated: bool.isRequired,
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    abtrigger: shape({
      hasError: bool.isRequired,
      success: bool.isRequired
    })
  },
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _triggerBuild(e){
    e.preventDefault();
    this.context.executeAction(triggerBuild, {
      JWT: this.props.JWT,
      name: this.props.name,
      namespace: this.props.namespace
    });
  },
  render(){
    const {
      isAutomated,
      abtrigger
      } = this.props;

    let trigButton = null;
    if (isAutomated) {
      let triggerIntent = 'warning';
      let triggerText = ' Trigger a Build';
      if (abtrigger.success) {
        triggerIntent = 'success';
        triggerText = 'Success';
      } else if (abtrigger.hasError) {
        triggerIntent = 'alert';
        triggerText = 'Error Occurred';
      }
      trigButton = (
        <Button variant={triggerIntent} size='small' onClick={this._triggerBuild}>
          <FA icon='fa-play'/>{triggerText}
        </Button>
      );
    }
    return trigButton;
  }
});

export default connectToStores(TriggerButton,
  [ TriggerBuildStore ],
  function({ getStore }, props) {
    return getStore(TriggerBuildStore).getState();
  });
