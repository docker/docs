'use strict';

import React, {
  PropTypes,
  createClass
  } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import BuildSettings from './builds/BuildSettings.jsx';
import RepoPushTrigger from './builds/RepoPushTrigger.jsx';
import LinkAutoTrigger from './builds/LinkAutoTrigger.jsx';
import BuildTrigger from './builds/BuildTrigger.jsx';
import DeployKeys from './builds/DeployKeys.jsx';
import AutoBuildSettingsStore from '../../../stores/AutoBuildSettingsStore.js';
import FA from '../../common/FontAwesome.jsx';
import styles from './Builds.css';

const { bool, func, number, object, shape, string } = PropTypes;
const debug = require('debug')('AutoBuildSettings');

var AutoBuildSettings = createClass({
  displayName: 'AutoBuildSettings',
  contextTypes: {
    executeAction: func.isRequired
  },
  propTypes: {
    autoBuildStore: shape({
      deployKey: string,
      provider: string.isRequired,
      repo_web_url: string.isRequired
    }),
    autoBuildBlankSlate: object.isRequired,
    description: string.isRequired,
    fullDescription: string.isRequired,
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    status: number.isRequired,
    JWT: string.isRequired,
    validations: object.isRequired
  },
  render() {
    debug('Build Settings');
    const {
      autoBuildBlankSlate,
      autoBuildStore,
      autoBuildLinks,
      JWT,
      name,
      namespace,
      newTags,
      triggerLinkForm,
      triggerLogs,
      triggerStatus,
      validations,
      STATUS
      } = this.props;

    let deployKeysPage;
    let deployKey = autoBuildStore.deploykey;
    if (deployKey) {
      deployKeysPage = (
        <DeployKeys autoBuildStore={autoBuildStore} />
      );
    }

    return (
      <div className={styles.buildsContent}>
        <div className='row'>
          <div className='columns large-12'>
            <BuildSettings autoBuildStore={autoBuildStore}
                           JWT={JWT}
                           name={name}
                           namespace={namespace}/>
          </div>
        </div>
        <div className='row'>
          <div className='columns large-12'>
            <RepoPushTrigger autoBuildStore={autoBuildStore}
                             newTags={newTags}
                             autoBuildBlankSlate={autoBuildBlankSlate}
                             JWT={JWT}
                             name={name}
                             namespace={namespace}
                             validations={validations}
                             STATUS={STATUS}/>
          </div>
        </div>
        <div className='row'>
          <div className='columns large-12'>
            <LinkAutoTrigger autoBuildStore={autoBuildStore}
                             autoBuildLinks={autoBuildLinks}
                             namespace={namespace}
                             name={name}
                             JWT={JWT}
                             validations={validations}
                             triggerLinkForm={triggerLinkForm}/>
          </div>
        </div>
        <div className='row'>
          <div className='columns large-12'>
            <BuildTrigger autoBuildStore={autoBuildStore}
                          triggerStatus={triggerStatus}
                          triggerLogs={triggerLogs}
                          namespace={namespace}
                          name={name}
                          JWT={JWT}/>
          </div>
        </div>
        <div className='row'>
          <div className='columns large-12'>
            {deployKeysPage}
          </div>
        </div>
      </div>
    );
  }
});

export default connectToStores(AutoBuildSettings,
  [ AutoBuildSettingsStore ],
  function({ getStore }, props) {
    return getStore(AutoBuildSettingsStore).getState();
  });
