'use strict';

import React, { Component, PropTypes } from 'react';
import styles from 'components/scenes/settings/formstyle.css';
// components
import RadioGroup from 'components/common/radioGroup';
import InputLabel from 'components/common/inputLabel';
import FileSystem from './filesystem';
import S3Settings from './s3';
import Azure from './azure';
import Swift from './swift';
// !IMPORTANT - @camacho 8/18/16 GCS
// DO NOT DELETE - THIS IS SHIPPING DARK TO BE TURNED ON IN >=2.2
// UNCOMMENT THE FOLLOWING IMPORT TO MAKE GCS INTERFACE VISIBLE
// import GCS from './gcs';
import ui from 'redux-ui';
// props
import { createStructuredSelector } from 'reselect';
import { getStorageType } from 'selectors/settings';
import { connect } from 'react-redux';

const mapState = createStructuredSelector({
    getStorageType
});

@connect(mapState)
@ui({
    state: {
        choices: [
            {
                label: 'Filesystem',
                value: 'filesystem'
            },
            {
                label: 'S3',
                value: 's3'
            },
            {
                label: 'Azure',
                value: 'azure'
            },
            {
                label: 'Swift',
                value: 'swift'
            }
            // !IMPORTANT - @camacho 8/18/16 GCS
            // DO NOT DELETE - THIS IS SHIPPING DARK TO BE TURNED ON IN >=2.2
            // UNCOMMENT THE FOLLOWING OBJECT TO MAKE GCS INTERFACE VISIBLE
            // {
            //   label: 'Google Cloud Storage',
            //   value: 'gcs'
            // }
        ],
        selected: (props) => props.getStorageType
    }
})
export default class ManualStorage extends Component {

    static propTypes = {
        ui: PropTypes.object,
        updateUI: PropTypes.func,
        onSubmit: PropTypes.func,
        getStorageType: PropTypes.string
    }

    chooseScene = (sceneChoice) => {
        this.props.updateUI({
            selected: sceneChoice
        });
    }

    render() {

        return (
            <div>
                <div className={ styles.formbox }>
                    <InputLabel>Storage backend</InputLabel>
                    <RadioGroup
                        initialChoice={ this.props.ui.selected }
                        choices={ this.props.ui.choices }
                        onChange={ ::this.chooseScene } />
                </div>
            { (() => {
                switch(this.props.ui.selected) {
                    case 's3': return <S3Settings saveFormData={ this.props.onSubmit } />;
                    case 'azure': return <Azure saveFormData={ this.props.onSubmit } />;
                    case 'swift': return <Swift saveFormData={ this.props.onSubmit } />;
                    // !IMPORTANT - @camacho 8/18/16 GCS
                    // DO NOT DELETE - THIS IS SHIPPING DARK TO BE TURNED ON IN >=2.2
                    // UNCOMMENT THE FOLLOWING CASE TO MAKE GCS INTERFACE VISIBLE
                    // case 'gcs': return <GCS saveFormData={ this.props.onSubmit } />;
                    default: return <FileSystem saveFormData={ this.props.onSubmit } />;
                }
            })() }

            </div>
        );
    }
}
