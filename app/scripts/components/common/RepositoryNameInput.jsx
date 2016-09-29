'use strict';

import React, { Component, PropTypes } from 'react';
import SimpleInput from './SimpleInput';
const {array, func, string} = PropTypes;

export default class RepositoryNameInput extends Component {
  static propTypes = {
    namespaces: array.isRequired,
    selectedNamespace: string.isRequired,
    repoName: string.isRequired,
    onRepoNameChange: func.isRequired,
    onNamespaceChange: func.isRequired,
    inputClass: string
  }

  render() {

    var namespaceOptions = this.props.namespaces.map(function(item, idx) {
      return (<option key={item + idx} value={item}>{item}</option>);
    });

    return (
      <div className='row'>
        <div className="small-4 columns">
          <select defaultValue={this.props.selectedNamespace}
                  onChange={this.props.onNamespaceChange}>
            {namespaceOptions}
          </select>
        </div>
        <div className="small-8 columns">
          <SimpleInput type="text"
                       onChange={this.props.onRepoNameChange}
                       value={this.props.repoName}
                       required
                       className={this.props.inputClass}
                       placeholder='Enter Name' />
        </div>
      </div>
    );
  }
}
