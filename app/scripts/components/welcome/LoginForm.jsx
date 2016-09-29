'use strict';

import React, { createClass, PropTypes } from 'react';
import _ from 'lodash';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';
var debug = require('debug')('LoginForm');

import onChange from '../../actions/common/onChangeUtil';
import FancyInput from 'common/FancyInput';
import LoginStore from '../../stores/LoginStore';
import loginUpdateFormField from '../../actions/loginUpdateFormField';
import styles from './LoginForm.css';

var LoginForm = createClass({
  displayName: 'LoginForm',
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  propTypes: {
    loginAction: PropTypes.func.isRequired
  },
  _onChange: onChange({
    storePrefix: 'LOGIN'
  }),
  _handleSubmit(event) {
    event.preventDefault();
    var loginPayload = this.props.values;
    loginPayload.username = loginPayload.username.toLowerCase();
    this.context.executeAction(this.props.loginAction, loginPayload);
  },
  render() {
    const { fields, values, variant, includeHelp } = this.props;

    var globalFormError = null;
    if(this.props.globalFormError) {
      globalFormError = <p className='alert-box alert'>{this.props.globalFormError}</p>;
    }

    let loginButton = (<Button type='submit'>Log In</Button>);
    if(this.props.STATUS === 'ATTEMPTING_LOGIN') {
      loginButton = (<Button type='submit' disabled>Logging In...</Button>);
    }
    let help;
    if (includeHelp) {
      help = (
          <Link to="/reset-password/">Can't Login?</Link>
      );
    }
    return (
      <form onSubmit={this._handleSubmit}
            className={styles.formWrapper}>
        {globalFormError}
        <FancyInput placeholder='Username'
                  onChange={this._onChange('username')}
                  hasError={fields.username.hasError}
                  error={fields.username.error}
                  value={values.username}
                  variant={variant}/>
        <FancyInput placeholder='Password'
                  type='password'
                  onChange={this._onChange('password')}
                  hasError={fields.password.hasError}
                  error={fields.password.error}
                  value={values.password}
                  variant={variant}/>
        <div className={styles.buttonWrapper}>
          <div className={styles.help}>
            {help}
          </div>
          {loginButton}
        </div>
      </form>
    );
  }
});

export default connectToStores(LoginForm,
                               [LoginStore],
                               function({ getStore }, props){
                                 return getStore(LoginStore).getState();
                               });
