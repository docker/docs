'use strict';

import React, { Component, PropTypes } from 'react';
import { mkAvatarForNamespace } from 'utils/avatar';
const { shape, string } = PropTypes;

import styles from './Gravatar.css';

export default class GravatarValue extends Component {

  static propTypes = {
    placeholder: string,
    value: shape({
      value: string.isRequired,
      label: string.isRequired
    })
  }

  render() {
    const { value: selectedObject, placeholder } = this.props;
    let gravatarValue = placeholder;
    if (selectedObject) {
      gravatarValue = (
        <div>
          <img className={styles.gravatar} src={mkAvatarForNamespace(selectedObject.value)}/>
          {selectedObject.label}
        </div>
      );
    }
    return (
      <div className="Select-placeholder">
        { gravatarValue }
      </div>
    );
  }

}
