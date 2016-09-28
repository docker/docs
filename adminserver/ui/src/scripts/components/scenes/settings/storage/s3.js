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
import { createValidator } from 'validation';

const mapState = createStructuredSelector({
    storage: getStorageRegistry
});

const formFields = [
    'rootdirectory',
    'region',
    'regionendpoint',
    'bucket',
    'accesskey',
    'secretkey',
    'v4auth',
    'secure'
];

@connect(mapState)
@ui({
    state: {
        showAdvanced: false
    }
})
@reduxForm({
    form: 's3',
    fields: formFields,
    validate: createValidator([
    ])
}, (state, props) => ({
    initialValues: {...props.storage.s3}
}))
export default class S3Settings extends Component {

    static propTypes = {
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
        this.props.saveFormData(data, 's3');
    }

    render() {

        const {
            fields: {
                rootdirectory,
                regionendpoint,
                region,
                bucket,
                accesskey,
                secretkey,
                v4auth,
                secure
            },
            handleSubmit,
			storage
        } = this.props;

        return (
            <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) }>
                <div className={ styles.formbox }>
                    <h2>S3 settings</h2>
                    <InputLabel
                        tip='Defaults to the empty string (bucket root).'
                        hint='Where registry files will be stored'>Root directory</InputLabel>
                    <Input type='text' formfield={ rootdirectory } />
                    <div className={ styles.row }>
                    <div className={ styles.halfColumn }>
                        <InputLabel
                            hint='Where you want to store objects'>AWS region name</InputLabel>
                        <Input type='text' placeholder='us-east-1' formfield={ region } />
                    </div>
                    <div className={ styles.halfColumn }>
                        <InputLabel
                            tip='Needs to be created before driver initialization.'
                            hint='Where you want to store objects'>S3 bucket name</InputLabel>
                        <Input type='text' formfield={ bucket } />
                    </div>
                    </div>
                    <InputLabel
                        tip='You can provide empty strings for your key if you plan on running the driver on an EC2 instance using IAM role grant credentials.'>AWS access key</InputLabel>
                    <Input type='text' formfield={ accesskey } />
                    <InputLabel
                        tip='You can provide empty strings for your key if you plan on running the driver on an EC2 instance using IAM role grant credentials.'>AWS secret key</InputLabel>
                    <Input type='text' formfield={ secretkey } />

                    <p><a href='#' onClick={ ::this.toggleAdvanced }>{ this.props.ui.showAdvanced ? 'Hide' : 'Show' } advanced settings</a></p>

                  { this.props.ui.showAdvanced &&
                    <div>
                      <InputLabel
                          isOptional={ true }
                          tip='Leave blank unless using a service other than AWS.'
                          hint='Where to find the S3 region'>Region Endpoint</InputLabel>
                      <Input type='text' formfield={ regionendpoint } />
                      <div>
                          <ToggleWithLabel labelOptions={
                              {
                                  labelText: 'Signature Version 4 Auth',
                                  inline: true,
                                  optional: true,
                                  tip: 'Indicates whether the registry uses Version 4 of AWS\'s authentication. Generally, you should set this to true.'
                              }
                          }
                          initial={ storage.s3 ? storage.s3.v4auth : false }
                          formField={ v4auth } />
                      </div>
                      <div>
                          <ToggleWithLabel labelOptions={
                              {
                                  labelText: 'Use HTTPS',
                                  inline: true,
                                  optional: true,
                                  tip: 'Indicates whether to use HTTPS instead of HTTP.'
                              }
                          }
                          initial={ storage.s3 ? storage.s3.secure : true }
                          formField={ secure } />
                    </div>
                  </div> }
                </div>
                <Button variant='primary'>Save</Button>
            </form>
        );
    }
}
