'use strict';

import React, { PropTypes } from 'react';
import { FlexTable, FlexRow, FlexHeader, FlexItem } from 'common/FlexTable.jsx';
import styles from './EmailComponents.css';
import PrimaryEmail from './PrimaryEmail';
import VerifiedOrResend from './VerifiedOrResend';
import DeleteEmailElement from './DeleteEmailElement';

export default React.createClass({
  displayName: 'EmailElement',
  propTypes: {
    user: PropTypes.string.isRequired,
    isPrimaryEmail: PropTypes.bool.isRequired,
    emailid: PropTypes.number.isRequired,
    email: PropTypes.string.isRequired,
    isVerified: PropTypes.bool.isRequired,
    JWT: PropTypes.string.isRequired,
    setNewPrimary: PropTypes.func.isRequired,
    STATUS: PropTypes.string
  },
  render(){
    const {
      JWT,
      user,
      isVerified,
      email,
      emailid,
      isPrimaryEmail,
      setNewPrimary,
      STATUS
    } = this.props;
    return (
            <FlexRow>
              <FlexItem grow={6}>
               <div className={styles.emailAddress}>{email}</div>
              </FlexItem>
              <VerifiedOrResend isVerified={isVerified}
                                email={email}
                                JWT={JWT}
                                emailid={emailid}
                                STATUS={STATUS}/>
              <PrimaryEmail checked={isPrimaryEmail}
                            JWT={JWT}
                            emailID={emailid}
                            isVerified={isVerified}
                            isPrimaryEmail={isPrimaryEmail}
                            setNewPrimary={setNewPrimary}/>
              <DeleteEmailElement emailid={emailid}
                                  user={user}
                                  JWT={JWT}
                                  isPrimaryEmail={isPrimaryEmail}/>
            </FlexRow>
    );
  }
});
