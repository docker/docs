'use strict';

import React, { PropTypes } from 'react';
import deleteEmail from 'actions/deleteEmail.js';
import { FlexItem } from 'common/FlexTable.jsx';
import styles from './EmailComponents.css';
import FA from 'common/FontAwesome.jsx';

export default React.createClass({
  displayName: 'deleteEmailElement',
  propTypes: {
    isPrimaryEmail: PropTypes.bool.isRequired,
    emailid: PropTypes.number.isRequired,
    user: PropTypes.string.isRequired,
    JWT: PropTypes.string.isRequired
  },
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  deleteEmail: function(e) {
    e.preventDefault();
    var payload = {
      JWT: this.props.JWT,
      username: this.props.user,
      delEmailID: this.props.emailid
    };
    this.context.executeAction(deleteEmail, payload);
  },
  render() {
    if(this.props.isPrimaryEmail) {
      return (
        <FlexItem grow={0}/>
      );
    } else {
      return (
        <FlexItem end grow={0}>
          <a href="#" onClick={this.deleteEmail} className={styles.remove}>
            <FA icon='fa-trash-o' size='lg'/>
          </a>
        </FlexItem>
      );
    }
  }
});
