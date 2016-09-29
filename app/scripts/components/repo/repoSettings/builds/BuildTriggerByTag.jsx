'use strict';

import React, { PropTypes, Component } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';
import FA from 'common/FontAwesome.jsx';
import AutobuildTriggerByTagStore from 'stores/AutobuildTriggerByTagStore';
import triggerBuildByTagAction from 'actions/triggerBuildByTag.js';
const debug = require('debug')('BuildTriggerByTag');
import findIndex from 'lodash/array/findIndex';
const { array, bool, number, func, shape, string } = PropTypes;
import {STATUS as COMMONSTATUS} from 'stores/common/Constants';
const {
  ATTEMPTING,
  DEFAULT
} = COMMONSTATUS;

import styles from './BuildTriggerByTag.css';

class BuildTriggerByTag extends Component {
  constructor(props) {
    super(props);
  }

  static contextTypes = {
    executeAction: func.isRequired
  };

  //TODO: arrayOfType
  static propTypes = {
    isRegexRule: bool.isRequired,
    isNew: bool.isRequired,
    repoInfo: shape({
      JWT: string.isRequired,
      name: string.isRequired,
      namespace: string.isRequired
    }),
    buildTag: shape({
      source_name: string.isRequired,
      source_type: string.isRequired,
      dockerfile_location: string.isRequired,
      name: string.isRequired
    }),
    triggerByTagStore: shape({
      triggers: array,
      tagStatuses: array
    }),
    triggerId: number.isRequired
  };

  triggerBuildByTag = (e) => {
    e.preventDefault();
    const {dockerfile_location, source_name, source_type} = this.props.buildTag;
    const {JWT, name, namespace} = this.props.repoInfo;
    const {triggerId} = this.props;
    if (source_name) {
      let buildTag = {
        dockerfileLocation: dockerfile_location,
        name: name,
        namespace: namespace,
        sourceType: source_type,
        sourceName: source_name
      };
      this.context.executeAction(triggerBuildByTagAction, {JWT, tag: buildTag, triggerId});
    }
  };

  render() {

    const {
      isNew,
      isRegexRule,
      triggerId
    } = this.props;

    const {
      triggers,
      tagStatuses
    } = this.props.triggerByTagStore;

    const statusIndex = findIndex(tagStatuses, (s) => {
      return s.id === triggerId;
    });
    const triggerIndex = findIndex(triggers, (t) => {
      return t.id === triggerId;
    });

    let buttonVariant = 'primary',
        buttonText = 'Trigger',
        buttonTooltip = 'Manually trigger a build',
        triggerIcon = 'fa-play-circle-o';

    if (isNew) {
      buttonTooltip = 'Cannot trigger on new rules.';
    } else if (isRegexRule) {
      buttonTooltip = 'Manual trigger not available on regex based rules.';
    }

    const currentStatus = statusIndex > -1 ? tagStatuses[statusIndex].status : DEFAULT;
    const triggerSuccess = triggerIndex > -1 ? triggers[triggerIndex].success : '';
    const triggerError = triggerIndex > -1 ? triggers[triggerIndex].error : '';

    if (currentStatus === ATTEMPTING) {
      triggerIcon = 'fa-spin fa-spinner';
      buttonText = 'Triggering';
    }

    if (triggerSuccess) {
      buttonVariant = 'success';
      buttonText = 'Triggered';
      triggerIcon = 'fa-check-circle-o';
    } else if (triggerError) {
      buttonVariant = 'alert';
      buttonText = 'Failed';
      triggerIcon = 'fa-times-circle-o';
    }

    return (
      <span title={buttonTooltip}
            className={isNew || isRegexRule ? styles.invisible : styles.visible}>
        <Button ghost
                variant={buttonVariant}
                onClick={this.triggerBuildByTag}>
          <FA icon={triggerIcon} /> {buttonText}
        </Button>
      </span>
    );
  }
}

export default connectToStores(BuildTriggerByTag,
  [
    AutobuildTriggerByTagStore
  ],
  function({ getStore }, props) {
    return { triggerByTagStore: getStore(AutobuildTriggerByTagStore).getState() };
  });
