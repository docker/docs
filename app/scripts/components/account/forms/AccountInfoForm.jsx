'use strict';

import React from 'react';
import classnames from 'classnames';
import connectToStores from 'fluxible-addons-react/connectToStores';

import AccountInfoFormStore from '../../../stores/AccountInfoFormStore';
import SimpleInput from 'common/SimpleInput.jsx';
import SaveSettings from '../../../actions/saveSettingsData';
import updateAccountInfoFormField from '../../../actions/updateAccountInfoFormField';
import Button from '@dux/element-button';
import { SplitSection } from '../../common/Sections.jsx';

import styles from './AccountInfoForm.css';
import FA from 'common/FontAwesome.jsx';
import {STATUS as BASE_STATUS} from 'stores/common/Constants.js';

var AccountInfoForm = React.createClass({
  contextTypes: {
    executeAction: React.PropTypes.func.isRequired
  },
  _onChange(fieldKey) {
    return (e) => {
      this.context.executeAction(updateAccountInfoFormField, {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  },
  onSubmit: function(e) {
    e.preventDefault();
    this.context.executeAction(SaveSettings, {
      JWT: this.props.JWT,
      username: this.props.user.username,
      updateData: this.props.values
    });
  },
  renderButton: function(STATUS) {
    switch (STATUS) {
      case BASE_STATUS.SUCCESSFUL:
        return (
          <Button type='submit'
                  variant='success'>Saved</Button>
              );
      case BASE_STATUS.ERROR:
        return (
          <Button type='submit'
                  variant='alert'>Save</Button>
              );
      case BASE_STATUS.ATTEMPTING:
        return (
          <Button type='submit'
                  variant='warning'>Saving <FA icon='fa-spinner fa-spin'/></Button>
              );
      default:
        return (
          <Button type='submit'
                  variant='primary'>Save</Button>
              );
    }
  },
  render: function() {
    const {
      full_name,
      company,
      location,
      profile_url,
      gravatar_email
    } = this.props.fields;
    return (
      <SplitSection title="Account Information"
                    subtitle={<p>This information will be visible to all users of Docker Hub.</p>}>
          <form className={styles.accountForm} onSubmit={this.onSubmit} >

              <div className={styles.lostFormInputs}>
                  <div className={styles.error}>
                    {full_name.error}
                  </div>
                  <SimpleInput placeholder="Full Name"
                               onChange={this._onChange('full_name')}
                               value={this.props.values.full_name}
                               hasError={full_name.hasError}/>

                  <div className={styles.error}>
                    {company.error}
                  </div>
                  <SimpleInput placeholder="Company"
                               onChange={this._onChange('company')}
                               value={this.props.values.company}
                               hasError={company.hasError}/>

                  <div className={styles.error}>
                    {location.error}
                  </div>
                  <SimpleInput placeholder="Location"
                               onChange={this._onChange('location')}
                               value={this.props.values.location}
                               hasError={location.hasError}/>

                  <div className={styles.error}>
                    {profile_url.error}
                  </div>
                  <SimpleInput placeholder="Profile Url"
                            onChange={this._onChange('profile_url')}
                            value={this.props.values.profile_url}
                            hasError={profile_url.hasError}/>

                  <div className={styles.error}>
                    {gravatar_email.error}
                  </div>
                  <SimpleInput placeholder="Gravatar Email"
                               onChange={this._onChange('gravatar_email')}
                               value={this.props.values.gravatar_email}
                               hasError={gravatar_email.hasError}/>
              </div>

              <div className={styles.accountButtonWrapper}>
                { this.renderButton(this.props.STATUS) }
              </div>

          </form>
      </SplitSection>
    );
  }
});

export default connectToStores(AccountInfoForm,
                               [
                                 AccountInfoFormStore
                               ],
                               function({ getStore }, props) {
                                 return getStore(AccountInfoFormStore).getState();
                               });
