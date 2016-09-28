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
import _ from 'lodash';

class Application extends Component {

  render() {
    //Destructure router and props
    const { JWT, user, location, history, params } = this.props;
    //We use !JWT to check if a user is in logged out state
    debug(this.props);
    if (!JWT && location.pathname === '/') {
      debug('Out Home');
      return (
        <main className={styles.main}>
          <MainNav
            isLoggedIn={false}
            isHomePage={true}
            history={history}/>
          <Welcome location={location} history={history}/>
        </main>
        );
    } else if (history.isActive('/login/') || history.isActive('/reset-password/')) {
      debug('Application');
      return (
        <main className={styles.main}>
          {this.props.children}
        </main>
      );
    } else {
      debug('Root');
      const ifJWT = JWT ? JWT : '';
      return (
        <main className={styles.main}>
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
                                 return _.merge({},
                                                getStore(ApplicationStore).getState(),
                                                {
                                                  JWT: getStore(JWTStore).getJWT(),
                                                  user: getStore(UserStore).getState()
                                                });
                               });
