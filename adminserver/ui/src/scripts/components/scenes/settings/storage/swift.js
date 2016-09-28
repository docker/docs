'use strict';

import React, { Component, PropTypes } from 'react';
import styles from 'components/scenes/settings/formstyle.css';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import Button from 'components/common/button';
import { ToggleWithLabel } from 'components/common/toggleSwitch';
import ui from 'redux-ui';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { getStorageRegistry } from 'selectors/settings';
import { checkRequiredFields, reportMissingFields } from 'validation';

const mapState = createStructuredSelector({
    storage: getStorageRegistry
});

@connect(mapState)
@ui({
    state: {
        showAdvanced: false
    }
})
@reduxForm({
    form: 'swiftstorage',
    fields: [
        'authurl',
        'username',
        'password',
        'container',
        'tenant',
        'domain',
        'region',
        'rootdirectory',
        'chunksize',
        'insecureskipverify'
    ],
    validate: (values) => {
        const requiredNames = [
            'authurl',
            'username',
            'password',
            'container'
        ];
        return reportMissingFields(checkRequiredFields(requiredNames, values), {});
    }
}, (state, props) => ({
    initialValues: {
        ...props.storage.swift
    }
}))
export default class Swift extends Component {

    static propTypes = {
        ui: PropTypes.object,
        updateUI: PropTypes.func,
        fields: PropTypes.object,
        handleSubmit: PropTypes.func,
        saveFormData: PropTypes.func,
        storage: PropTypes.object
    }

    toggleAdvanced = (evt) => {
        evt.preventDefault();
        this.props.updateUI({ showAdvanced: !this.props.ui.showAdvanced });
    }

    onSubmit = (data) => {

        // chunksize is a number
        if (data.chunksize) {
            data.chunksize = parseInt(data.chunksize);
        }

        this.props.saveFormData(data, 'swift');
    }

    render() {

        const {
            fields: {
                authurl,
                username,
                password,
                container,
                tenant,
                domain,
                region,
                rootdirectory,
                chunksize,
                insecureskipverify
            },
            handleSubmit,
            storage
        } = this.props;

        return (
            <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) }>
                <div className={ styles.formbox }>
                    <h2>Swift settings</h2>
                    <InputLabel
                        hint='URL to get an auth token'>Authorization URL</InputLabel>
                    <Input type='text' formfield={ authurl } />
                    <div className={ styles.row }>
                    <div className={ styles.halfColumn }>
                        <InputLabel
                            hint='Your OpenStack username'>Username</InputLabel>
                    <Input type='text' formfield={ username } />
                    </div>
                    <div className={ styles.halfColumn }>
                        <InputLabel
                            hint='Your OpenStack password'>Password</InputLabel>
                    <Input type='password' formfield={ password } />
                    </div>
                    </div>
                    <InputLabel
                        hint='Name of your Swift container where you want to store objects'>Container</InputLabel>
                    <Input type='text' formfield={ container } />

                <p><a href='#' onClick={ ::this.toggleAdvanced }>{ this.props.ui.showAdvanced ? 'Hide' : 'Show' } advanced settings</a></p>

                { this.props.ui.showAdvanced &&
                  <div>
                    <InputLabel
                        isOptional={ true }
                        hint='Your OpenStack Tenant name or Tenant ID'>Tenant name or ID</InputLabel>
                    <input type='text' {...tenant} />
                    <InputLabel
                        isOptional={ true }
                        hint='Your OpenStack domain name or ID for identity v3 API'>Domain name or ID</InputLabel>
                    <input type='text' {...domain} />
                    <InputLabel
                        isOptional={ true }
                        hint='OpenStack region name where you want to store objects'>Region</InputLabel>
                    <input type='text' {...region} />
                    <InputLabel
                        isOptional={ true }
                        tip="Defaults to the empty string which is the container's root."
                        hint='Root directory where you want to store registry files'>Prefix</InputLabel>
                    <input type='text' {...rootdirectory} />
                    <InputLabel
                        isOptional={ true }
                        tip='Default is 5 MB. You might experience better performance for larger chunk sizes depending on the speed of your connection to Swift.'
                        hint='Segment size for Dynamic Large Objects uploads (performed by WriteStream)'>Chunk size</InputLabel>
                    <input type='number' {...chunksize} />

                    <div style={
                        { float: 'left' }
                    }>
                        <ToggleWithLabel labelOptions={
                            {
                                labelText: 'Skip TLS Verification',
                                inline: true,
                                optional: true,
                                tip: 'Enable to skip TLS verification for your OpenStack provider. The driver disables this by default.'
                            }
                        }
                        initial={ storage.swift ? storage.swift.insecureskipverify : false }
                        formField={ insecureskipverify } />
                    </div>
                  </div> }
                </div>
                <Button variant='primary'>Save</Button>
            </form>
        );
    }
}
