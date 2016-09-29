'use strict';

const debug = require('debug')('Application');
import ApplicationStore from '../stores/ApplicationStore';
import JWTStore from '../stores/JWTStore';
import UserStore from '../stores/UserStore.js';
import React, { PropTypes, Component, cloneElement } from 'react';
import Welcome from './Welcome.jsx';
import connectToStores from 'fluxible-addons-react/connectToStores';
import Spinner from './Spinner.jsx';
import MainNav from './Nav.jsx';
import styles from './Application.css';
import DocumentTitle from 'react-document-title';
import merge from 'lodash/object/merge';
import { STORE_OFFICIAL_REPONAME_ID_MAP } from './store-promotion';

class Application extends Component {

  renderGenericStorePromo() {
    return (
      <div className={styles.storePromo}>
        <a href="https://store.docker.com">
          Docker Store is the new place to find all Docker content.&nbsp;
          <span className={styles.storePromoUnderline}>
            Try out the beta now →
          </span>
        </a>
      </div>
    );
  }

  renderRepoStorePromo(name, id) {
    return (
      <div className={styles.storePromo}>
        <a href={`https://store.docker.com/images/${id}`}>
          <strong>{name}</strong> is now available in the
        </a>
        <a href="https://store.docker.com"> Docker Store,</a>
        <a href={`https://store.docker.com/images/${id}`}>
          &nbsp;the new place to find Docker content.&nbsp;
          <span className={styles.storePromoUnderline}>
            Check it out →
          </span>
        </a>
      </div>
    );
  }


  render() {
    //Destructure router and props
    const { JWT, user, location, history, params } = this.props;
    //We use !JWT to check if a user is in logged out state

    if (!JWT && location.pathname === '/') {
      debug('Out Home');
      return (
        <DocumentTitle title='Docker Hub'>
          <main className={styles.main}>
            <MainNav
              isLoggedIn={false}
              isHomePage={true}
              history={history}/>
            <Welcome location={location} history={history}/>
          </main>
        </DocumentTitle>
        );
    } else if (history.isActive('/login/') || history.isActive('/reset-password/')) {
      debug('Application');
      return (
        <DocumentTitle title='Docker Hub'>
          <main className={styles.main}>
            {this.props.children}
          </main>
        </DocumentTitle>
      );
    } else {
      debug('Root');

      // Render store promotion based on the current location and only if the
      // store-promo flag is set.
      let promo = null;
      const storePromoFlag = this.props.location.query['store-promo'] === '1';
      if (storePromoFlag) {
        if (history.isActive('/_/') ||
            location.pathname.indexOf('/r/library/') === 0) {
          const name = params.splat;
          const id = STORE_OFFICIAL_REPONAME_ID_MAP[name];
          if (name && id) {
            promo = this.renderRepoStorePromo(name, id);
          }
        } else if (history.isActive('/search/')) {
          promo = this.renderGenericStorePromo();
        } else if (history.isActive('/explore/')) {
          promo = this.renderGenericStorePromo();
        }
      }

      const ifJWT = JWT ? JWT : '';
      return (
        <DocumentTitle title='Docker Hub'>
          <main className={styles.main}>
            {promo}
            <MainNav isLoggedIn={!!JWT}
              JWT={ifJWT}
              user={user}
              location={location}
              history={history}
              params={params}/>
            {this.props.children && cloneElement(this.props.children, {
              JWT: ifJWT,
              user,
              location
            })}
          </main>
        </DocumentTitle>
      );
    }
  }
}

export default connectToStores(Application,
                               [
                                 ApplicationStore,
                                 JWTStore,
                                 UserStore
                               ],
                               ({ getStore }, props) => {
                                 return merge({},
                                                getStore(ApplicationStore).getState(),
                                                {
                                                  JWT: getStore(JWTStore).getJWT(),
                                                  user: getStore(UserStore).getState()
                                                });
                               });
