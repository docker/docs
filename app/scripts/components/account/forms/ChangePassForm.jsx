'use strict';

import React, { PropTypes } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import _ from 'lodash';
import classnames from 'classnames';

import ChangePasswordStore from '../../../stores/ChangePasswordStore';
import SimpleInput from 'common/SimpleInput.jsx';
import changePass from '../../../actions/changePassword';
import updateChangePassStore from '../../../actions/updateChangePassStore';
import Button from '@dux/element-button';
import { SplitSection } from '../../common/Sections.jsx';

import styles from './ChangePassForm.css';

var debug = require('debug')('ChangePassForm');

var ChangePassForm = React.createClass({
  contextTypes: {
    getStore: PropTypes.func.isRequired,
    executeAction: PropTypes.func.isRequired
  },
  getDefaultProps: function() {
    return {
      JWT: '',
      username: ''
    };
  },
  getInitialState: function() {
    return {
      confpassErr: ''
    };
  },
  onStoreChange: function() {
    var store = this.context.getStore(ChangePasswordStore);
    this.setState(store.getState());
  },
  oldPassChange: function(e) {
    e.preventDefault();
    var oldpass = e.target.value;
    var payload = {
      oldpass: oldpass,
      newpass: this.props.changePasswordStore.newpass,
      confpass: this.props.changePasswordStore.confpass
    };
    this.context.executeAction(updateChangePassStore, payload);
  },
  newPassChange: function(e) {
    e.preventDefault();
    var newpass = e.target.value;
    var payload = {
      oldpass: this.props.changePasswordStore.oldpass,
      newpass: newpass,
      confpass: this.props.changePasswordStore.confpass
    };
    this.context.executeAction(updateChangePassStore, payload);
  },
  confirmChange: function(e) {
    var confpass = e.target.value;
    var payload = {
      oldpass: this.props.changePasswordStore.oldpass,
      newpass: this.props.changePasswordStore.newpass,
      confpass: confpass
    };
    this.context.executeAction(updateChangePassStore, payload);
    this.setState({
      confpassErr: false
    });
  },
  onSubmit: function(e) {
    e.preventDefault();
    var store = this.props.changePasswordStore;
    if (store.confpass === store.newpass) {
      var payload = {
        JWT: this.props.JWT,
        username: this.props.username,
        oldpassword: this.props.changePasswordStore.oldpass,
        newpassword: this.props.changePasswordStore.newpass
      };
      this.context.executeAction(changePass, payload);
    } else {
      debug('passwords do not match');
      this.setState({
        confpassErr: true
      });
    }
  },
  render: function() {
    var store = this.props.changePasswordStore;
    var oldHasErr = _.has(store.err, 'old_password') || _.has(store.err, 'non_field_errors');
    var oldError;
    var newHasErr = _.has(store.err, 'new_password');
    var newError;
    if (oldHasErr) {
      if (_.has(store.err, 'old_password')) {
        oldError = store.err.old_password[0];
      } else {
        oldError = store.err.non_field_errors[0];
      }
    }
    if (newHasErr) {
      newError = store.err.new_password[0];
    }

    return (
      <SplitSection title="Change Password"
                    subtitle={<p>Please choose a password which is longer than 6 characters.</p>}>
          <form onSubmit={this.onSubmit} >
            <div className={styles.wrapper}>
              <div className={styles.lostFormInputs}>

                <label className={styles.label}>Old password</label>
                <div className={styles.error}>
                  { oldError }
                </div>
                <SimpleInput
                  type="password"
                  placeholder="Enter current password"
                  value={store.oldpass}
                  hasError={oldHasErr}
                  onChange={this.oldPassChange}/>

                <label className={styles.label}>New password</label>
                <div className={styles.error}>
                  { newError }
                </div>
                <SimpleInput
                  type="password"
                  placeholder="Enter new password"
                  value={store.newpass}
                  hasError={newHasErr}
                  onChange={this.newPassChange}/>

                <label className={styles.label}>Confirm password</label>
                <div className={styles.error}>
                  { this.state.confpassErr ? 'Make sure passwords are identical' : '' }
                </div>
                <SimpleInput
                  type="password"
                  placeholder="Confirm new password"
                  hasError={!!this.state.confpassErr}
                  value={store.confpass}
                  onChange={this.confirmChange}/>

              </div>
              <div className={styles.buttons}>
                <Button type='submit'>Save</Button>
              </div>
            </div>
          </form>
      </SplitSection>
    );
  }
});

export default connectToStores(ChangePassForm,
  [ChangePasswordStore],
  function({ getStore }, props) {
    return {
      changePasswordStore: getStore(ChangePasswordStore).getState()
    };
  });
