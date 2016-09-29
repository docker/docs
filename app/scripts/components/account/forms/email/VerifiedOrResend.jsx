'use strict';

import React, { PropTypes } from 'react';
import resendConfirmationEmail from 'actions/resendConfirmationEmail';
import { FlexItem } from 'common/FlexTable.jsx';
import styles from './EmailComponents.css';
import FA from 'common/FontAwesome.jsx';
import { EMAILSTATUS } from 'stores/emailsstore/Constants';

export default React.createClass({
  displayName: 'VerifiedOrResend',
  propTypes: {
    isVerified: PropTypes.bool.isRequired,
    email: PropTypes.string.isRequired,
    emailid: PropTypes.number.isRequired,
    STATUS: PropTypes.string
  },
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  _resendConfirmation(e) {
    e.preventDefault();
    this.context.executeAction(resendConfirmationEmail, {
      JWT: this.props.JWT,
      emailID: this.props.emailid
    });
  },
  render() {
    const { STATUS } = this.props;
    if(this.props.isVerified){
      return (<FlexItem grow={2}>
                <span className={styles.emphasis}>verified</span>
              </FlexItem>);
    } else if (STATUS === EMAILSTATUS.ATTEMPTING) {
      return (<FlexItem grow={2}>Sending <FA icon='fa-spinner fa-spin'/></FlexItem>);
    } else if (STATUS === EMAILSTATUS.SUCCESS) {
      return (<FlexItem grow={2}>Email Sent!</FlexItem>);
    } else if (STATUS === EMAILSTATUS.FAILED) {
      return (
        <FlexItem grow={2}>
          <div className={styles.failed}>Failed</div>
        </FlexItem>
      );
    } else {
      return (
        <FlexItem grow={2}>
          <a href="#"
             onClick={this._resendConfirmation}>Resend Email</a>
         </FlexItem>
      );
      }
  }
});
