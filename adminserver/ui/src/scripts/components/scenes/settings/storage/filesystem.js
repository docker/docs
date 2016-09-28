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

const mapState = createStructuredSelector({
    storage: getStorageRegistry
});

@connect(mapState)
@reduxForm({
    form: 'filesystem',
    fields: ['rootdirectory']
}, (state, props) => ({
    initialValues: {...props.storage.filesystem}
}))
export default class FileSystem extends Component {

    static propTypes = {
        fields: PropTypes.object,
        handleSubmit: PropTypes.func,
        saveFormData: PropTypes.func,
        storage: PropTypes.object
    }

    onSubmit = (data) => {
        this.props.saveFormData(data, 'filesystem');
    }

    render() {

        const {
            fields: { rootdirectory },
            handleSubmit
        } = this.props;

        return (
                <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) }>
                    <div className={ styles.formbox }>
                        <h2>Filesystem settings</h2>
                        <InputLabel
                            tip={ <span>This is where you need to mount your distributed filesystem if you are using one.</span> }
                            hint='Where registry files will be stored'>Storage backend</InputLabel>
                        <Input type='text' placeholder='Set automatically by DTR' disabled formfield={ rootdirectory } />
                    </div>
                    <Button variant='primary' type='submit'>Save</Button>
                </form>
        );
    }
}
