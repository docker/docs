'use strict';

var debug = require('debug')('SignupForm');

import React, { PropTypes } from 'react';
import _ from 'lodash';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';

import FA from 'common/FontAwesome';
import FancyInput from 'common/FancyInput';
import SignupStore from '../../stores/SignupStore';
import attemptSignup from '../../actions/attemptSignup';
import onChange from '../../actions/common/onChangeUtil';
import { STATUS } from '../../stores/signupstore/Constants';
import { handleFormErrors } from './_utils';
import styles from './SignupForm.css';

var SignupForm = React.createClass({
  displayName: 'SignupForm',
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  PropTypes: {
    location: PropTypes.object
  },
  _onSubmit(e) {
    e.preventDefault();
    const { partner_value, redirect_value } = this.props.location.query;
    var payload = this.props.values;
    payload.username = payload.username.toLowerCase();
    if (partner_value) {
      payload.partner_value = partner_value;
    }
    if (redirect_value) {
      payload.redirect_value = redirect_value;
    }
    debug(payload);
    this.context.executeAction(attemptSignup, payload);
  },
  onChange: onChange({
    storePrefix: 'SIGNUP'
  }),
  render() {
    debug(this.props);
    if(this.props.STATUS === STATUS.SUCCESSFUL_SIGNUP) {
      return (
        <div className={styles.success}>
          <div className={styles.heading}>Sweet! You're almost ready to go!</div>
          <p className={styles.subtext}>Please check your email to activate your account.</p>
          <FA icon='fa-envelope' size='2x'/>
        </div>
      );
    } else {
      return (
        <div>
          <div className='row'>
            <div className='large-12 columns'>
              <h3 className={styles.heading}>New to Docker?</h3>
              <p className={styles.subtext}>Create your free Docker ID to get started.</p>
            </div>
          </div>
          <form onSubmit={this._onSubmit}
                className='row'>
            <div className='large-12 columns'>
              <FancyInput placeholder='Choose a Docker Hub ID'
                        hasError={ this.props.fields.username.hasError }
                        error={ this.props.fields.username.error }
                        onChange={this.onChange('username')}
                        value={this.props.values.username}
                        variant='white'/>
            </div>
            <div className='large-12 columns'>
              <FancyInput placeholder='Enter your email address'
                        type='email'
                        hasError={ this.props.fields.email.hasError }
                        error={ this.props.fields.email.error }
                        onChange={this.onChange('email')}
                        value={this.props.values.email}
                        variant='white'/>
            </div>
            <div className='large-12 columns'>
              <FancyInput placeholder='Choose a password'
                        type='password'
                        hasError={ this.props.fields.password.hasError }
                        error={ this.props.fields.password.error }
                        onChange={this.onChange('password')}
                        value={this.props.values.password}
                        variant='white'/>
            </div>
            <div className={'large-12 columns ' + styles.submit}>
              <Button type='submit'>Sign Up</Button>
            </div>
          </form>
        </div>
      );
    }
  }
});

export default connectToStores(SignupForm,
                               [SignupStore],
                               function({ getStore }, props) {
                                 return getStore(SignupStore).getState();
                               });
