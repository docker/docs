'use strict';

// !IMPORTANT - @camacho 8/18/16 GCS
// DO NOT DELETE - THIS IS SHIPPING DARK TO BE TURNED ON IN >=2.2
// SEARCH CODEBASE FOR "!IMPORTANT - @camacho 8/18/16 GCS" FOR NECESSARY CHANGES
// TO MAKE GCS INTERFACE VISIBLE

import React, { Component, PropTypes } from 'react';
import styles from 'components/scenes/settings/formstyle.css';
import buttonStyles from 'components/common/button/button.css';
import Input from 'components/common/input';
import InputLabel from 'components/common/inputLabel';
import Button from 'components/common/button';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { getStorageRegistry } from 'selectors/settings';
import {
  createValidator,
  required,
  json
} from 'validation';

const mapState = createStructuredSelector({
  storage: getStorageRegistry
});

const formFields = [
  'rootdirectory',
  'bucket',
  'credentials'
];

@connect(mapState)
@reduxForm({
  form: 'gcs',
  fields: formFields,
  validate: createValidator({
    'bucket': [required],
    'credentials': [required, json]
  })
}, (state, props) => ({
  initialValues: {
    ...props.storage.gcs
  }
}))
export default class S3Settings extends Component {

  static propTypes = {
    fields: PropTypes.object,
    handleSubmit: PropTypes.func,
    saveFormData: PropTypes.func,
    storage: PropTypes.object,
    pristine: PropTypes.bool,
    invalid: PropTypes.bool
  }

  upload = () => {
    const { files } = this.refs.file;
    if (files.length < 1) {
      return;
    }
    let reader = new FileReader();
    reader.onload = this.onload;
    reader.readAsText(files[0]);
  }

  onload = (evt) => {
    // extract the json from the filereader event
    const content = evt.target.result;
    this.props.fields.credentials.onChange(content);
  }

  onSubmit = (data) => {
    this.props.saveFormData(data, 'gcs');
  }

  render() {
    const {
      fields: {
        rootdirectory,
        bucket,
        credentials
      },
      handleSubmit,
      invalid,
      pristine
    } = this.props;

    return (
      <form method='POST'
        onSubmit={ handleSubmit(::this.onSubmit) }
        style={ { marginBottom: '1rem' } }
      >
        <div className={ styles.formbox }>
          <h2>Google Cloud Storage settings</h2>
          <InputLabel
            tip='Defaults to the empty string (bucket root).'
            hint='Where registry files will be stored'
          >
            Root directory
          </InputLabel>
          <Input type='text' formfield={ rootdirectory } />
          <InputLabel
            tip='A Google Cloud Storage bucket. Create one before configuring this driver.'
            hint='Where you want to store objects'
          >
            Bucket name
          </InputLabel>
          <Input type='text' formfield={ bucket } />
          <InputLabel
            hint='Downloaded from Google Cloud Platform'
          >
            Credentials
          </InputLabel>
          <div className={ styles.row }>
            <div className={ styles.gcsUploadColumn }>
              <Input
                type='text'
                isTextarea
                placholder={
`{
  "type": "service_account",
  "project_id": "example-000000",
  "private_key_id": "0000000000000",
  "private_key": "-----BEGIN PRIVATE KEY-----\-----END PRIVATE KEY-----\n",
  "client_email": "example@example.com",
  "client_id": "000000000000000000",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/…."
}`
                }
                formfield={ credentials }
              />
            </div>
            <div className={ styles.gcsUploadLabelColumn }>
              OR
            </div>
            <div className={ styles.gcsUploadColumn }>
              <input
                type='file'
                className={ styles.gcsUpload }
                ref='file'
                accept='.json,.txt'
                name='gcs-upload'
                id='gcs-upload'
                onChange={ this.upload }
              />
              <label htmlFor='gcs-upload' className={ [
                buttonStyles.button,
                buttonStyles.primary,
                buttonStyles.outline
              ].join(' ') }>
                Upload JSON file
              </label>
            </div>
          </div>
        </div>
        <Button
          variant='primary'
          disabled={ pristine || invalid }
        >Save</Button>
      </form>
    );
  }
}
