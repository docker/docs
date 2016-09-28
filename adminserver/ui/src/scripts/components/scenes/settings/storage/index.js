'use strict';

import React, { Component, PropTypes } from 'react';
import RadioGroup from 'components/common/radioGroup';
import ManualStorage from './manualStorage';
import YAMLFile from './yaml.js';
import ui from 'redux-ui';
import Spinner from 'components/common/spinner';
// actions
import autoaction from 'autoaction';
import { connect } from 'react-redux';
import { mapActions } from 'utils';
import {
  saveFormStorage,
  getStorageSettings
} from 'actions/settings';
// selectors
import { getNonStorageRegistry } from 'selectors/settings';
import { createStructuredSelector } from 'reselect';
import consts from 'consts';

const mapState = createStructuredSelector({
    getNonStorageRegistry
});

@autoaction({
    getStorageSettings: []
}, {
    getStorageSettings
})
@connect(
    mapState,
    mapActions({
        saveFormStorage
    })
)
@ui({
    state: {
        choices: [
            {
                label: 'Manual form',
                value: 'MANUAL'
            },
            {
                label: 'YAML file',
                value: 'YAML'
            }
        ],
        selected: 'MANUAL'
    }
})
export default class StorageSettings extends Component {

    static propTypes = {
        ui: PropTypes.object,
        updateUI: PropTypes.func,
        getNonStorageRegistry: PropTypes.object,
        actions: PropTypes.object
    };

    chooseScene = (sceneChoice) => {
        this.props.updateUI({ selected: sceneChoice });
    };

    onSubmit = (data, form) => {
        const settings = {
            ...this.props.getNonStorageRegistry,
            [form]: data
        };
        this.props.actions.saveFormStorage(settings);
    };

    render() {
        const status = [
            [consts.settings.GET_STORAGE_SETTINGS]
        ];

        return (
        <div>
          <Spinner loadingStatus={ status }>
            <h2>Set up your storage with a
                <RadioGroup
                    initialChoice={ this.props.ui.selected }
                    choices={ this.props.ui.choices }
                    onChange={ ::this.chooseScene } />
            </h2>
            { (() => {
                switch(this.props.ui.selected) {
                    case 'YAML': return <YAMLFile />;
                    default: return <ManualStorage onSubmit={ ::this.onSubmit } />;
                }
            })() }
          </Spinner>
        </div>
    );
  }
}
