 'use strict';
import React from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import PrivateRepoUsageStore from '../../stores/PrivateRepoUsageStore';
import styles from './PrivateRepoStatusBlock.css';
import classnames from 'classnames';

var debug = require('debug')('PrivateRepoStatusBlock');

//This is the block that shows how many private repos the user has available and how many have been used currently

var PrivateRepoStatusBlock = React.createClass({
  propTypes: {
    privateRepoUsed: React.PropTypes.number.isRequired,
    numFreePrivateRepos: React.PropTypes.number.isRequired,
    defaultRepoVisibility: React.PropTypes.string.isRequired,
    privateRepoAvailable: React.PropTypes.number.isRequired,
    privateRepoPercentUsed: React.PropTypes.number.isRequired,
    privateRepoLimit: React.PropTypes.number.isRequired,
    notAvailable: React.PropTypes.bool,
    isOrg: React.PropTypes.bool,
    orgNamespace: React.PropTypes.string.isRequired
  },
  getDefaultProps: function() {
    return {
      isOrg: false,
      orgNamespace: ''
    };
  },
  render: function() {
    const linkToBilling = this.props.isOrg ? `/u/${this.props.orgNamespace}/dashboard/billing/` : '/account/billing-plans/';
    if (!this.props.notAvailable) {
      return (
        <ul className="right">
          <li className={ styles.prLabel }>Private Repositories:</li>
          <li className={ styles.labelText }>Using {this.props.privateRepoUsed} of {this.props.privateRepoLimit}</li>
          <li><Link to={ linkToBilling }>Get more</Link></li>
        </ul>
      );
    } else {
      return <span />;
    }
  }
});

export default connectToStores(PrivateRepoStatusBlock,
                               [
                                 PrivateRepoUsageStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(PrivateRepoUsageStore).getState();
                               });
