'use strict';

import React, { cloneElement } from 'react';
import { Link } from 'react-router';
import FluxibleMixin from 'fluxible-addons-react/FluxibleMixin';
import OrganizationStore from '../../stores/OrganizationStore';
import { PageHeader } from 'dux';
import FA from '../common/FontAwesome';
const debug = require('debug')('OrganizationSettings');

var OrganizationSettings = React.createClass({
  displayName: 'OrganizationSettings',
  mixins: [FluxibleMixin],
  getInitialState: function() {
    return {
      orgs: this.context.getStore(OrganizationStore).getOrgs()
    };
  },
  statics: {
    storeListeners: {
      onOrgStoreChange: [OrganizationStore]
    }
  },
  onOrgStoreChange: function() {
    this.setState({
      orgs: this.context.getStore(OrganizationStore).getOrgs()
    });
  },
  componentDidMount: function() {
    this.setState({
      orgs: this.context.getStore(OrganizationStore).getOrgs()
    });
  },
  render: function() {

    var maybeCreateOrgBtn;
    var pathname = this.props.location.pathname;
    if (pathname.indexOf('/add/') === -1) {
      maybeCreateOrgBtn = <Link to="/organizations/add/" className='button'>Create Organization <FA icon='fa-plus' /></Link>;
    }

    return (
      <div className="orgs-settings">
        <PageHeader title='Organizations & Teams'>
          {maybeCreateOrgBtn}
        </PageHeader>
        <div className="orgs-body">
          {cloneElement(this.props.children, {
            user: this.props.user,
            JWT: this.props.JWT,
            orgs: this.state.orgs
          })}
        </div>
      </div>
    );
  }
});

module.exports = OrganizationSettings;
