'use strict';

import React, { Component, PropTypes } from 'react';
const { bool, object, func, array, instanceOf } = PropTypes;
import { connect } from 'react-redux';
import { reset as resetReduxForm } from 'redux-form';
import { Map } from 'immutable';
import consts from 'consts';
import { mapActions } from 'utils';
import * as OrganizationActions from 'actions/organizations';
import uiDecorator from 'redux-ui';
import autoaction from 'autoaction';
// Components
import { VelocityComponent } from 'velocity-react';
import Button from 'components/common/button';
import Spinner from 'components/common/spinner';
import NewOrganization from './newOrganization';
import OrganizationList from './organizationList.js';
// Selectors
import { orgNamesSelector } from 'selectors/organizations';
import { isAdminSelector } from 'selectors/users';
import { getReposByNamespace } from 'selectors/repositories';
import { createStructuredSelector } from 'reselect';
import { currentUserSelector } from 'selectors/users';
// Misc
import styles from './list.css';

const mapState = createStructuredSelector({
  orgs: orgNamesSelector,
  isAdmin: isAdminSelector,
  reposByOrg: getReposByNamespace,
  user: currentUserSelector
});

@uiDecorator({
  key: 'organizations',
  state: {
    isFormVisible: false
  },
  reducer: (state, action) => {
    if (action.type === consts.organizations.CREATE_ORGANIZATION && action.ready && !action.error) {
      return state.set('isFormVisible', false);
    }
    return state;
  }
})
@connect(mapState, mapActions({...OrganizationActions, resetReduxForm }))
@autoaction({
  listOrganizations: []
}, OrganizationActions)
export default class Organizations extends Component {

  static propTypes = {
    actions: object,
    orgs: array,
    isAdmin: bool,
    reposByOrg: instanceOf(Map),

    ui: object,
    updateUI: func,
    location: object,
    user: object
  }

  static childContextTypes = {
      location: object
  }

  getChildContext() {
      return {
          location: this.props.location
      };
  }

  createOrganization(data) {
    data.user = this.props.user;
    return this.props.actions.createOrganization(data).then(() => {
      // On successful submission, reset the form.
      // On error, an {_error: ... } object will be returned from action creator instead
      this.props.actions.resetReduxForm('createOrganizationForm');
    });
  }

  showNewForm(visible = true) {
    this.props.updateUI('isFormVisible', visible);
  }

  renderActions() {
    if (this.props.isAdmin) {
      // Only super users can create new orgs
      return (
        <div className={ styles.controls }>
          <span className={ styles.createFormButtonWrapper }>
            <Button
              id='new-organization-button'
              variant='secondary'
              onClick={ () => this.showNewForm() }
              disabled={ this.props.ui.isFormVisible }>
              New organization
            </Button>
          </span>
        </div>
      );
    }
  }

  renderCreateOrgForm() {
    return (
      <VelocityComponent animation={ this.props.ui.isFormVisible ? 'slideDown' : 'slideUp' } duration={ 250 }>
        <div className={ styles.newFormWrapper }>
          <NewOrganization
            dismissForm={ () => this.showNewForm(false) }
            onSubmit={ ::this.createOrganization }
            initialValues={ { name: '' } } />
        </div>
      </VelocityComponent>
    );
  }

  renderZeroState() {
    const showZeroState = !this.props.orgs || this.props.orgs.length || this.props.ui.isFormVisible;
    if (showZeroState) {
      return null;
    }
    if (this.props.isAdmin) {
      return (
        <div className={ styles.zero }>
          <img src='/public/img/octopus.png' alt='' />
          <p>
            You have <b>no organizations</b>...<br />
            make one now!
          </p>
          <div className={ styles.arrow }></div>
        </div>
      );
    } else {
      return (
        <div className={ styles.zero }>
          <img src='/public/img/octopus.png' alt='No organizations!' />
          <div>
            There are no organizations...
            Ask your System Administrator to make one!
          </div>
        </div>
      );
    }
  }

  render() {
    const listOrgsStatus = [consts.organizations.LIST_ORGANIZATIONS];
    return (
      <div className={ styles.wrapper }>
        <Spinner className={ styles.container } loadingStatus={ [listOrgsStatus] }>
          { this.renderActions() }
          { this.renderCreateOrgForm() }
          { this.renderZeroState() }
          <OrganizationList
            reposByOrg={ this.props.reposByOrg }
            organizations={ this.props.orgs } />
        </Spinner>
      </div>
    );
  }

}
