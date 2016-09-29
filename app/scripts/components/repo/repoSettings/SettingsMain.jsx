'use strict';

import React, { PropTypes, Component } from 'react';
import VisibilityForm from './settingsMain/VisibilityForm';
import DeleteRepositoryForm from './settingsMain/DeleteRepositoryForm';
import styles from './SettingsMain.css';
const { func, string, object } = PropTypes;
const debug = require('debug')('SettingsMain');

export default class repoSettingsMain extends Component {
  static propTypes = {
    JWT: string.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    user: object.isRequired
  }

  static contextTypes = {
    executeAction: func.isRequired
  }

  render() {
    const { JWT,
            name,
            user,
            namespace } = this.props;
    return (
      <div className={styles.settingsContent}>
        <div className='row'>
          <div className={'columns large-8 ' + styles.cards}>
            <VisibilityForm JWT={JWT}
                            name={name}
                            namespace={namespace}
                            user={user} />
            <DeleteRepositoryForm
                            JWT={JWT}
                            name={name}
                            params={this.props.params}/>
          </div>
        </div>
      </div>
    );
  }
}
