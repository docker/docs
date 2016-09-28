'use strict';

import React, { Component, PropTypes } from 'react';
import styles from 'components/scenes/settings/formstyle.css';
import InputLabel from 'components/common/inputLabel';
// state
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { getStorageYAML } from 'selectors/settings';
// actions
import {
  saveYamlStorage,
  getStorageSettings
} from 'actions/settings';
import { mapActions } from 'utils';

const mapState = createStructuredSelector({
  yaml: getStorageYAML
});

@connect(mapState, mapActions({ saveYamlStorage, getStorageSettings }))
export default class YAMLFile extends Component {
  static propTypes = {
    actions: PropTypes.object,
    yaml: PropTypes.string.isRequired
  }

  downloadYAML() {
    const a = document.createElement('a');
    const b = new Blob([this.props.yaml], { type: 'application/octet-binary' });
    const u = URL.createObjectURL(b);
    a.href = u;
    a.download = 'storage.yml';
    a.click();
    URL.revokeObjectURL(u);
  }

  upload() {
    const { files } = this.refs.file;
    if (files.length < 1) {
      return;
    }
    let reader = new FileReader();
    reader.onload = this.onload.bind(this);
    reader.readAsText(files[0]);
  }

  onload(evt) {
    // extract the yaml from the filereader event
    const yaml = evt.target.result;
    this.props.actions.saveYamlStorage({
      config: yaml
    }).then(() => this.props.actions.getStorageSettings());
  }

    render() {
        return (
            <div className={ styles.formbox }>
                <h2>YAML file</h2>
                <InputLabel>Upload</InputLabel>
                <input type='file' className={ styles.fileInput } ref='file' accept='.yml,.yaml,.txt' />

                <button onClick={ ::this.upload }>Upload</button>

                <InputLabel>Download</InputLabel>
                <button onClick={ ::this.downloadYAML }>Download</button>
            </div>
        );
    }
}
