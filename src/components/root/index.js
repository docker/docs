import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import styles from './styles.css';
import routes from 'lib/constants/routes';
import classnames from 'classnames';
import { accountFetchCurrentUser } from 'actions/account';
import TopNav from './TopNav';
import Footer from './Footer';
const { func, node, object } = PropTypes;

const mapStateToProps = ({ account }) => {
  const { currentUser, userEmails } = account;
  return {
    currentUser,
    userEmails,
  };
};

const dispatcher = {
  fetchCurrentUser: accountFetchCurrentUser,
};

@connect(mapStateToProps, dispatcher)
export default class Root extends Component {
  static propTypes = {
    children: node.isRequired,
    currentUser: object,
    userEmails: object,
    fetchCurrentUser: func.isRequired,
    location: object.isRequired,
    params: object.isRequired,
  }

  // Note: "display: 'none' is to hide the app until the css has been loaded
  // to avoid rendering unstyled markup. This is only in development
  render() {
    const { location, params } = this.props;
    const routesWithTransparentNavBar = [routes.home(), routes.beta()];
    const showTransparentNavBar =
      routesWithTransparentNavBar.indexOf(location.pathname) >= 0;
    const { id = 'puttingSomethingHereSoItDoesntBreak' } = params;
    const unWrappedRoutes = [
      routes.beta(),
      routes.bundleDetail({ id }),
      routes.home(),
      routes.imageDetail({ id }),
      routes.login(),
      routes.publisherAddProduct(),
    ];
    const isWrapped = unWrappedRoutes.indexOf(location.pathname) < 0;
    // Routes that do not have a footer or a topNav
    const noFooterOrTopNavRoutes = [
      routes.login(),
    ];
    let topNav;
    let footer;
    if (noFooterOrTopNavRoutes.indexOf(location.pathname) === -1) {
      topNav = (
        <TopNav
          showTransparentNavBar={showTransparentNavBar}
          location={location}
        />
      );
      footer = <Footer />;
    }

    const mainClasses = classnames({
      wrapped: isWrapped,
      [styles.main]: true,
    });

    return (
      <div style={{ display: 'none' }} className={styles.app}>
        {topNav}
        <div className={mainClasses}>
          {this.props.children}
        </div>
        {footer}
      </div>
    );
  }
}
