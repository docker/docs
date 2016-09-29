/**
MjA2MTU is a base64 encoding of the pk for the user account in hub

dhiltgen [11:07 AM]
41y-d0f275462b1678b1aeca is a random generated token
 **/
'use strict';

import React, { PropTypes } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';
import FancyInput from 'common/FancyInput';
import { Link } from 'react-router';
import resetPasswordSubmit from '../../actions/resetPasswordSubmit.js';
import ChangePasswordStore from '../../stores/ChangePasswordStore.js';
import clearChangePasswordStore from '../../actions/clearChangePasswordStore';
import styles from './ResetPass.css';
var debug = require('debug')('Password Reset Confirmation: ');

var ResetPassword = React.createClass({
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  PropTypes: {
    changePassStore: PropTypes.shape({
      reset: PropTypes.bool.isRequired,
      resetErr: PropTypes.bool.isRequired
    })
  },
  getInitialState: function() {
    return {
      pass1: '',
      pass2: '',
      passErr: false
    };
  },
  onPassChange: function(e) {
    e.preventDefault();
    this.setState({pass1: e.target.value});
  },
  onConfirmChange: function(e) {
    e.preventDefault();
    var confirm = e.target.value;
    this.setState({pass2: confirm, passErr: false});
  },
  onErrorRetry: function(e) {
    e.preventDefault();
    this.setState({
      pass1: '',
      pass2: ''
    });
    this.context.executeAction(clearChangePasswordStore);
  },
  transitionLogin: function(e) {
    e.preventDefault();
    this.props.history.pushState('/login/');
  },
  resetPassSubmit: function(e) {
    e.preventDefault();
    debug('reset pass submit');
    if (this.state.pass1 === this.state.pass2) {
      var { uidb64, reset_token } = this.props.params;
      this.context.executeAction(resetPasswordSubmit,
        {uidb64: uidb64,
          reset_token: reset_token,
          password_1: this.state.pass1,
          password_2: this.state.pass2});
      this.setState({reset: true});
    } else {
      this.setState({ passErr: true });
    }
  },
  render: function() {
    var disabledState = (this.state.passErr || !this.state.pass1);

    if (!this.props.changePassStore.reset && !this.props.changePassStore.hasErr) {
      let error;
      if (this.state.passErr) {
        error = 'Make sure passwords are identical';
      }
      return (
        <div className={styles.resetPassPage}>
          <div className='row'>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <div className={styles.head}>Password Reset</div>
            <div className={styles.subHead}>Enter your new password.</div>
            </div>
          </div>
          <div className='row'>
            <div className='small-6 small-centered large-3 large-centered columns'>
              <form onSubmit={this.resetPassSubmit} className={styles.resetPassSubmit}>
                <FancyInput placeholder='New Password'
                          onChange={this.onPassChange}
                          type="password"
                          value={this.state.pass1}
                          variant='white'/>
                <FancyInput placeholder='Confirm'
                  onChange={this.onConfirmChange}
                  type="password"
                  hasError={this.state.passErr}
                  error={error}
                  value={this.state.pass2}
                  variant='white'/>
                <Button type="submit" disabled={disabledState}>Reset Password</Button>
              </form>
            </div>
          </div>
        </div>
      );
    } else if (this.props.changePassStore.hasErr) {
      return (
        <div className={styles.resetPassPage}>
          <div className='row'>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <div className={styles.head}>There was an error!</div>
            <div className={styles.subHead}>Your password has not been reset, please try again</div>
            </div>
          </div>
          <div className='row'>
            <div className={'small-4 small-centered large-3 large-centered columns ' + styles.back}>
              <Button onClick={this.onErrorRetry}>Back</Button>
            </div>
          </div>
        </div>
      );
    } else {
      return (
        <div className={styles.resetPassPage}>
          <div className='row'>
            <div className={'small-6 small-centered large-4 large-centered columns ' + styles.header}>
              <div className={styles.head}>Your password has been reset</div>
            <div className={styles.subHead}>You may now login with your new password</div>
            </div>
          </div>
          <div className='row'>
            <div className={'small-4 small-centered large-3 large-centered columns ' + styles.back}>
              <Button onClick={this.transitionLogin}>Back to Login</Button>
            </div>
          </div>
        </div>
      );
    }
  }
});

module.exports = connectToStores(ResetPassword,
  [ChangePasswordStore],
  function({ getStore }, props) {
    return {
      changePassStore: getStore(ChangePasswordStore).getState()
    };
  });
