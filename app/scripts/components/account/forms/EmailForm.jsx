'use strict';

import React, { PropTypes } from 'react';
import has from 'lodash/object/has';
import classnames from 'classnames';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Button from '@dux/element-button';
import Card, { Block } from '@dux/element-card';

import { SplitSection } from 'common/Sections.jsx';
import { FlexTable, FlexRow, FlexHeader, FlexItem } from 'common/FlexTable.jsx';
import SimpleInput from 'common/SimpleInput';
import EmailElement from './email/EmailElement.jsx';
import addUserEmail from 'actions/addUserEmail.js';
import addUserEmailChange from 'actions/addUserEmailChange.js';
import setNewPrimaryEmail from 'actions/setNewPrimaryEmail';
import EmailsStore from 'stores/EmailsStore';

import styles from './EmailForm.css';

var debug = require('debug')('EmailForm');

var EmailForm = React.createClass({

  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  propTypes: {
    user: PropTypes.object.isRequired,
    JWT: PropTypes.string.isRequired,
    emails: PropTypes.array.isRequired,
    addEmail: PropTypes.string,
    addError: PropTypes.string,
    emailConfirmations: PropTypes.object
  },
  onAddEmailChange: function(e) {
    e.preventDefault();
    this.context.executeAction(addUserEmailChange, {email: e.target.value});
  },
  saveNewEmail: function(e) {
    e.preventDefault();
    this.context.executeAction(addUserEmail, {
      JWT: this.props.JWT,
      newEmail: this.props.addEmail,
      user: this.props.user
    });
  },
  setNewPrimary(id) {
    return (e) => {
      e.preventDefault();
      var payload = {
        JWT: this.props.JWT,
        username: this.props.user.username,
        emailId: id
      };
      this.context.executeAction(setNewPrimaryEmail, payload);
    };
  },
  _mkEmailElement({id, user, primary, email, verified}) {
    const {emailConfirmations} = this.props;
    let emailStatus = '';
    if (has(emailConfirmations, id)) {
      emailStatus = emailConfirmations[id];
    }
    return (<EmailElement emailid={id}
                          user={user}
                          isPrimaryEmail={primary}
                          email={email}
                          isVerified={verified}
                          JWT={this.props.JWT}
                          STATUS={emailStatus}
                          setNewPrimary={this.setNewPrimary}
                          key={id}
    />);
  },
  render: function() {
    const buttonVariant = !this.props.addError ? 'primary' : 'alert';
    let errorMsg;
    if (this.props.addError) {
      errorMsg = (<div className={styles.addError}>{this.props.addError}</div>);
    }
    return (
      <SplitSection title="Email Addresses"
                    module={false}
                    subtitle={(<div><p>This email address will be used for all notifications and correspondence from Docker.</p>
                                       <p>If you wish to designate a different email address as primary, first add a new address to your account and then click "make primary".</p></div>)}>
        <FlexTable>
          <FlexHeader>
            <FlexItem grow={3}>
                <form onSubmit={this.saveNewEmail} className={styles.noBottomMargin}>
                  <SimpleInput
                    type="email"
                    placeholder="New Email"
                    hasError={!!this.props.addError}
                    value={this.props.addEmail}
                    onChange={this.onAddEmailChange}
                    />
                  {errorMsg}
                </form>
            </FlexItem>
            <FlexItem grow={1}>
                <div className={styles.noBottomMargin}>
                  <Button onClick={this.saveNewEmail}
                          float="right"
                          variant={buttonVariant}>Add Email</Button>
                </div>
            </FlexItem>
          </FlexHeader>
          {this.props.emails.map(this._mkEmailElement)}
        </FlexTable>
      </SplitSection>
    );
  }
});

export default connectToStores(EmailForm,
                               [
                                 EmailsStore
                               ], function({ getStore }, props) {
                                 return getStore(EmailsStore).getState();
                               });
