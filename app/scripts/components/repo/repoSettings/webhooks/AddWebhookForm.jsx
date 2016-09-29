'use strict';

import React, {
  PropTypes,
  Component
} from 'react';

import connectToStores from 'fluxible-addons-react/connectToStores';
import AddWebhookFormStore from 'stores/AddWebhookFormStore';
import JWTStore from 'stores/JWTStore';

import {
  DEFAULT, ATTEMPTING, ERROR
} from 'stores/addwebhookformstore/Constants';

import addPipeline from 'actions/addPipeline';
import resetForm from 'actions/resetWebhookForm';

const debug = require('debug')('AddWebhookForm.jsx');
import styles from './AddWebhookForm.css';

class AddWebhookForm extends Component {
  static propTypes = {
    namespace: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    STATUS: PropTypes.oneOf([DEFAULT, ATTEMPTING, ERROR]),
    cancel: PropTypes.func.isRequired,
    serverErrors: PropTypes.object
  }

  static defaultProps = {
    STATUS: DEFAULT,
    serverErrors: {}
  }

  static contextTypes = {
    getStore: PropTypes.func.isRequired,
    executeAction: PropTypes.func.isRequired
  }

  state = {
    form: 'invalid',
    pipelineName: 'empty',
    hookUrl: 'empty'
  }

  componentDidMount() {
    this.refs.pipelineName.focus();
  }

  componentWillUnmount() {
    this.context.executeAction(resetForm);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.STATUS === ERROR && this.state.form !== 'invalid') {
      this.setState({form: 'valid'});
    }
  }

  validateForm() {
    const valid = this.refs.form.checkValidity();
    this.setState({
      form: valid ? 'valid' : 'invalid'
    });

    return valid;
  }

  validateInput(name) {
    const input = this.refs[name];
    const value = input.value;
    const valid = input.checkValidity();

    let inputState;

    if(value === '') {
      inputState = 'empty';
    } else {
      inputState = valid ? 'valid' : 'invalid';
    }

    this.setState({[name]: inputState});
    return valid;
  }

  possiblyRenderServerError() {
    if (this.props.STATUS !== ERROR) {
      return undefined;
    }

    const messages = [];

    try {
      const response = this.props.serverErrors.response.body;
      Object.keys(response).forEach((key) => {
        if (key === 'webhooks') {
          messages.push(response[key][0].hookUrl[0]);
        } else {
          messages.push(response[key]);
        }
      });
    } catch (e) {
      messages.push('Something went wrong! Please try again.');
    }

    return (
      <div className={styles.formError}>
        {messages.join(<br />)}
      </div>
    );
  }

  onInputUpdated(name) {
    this.validateInput(name);
    this.validateForm();
  }

  onClickCancel(event) {
    event.preventDefault();
    this.props.cancel();
  }

  onFormKeyUp(event) {
    // Look for press of escape key
    if (event.which === 27) {
      this.onClickCancel(event);
    }
  }

  onSubmit(event) {
    event.preventDefault();

    if (this.state.form !== 'valid') {
      return;
    }

    this.setState({
      form: 'submitting',
      serverErrors: null
    });

    const { jwt, namespace, name } = this.props;

    const webhooks = [{hookUrl: this.refs.hookUrl.value}];

    this.context.executeAction(addPipeline, {
      jwt,
      namespace,
      name,
      pipelineName: this.refs.pipelineName.value,
      webhooks
    });
  }

  render() {
    const formClasses = [styles.form];
    const pipelineNameClasses = [styles.input];
    const hookUrlClasses = [styles.input];

    formClasses.push(styles[`form_${this.state.form}`]);
    pipelineNameClasses.push(styles[`input_${this.state.pipelineName}`]);
    hookUrlClasses.push(styles[`input_${this.state.hookUrl}`]);

    return (
      <form
        className={formClasses.join(' ')}
        noValidate
        onSubmit={this.onSubmit.bind(this)}
        onKeyUp={this.onFormKeyUp.bind(this)}
        ref='form'
      >
        <div className={styles.formRow}>
          <div className={styles.inputGroup}>
            <input
              className={pipelineNameClasses.join(' ')}
              type='text'
              required
              minLength='3'
              name='pipelineName'
              id='pipelineName'
              ref='pipelineName'
              onChange={this.onInputUpdated.bind(this, 'pipelineName')}
            />
            <label
              className={styles.label}
              htmlFor="pipelineName"
            >
              Webhook name
            </label>
            <div className={styles.inputError}>
              Must be at least 3 chars
            </div>
          </div>
          <div className={styles.inputGroup}>
            <input
              className={hookUrlClasses.join(' ')}
              type='url'
              required
              placeholder='https://'
              name='hookUrl'
              id='hookUrl'
              ref='hookUrl'
              onChange={this.onInputUpdated.bind(this, 'hookUrl')}
            />
            <label
              className={styles.label}
              htmlFor="hookUrl"
            >
              Webhook URL
            </label>
            <div className={styles.inputError}>
              Valid URL required
            </div>
          </div>
        </div>
        <div className={[styles.formRow, styles.formRowActions].join(' ')}>
          {this.possiblyRenderServerError()}
          <button
            className={[styles.button, styles.buttonSecondary].join(' ')}
            type='button'
            onClick={this.onClickCancel.bind(this)}
          >
            Cancel
          </button>
          <button
            className={[styles.button, styles.buttonPrimary].join(' ')}
            type='submit'
            disabled={!(this.state.form === 'valid')}
          >
            {this.state.form === 'submitting' ? 'Submitting' : 'Save'}
          </button>
        </div>
      </form>
    );
  }
}
export default connectToStores(
  AddWebhookForm,
  [
   AddWebhookFormStore,
   JWTStore
  ],
  ({ getStore }, props) => {
   const webhookData = getStore(AddWebhookFormStore).getState();
   const { jwt } = getStore(JWTStore).getState();
   return {
     ...webhookData,
     jwt
   };
  }
);
