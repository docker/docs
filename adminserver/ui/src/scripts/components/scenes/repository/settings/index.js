'use strict';

import React, { Component } from 'react';
import styles from '../repository.css';
import css from 'react-css-modules';

import EditForm from './edit';
import DeleteForm from './delete';

@css(styles, {allowMultiple: true})
export class RepositorySettingsTab extends Component {
  static propTypes = {
    params: React.PropTypes.object
  }

  render() {
    return (
      <div styleName='repositorySettings'>
        <EditForm { ...this.props } />
        <DeleteForm { ...this.props } />
      </div>
    );
  }
}
