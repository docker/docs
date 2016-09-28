'use strict';

import styles from './RepositoryDetailsWrapper.css';
import React, { PropTypes, createClass, cloneElement } from 'react';
import connectToStores from 'fluxible-addons-react/connectToStores';
import RepoSecondaryNav from './RepoSecondaryNav.jsx';
import AutoBuildSettingsStore from '../../stores/AutoBuildSettingsStore';
import omit from 'lodash/object/omit';
import startsWith from 'lodash/string/startsWith';
import { STORE_OFFICIAL_REPONAME_ID_MAP } from './store-promotion';
const debug = require('debug')('hub::RepositoryDetailsWrapper');
const { string, shape, object, bool, number, func } = PropTypes;

const RepositoryDetails = createClass({
  displayName: 'RepositoryDetailsWrapper',
  propTypes: {
    description: string.isRequired,
    fullDescription: string.isRequired,
    isPrivate: bool.isRequired,
    isAutomated: bool.isRequired,
    name: string.isRequired,
    namespace: string.isRequired,
    status: number.isRequired,
    user: object,
    JWT: string,
    canEdit: bool,
    autoBuildStore: shape({
      provider: string.isRequired,
      repo_web_url: string.isRequired
    })
  },
  contextTypes: {
    getStore: func.isRequired
  },
  render() {
    const {
      name,
      namespace,
      isAutomated,
      JWT,
      user,
      canEdit,
      children,
      location
    } = this.props;

    const storePromoFlag = this.props.location.query['store-promo'] === '1';
    const storeID = STORE_OFFICIAL_REPONAME_ID_MAP[name];
    const libraryNamespace = namespace === 'library';
    const showPromo = storePromoFlag && storeID && libraryNamespace;

    return (
      <div>
        {showPromo && (
          <div className={styles.storePromo}>
            <a href={`https://store.docker.com/images/${storeID}`}>
              <strong>{name}</strong> is now available in the
            </a>
            <a href="https://store.docker.com"> Docker Store,</a>
            <a href={`https://store.docker.com/images/${storeID}`}>
              &nbsp;the new place to find Docker content. Check it out →
            </a>
          </div>
        )}
        <RepoSecondaryNav user={namespace}
                          splat={name}
                          canEdit={canEdit}
                          isAutomated={isAutomated}
                          isOfficialRoute={startsWith(location.pathname, '/_/')}/>
        <div className={styles.repoDetailsContent}>
          <div className='row'>
            <div className='large-12 columns'>
              {children && cloneElement(children, omit(this.props, 'children'))}
            </div>
          </div>
        </div>
      </div>
    );
  }
});

export default connectToStores(RepositoryDetails,
  [ AutoBuildSettingsStore ],
  function({ getStore }, props) {
    return getStore(AutoBuildSettingsStore).getState();
  });
