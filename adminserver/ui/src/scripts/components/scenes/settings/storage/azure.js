'use strict';

import React, { Component, PropTypes } from 'react';
import styles from 'components/scenes/settings/formstyle.css';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import Button from 'components/common/button';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { getStorageRegistry } from 'selectors/settings';
import { checkRequiredFields, reportMissingFields } from 'validation';

const mapState = createStructuredSelector({
    storage: getStorageRegistry
});

const formFields = [
    'accountname',
    'accountkey',
    'container',
    'realm'
];

@connect(mapState)
@reduxForm({
    form: 'azure',
    fields: formFields,
    validate: (values) => {
        return reportMissingFields(checkRequiredFields(formFields, values), {});
    }
}, (state, props) => ({
    initialValues: {...props.storage.azure}
}))
export default class Azure extends Component {

    static propTypes = {
        fields: PropTypes.object,
        handleSubmit: PropTypes.func,
        saveFormData: PropTypes.func,
        storage: PropTypes.object
    }

    onSubmit = (data) => {
        this.props.saveFormData(data, 'azure');
    }

    render() {

        const {
            fields: {
                accountname,
                accountkey,
                container,
                realm
            },
            handleSubmit
        } = this.props;

        return (
            <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) }>
                <div className={ styles.formbox }>
                    <h2>Azure settings</h2>
                    <InputLabel
                        hint='Your Azure storage account'>Account name</InputLabel>
                    <Input type='text' formfield={ accountname } />
                    <InputLabel
                        hint='Primary or secondary key'>Account key</InputLabel>
                    <Input type='text' formfield={ accountkey } />
                    <InputLabel
                        tip='Must comply the storage container name requirements.'
                        hint='Name of the root storage container in where registry data will be stored'>Container</InputLabel>
                    <Input type='text' formfield={ container } />
                    <InputLabel
                        tip='Defaults to core.windows.net. For example, the realm for "Azure in China" would be "core.chinacloudapi.cn" and the realm for "Azure Government" would be "core.usgovcloudapi.net".'
                        hint='Domain name suffix for the storage service API endpoint'>Realm</InputLabel>
                    <Input type='text' formfield={ realm } />
                </div>
                <Button variant='primary'>Save</Button>
            </form>
        );
    }
}
