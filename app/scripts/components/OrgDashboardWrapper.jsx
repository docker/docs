'use strict';

import React, { Component, PropTypes, cloneElement } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import includes from 'lodash/collection/includes';
import merge from 'lodash/object/merge';
import trunc from 'lodash/string/trunc';
import LiLink from './common/LiLink';
import FA from './common/FontAwesome';
import GravatarOption from './gravatar/GravatarOption';
import GravatarValue from './gravatar/GravatarValue';
import PrivateReposBlock from './dashboards/PrivateRepoStatusBlock.jsx';
import DashboardNamespacesStore from '../stores/DashboardNamespacesStore';
import OrganizationStore from '../stores/OrganizationStore.js';
import Route404 from 'common/RouteNotFound404Page.jsx';
import { SecondaryNav } from 'dux';
import Select from 'react-select';
import styles from './DashboardWrapper.css';

const { array, func, object, shape, string } = PropTypes;
const debug = require('debug')('OrgDashboardWrapper');

class OrgDashboardWrapper extends Component {

  static propTypes = {
    currentUserContext: string,
    JWT: string,
    org: object,
    user: shape({
      username: string
    }),
    namespaces: array,
    ownedNamespaces: array
  }

  mkOptions = (arr) => {
    return arr.map((namespace) => {
      return {value: namespace, label: trunc(namespace, 16)};
    });
  };

  changeContext = ({value: namespace, label}) => {
    if (namespace) {
      const { currentUserContext, user } = this.props;
      if (currentUserContext === namespace) {
        //do nothing
        debug('doing nothing. Namespace is the same');
      } else if (user.username === namespace) {
        //navigate to home
        debug('navigate to home; its a user namespace');
        this.props.history.pushState(null, '/');
      } else {
        //navigate to org dashboard
        debug('navigate to org dashboard');
        this.props.history.pushState(null, `/u/${namespace}/dashboard/`);
      }
    }
  };

  render() {

    const {
      currentUserContext,
      JWT,
      namespaces,
      org,
      ownedNamespaces,
      children
    } = this.props;

    if (JWT) {
      let Settings = (
        <LiLink to={`/u/${currentUserContext}/dashboard/settings/`}>
            <FA icon='fa-gear'/> Settings
        </LiLink>
      );
      let Billing = (
        <LiLink to={`/u/${currentUserContext}/dashboard/billing/`}>
          <FA icon='fa-money'/> Billing
        </LiLink>
      );
      let isOwner = true;
      if (!includes(ownedNamespaces, currentUserContext)) {
        Settings = null;
        Billing = null;
        isOwner = false;
      }
      return (
        <div className="dashboard">
          <SecondaryNav>
            <ul className='left'>
              <li>
                <Select onChange={this.changeContext}
                        className={styles.select}
                        clearable={false}
                        optionComponent={GravatarOption}
                        valueComponent={GravatarValue}
                        value={currentUserContext}
                        options={this.mkOptions(namespaces)} />
              </li>
              <LiLink to={`/u/${currentUserContext}/dashboard/`} onlyActiveOnIndex>
                    <FA icon='fa-book'/> Repositories
              </LiLink>
              <LiLink to={`/u/${currentUserContext}/dashboard/teams/`}>
                <FA icon='fa-users'/> Teams
              </LiLink>
              { Billing }
              { Settings }
            </ul>
            <PrivateReposBlock isOrg={ true } orgNamespace={ currentUserContext }/>
          </SecondaryNav>
          {cloneElement(children, {
            JWT,
            user: org,
            isOwner,
            currentUserContext
          })}
        </div>
      );
    } else {
      return (
        <Route404 />
      );
    }
  }
}

export default connectToStores(OrgDashboardWrapper,
                               [
                                 DashboardNamespacesStore,
                                 OrganizationStore
                               ],
                               function({ getStore }, props) {
                                 return merge(
                                   {},
                                   getStore(DashboardNamespacesStore).getState(),
                                   {org: getStore(OrganizationStore).getCurrentOrg()}
                                 );
                               });
