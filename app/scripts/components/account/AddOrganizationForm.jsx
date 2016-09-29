'use strict';

import React from 'react';
import OrganizationStore from '../../stores/OrganizationStore';
import AddOrganizationStore from '../../stores/AddOrganizationStore';
import createOrganizationAction from '../../actions/createOrganization';
import updateAddOrganizationFormField from '../../actions/updateAddOrganizationFormField';
import connectToStores from 'fluxible-addons-react/connectToStores';
const debug = require('debug')('AddOrganizationForm');
import SimpleInput from 'common/SimpleInput.jsx';
import Button from '@dux/element-button';
import styles from './AddOrganizationForm.css';
import { SplitSection } from '../common/Sections.jsx';

var AddOrganizationForm = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  propTypes: {
    JWT: React.PropTypes.string.isRequired
  },
  _onChange(fieldKey) {
    return (e) => {
      this.context.executeAction(updateAddOrganizationFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  _handleCreate: function(e) {
    e.preventDefault();
    /*eslint-disable camelcase */
    var newOrg = {
      orgname: this.props.values.orgname.toLowerCase(),
      full_name: this.props.values.full_name,
      gravatar_email: this.props.values.gravatar_email,
      company: this.props.values.company,
      location: this.props.values.location,
      profile_url: this.props.values.profile_url
      /*eslint-enable camelcase */
    };
    this.context.executeAction(createOrganizationAction, {
      jwt: this.props.JWT,
      organization: newOrg
    });
  },
  render: function() {
    return (
      <div className={styles.contentWrapper}>
        <SplitSection title='Create Organization'
                      subtitle={<p>
                                Organizations can have multiple Teams. Teams can have differing permissions. Namespace is
                                unique and this is where repositories for this organization will be created.
                                                                      </p>}>
          <form onSubmit={this._handleCreate}>
            <br />
            <SimpleInput placeholder='Namespace'
                         onChange={this._onChange('orgname')}
                         value={this.props.values.orgname}
                         {...this.props.fields.orgname} />

            <SimpleInput placeholder="Organization Full Name"
                        onChange={this._onChange('full_name')}
                        value={this.props.values.full_name}
                        {...this.props.fields.full_name} />
            <SimpleInput placeholder="Company"
                        onChange={this._onChange('company')}
                        value={this.props.values.company}
                        {...this.props.fields.company} />
            <SimpleInput placeholder="Location"
                        onChange={this._onChange('location')}
                        value={this.props.values.location}
                        {...this.props.fields.location} />
            <SimpleInput placeholder="Gravatar Email"
                        onChange={this._onChange('gravatar_email')}
                        value={this.props.values.gravatar_email}
                        {...this.props.fields.gravatar_email} />
            <SimpleInput placeholder="Website URL"
                        onChange={this._onChange('profile_url')}
                        value={this.props.values.profile_url}
                        {...this.props.fields.profile_url} />
            <Button type="submit">Create</Button>
          </form>
        </SplitSection>
      </div>
    );
  }
});

export default connectToStores(AddOrganizationForm,
                               [
                                 AddOrganizationStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(AddOrganizationStore).getState();
                               });
