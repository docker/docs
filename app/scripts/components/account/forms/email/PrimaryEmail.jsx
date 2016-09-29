'use strict';

import React, { PropTypes } from 'react';
import { FlexItem } from 'common/FlexTable.jsx';
import styles from './EmailComponents.css';

export default React.createClass({
  displayName: 'PrimaryEmail',
  propTypes: {
    isVerified: PropTypes.bool.isRequired,
    isPrimaryEmail: PropTypes.bool.isRequired,
    emailID: PropTypes.number.isRequired,
    setNewPrimary: PropTypes.func.isRequired
  },
  contextTypes: {
    executeAction: PropTypes.func.isRequired
  },
  render() {
    if (!this.props.isVerified) {
      return (
        <FlexItem grow={2} />
      );
    } else if (this.props.isPrimaryEmail) {
      return (
        <FlexItem grow={2}>
          <span className={styles.emphasis}>primary</span>
        </FlexItem>
      );
    } else {
      return (
        <FlexItem grow={2}>
          <a href="#" onClick={this.props.setNewPrimary(this.props.emailID)}>make primary</a>
        </FlexItem>
      );
    }
  }
});
