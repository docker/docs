'use strict';

import React, { Component, PropTypes, cloneElement } from 'react';
import { Link } from 'react-router';
import connectToStores from 'fluxible-addons-react/connectToStores';
import trunc from 'lodash/string/trunc';
import PrivateReposBlock from './dashboards/PrivateRepoStatusBlock.jsx';
import DashboardNamespacesStore from '../stores/DashboardNamespacesStore';
import { SecondaryNav } from 'dux';
import GravatarOption from './gravatar/GravatarOption';
import GravatarValue from './gravatar/GravatarValue';
import FA from './common/FontAwesome';
import LiLink from './common/LiLink';
import Select from 'react-select';
import styles from './DashboardWrapper.css';

const { array, func, shape, string } = PropTypes;
const debug = require('debug')('DashboardWrapper');

class DashboardWrapper extends Component {

  static propTypes = {
    //currentUserContext (dashboard user context, could be logged in user or one of the user's orgs)
    //namespaces (all the namespaces that the user can look at)
    dashboardNamespaces: shape({
      currentUserContext: string.isRequired,
      namespaces: array.isRequired
    }),
    JWT: string,
    user: shape({
      username: string
    })
  }

  mkOptions = (arr) => {
    return arr.map((namespace) => {
      return {value: namespace, label: trunc(namespace, 16)};
    });
  };

  changeContext = ({value: namespace, label}) => {
    if (namespace) {
      const { dashboardNamespaces, user, history, location } = this.props;
      let currentPage = parseInt(location.query.page, 10);
      if (dashboardNamespaces.currentUserContext === namespace) {
        //do nothing
        debug('doing nothing. Namespace is the same');
      } else if (user.username === namespace) {
        //navigate to home
        debug('navigate to home; its a user namespace');
        history.pushState(null, '/', {page: currentPage || 1});
      } else {
        //navigate to org dashboard
        debug('navigate to org dashboard');
        history.pushState(null, `/u/${namespace}/dashboard/`);
      }
    }
  };

  render() {
    const { children, dashboardNamespaces, JWT, user } = this.props;
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
                      value={dashboardNamespaces.currentUserContext}
                      options={this.mkOptions(dashboardNamespaces.namespaces)} />
            </li>
            <LiLink to='/' onlyActiveOnIndex>
              <FA icon='fa-book' /> Repositories
            </LiLink>
            <LiLink to='/stars/' activeClassName="active"><FA icon='fa-star' /> Stars</LiLink>
            <LiLink to='/contributed/' activeClassName="active"><FA icon='fa-pencil-square-o'/> Contributed</LiLink>
          </ul>
          <PrivateReposBlock />
        </SecondaryNav>
        {children && cloneElement(children, { JWT, user })}
      </div>
    );
  }

}

export default connectToStores(DashboardWrapper,
                               [
                                 DashboardNamespacesStore
                               ],
                               function({ getStore }, props) {
                                 return {
                                   dashboardNamespaces: getStore(DashboardNamespacesStore).getState()
                                 };
                               });
