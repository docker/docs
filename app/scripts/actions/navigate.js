'use strict';

var debug = require('debug')('hub:actions:navigate');
import _ from 'lodash';
import {
  JWT,
  Notifications,
  Repositories as Repos,
  Search,
  Users,
  Billing
} from 'hub-js-sdk';
import accountSettings from './navigate/accountSettings';
import addRepo from './navigate/addRepo';
import billingPlans from './navigate/billingPlans';
import bitbucketRedirect from './navigate/bitbucketRedirect';
import bitbucketUsersAndRepos from './navigate/bitbucketUsersAndRepos';
import buildsMain from './navigate/buildsMain';
import buildLogs from './navigate/buildLogs';
import autoBuildSettings from './navigate/autoBuildSettings';
import collaborators from './navigate/repoSettingsCollaborators';
import dashStars from './navigate/dashStars';
import dashContribs from './navigate/dashContribs';
import dockerfile from './navigate/dockerfile';
import explore from './navigate/explore';
import serverTrial from './navigate/serverTrial';
import serverTrialSuccess from './navigate/serverTrialSuccess';
import serverBilling from './navigate/serverBilling';
import cloudBilling from './navigate/cloudBilling';
import getNamespaces from './navigate/getNamespaces';
import githubUsersAndRepos from './navigate/githubUsersAndRepos';
import githubRedirect from './navigate/githubRedirect';
import home from './navigate/home';
import licenses from './navigate/licenses';
import linkedAccountSettings from './navigate/linkedAccountsSettings';
import notificationSettings from './navigate/notificationSettings';
import orgDashBilling from './navigate/orgBilling';
import orgDashTeams from './navigate/orgDashTeams';
import orgHome from './navigate/orgHome';
import orgSummary from './navigate/orgSummary';
import orgSettings from './navigate/orgSettings';
import repo from './navigate/repo';
import repoDetailsTags from './navigate/repoDetailsTags';
import repoDetailsScannedTag from './navigate/repoDetailsScannedTag';
import repoOfficial from './navigate/repoOfficial';
import repoSettings from './navigate/repoSettings';
import resetPass from './navigate/resetPass';
import search from './navigate/search';
import toOrg from './navigate/toOrg.js';
import UserStore from '../stores/UserStore';
import user from './navigate/user';
import userStars from './navigate/userStars';
import webhooks from './navigate/webhooks';

function noop({actionContext, payload, done, maybeData}) {
  done();
}

function routesHaveHandlerFor(route, routes){
  return _.has(routes, route);
}

function withUser(actionContext, payload, cb) {
  /**
   * If we are on the server and payload.cookies.jwt is a
   * jwt, use it
   */
  if (payload.cookies && payload.cookies.jwt) {
    let token = payload.cookies.jwt;
    actionContext.dispatch('RECEIVE_JWT', token);
    /**
     * Use the JWT to get the JWT's user's data.
     */
    JWT.getUser(token, function(err, res) {
      if (err) {
        debug('NOT REALLY EXPIRED. JUST AN ERROR');
        actionContext.dispatch('EXPIRED_SIGNATURE', null);
        return cb(null, err);
      } else {
        var cbData = res.body;
        cbData.isAdmin = res.body.is_admin;
        actionContext.dispatch('RECEIVE_USER', res.body);
        cb(null, {token: token, user: cbData});
      }
    });
  } else if (payload.jwt) {
    debug('payload has jwt');
    /**
     * if we have access to a jwt already, use it instead and assume
     * we already have the user data since we're likely on the client
     * (and we fill in the user data when a user logs in on the client)
     */
    cb(null, {
      token: payload.jwt,
      user: actionContext.getStore(UserStore).getState()
    });
  } else {
    debug('no jwt and payload has no cookies; This should not happen');
    /**
     * We have no jwt and no error? This shouldn't happen.
     */
    cb(null, {});
  }
}

module.exports = function(actionContext, payload, done) {
  var _done = done;
  done = function() {
    _done.apply(this, arguments);
  };

  if (!payload.location.pathname) {
    /**
     * if we don't have a pathname, react-router doesn't have a route.
     * ignore it. There's nothing we can do.
     */
    return done();
  }

  withUser(actionContext, payload, function(err, maybeData) {
    if (err) {
      debug(err);
      return done();
    }

    let routeName = payload.routes[payload.routes.length - 1].name;
    debug('routeName', routeName);
    let routes = {
      'accountSettings': accountSettings,
      'addRepo': addRepo,
      'addWebhook': noop,
      'addAutoBuild': linkedAccountSettings, // Questionable navigate Route
      'authServicesRoot': linkedAccountSettings,
      'autobuildBitbucket': getNamespaces,
      'autobuildBitbucketOrgs': bitbucketUsersAndRepos,
      'autobuildGithub': getNamespaces,
      'autobuildGithubOrgs': githubUsersAndRepos,
      'autobuildSettings': autoBuildSettings,
      'billingPlans': billingPlans,
      'bitbucketRedirect': bitbucketRedirect,
      'buildLogs': buildLogs,
      'buildsMain': buildsMain,
      'cloudBilling': cloudBilling,
      'collaborators': collaborators,
      'createOrgSubscription': orgDashBilling,
      'createSubscription': billingPlans,
      'dashboardHome': home,
      'dashContribs': dashContribs,
      'dashStars': dashStars,
      'dockerfile': dockerfile,
      'explore': explore,
      'githubRedirect': githubRedirect,
      'licenses': licenses,
      'notifications': notificationSettings,
      'orgDashBilling': orgDashBilling,
      'orgDashHome': orgHome,
      'orgDashSettings': orgSettings,
      'orgDashTeams': orgDashTeams,
      'orgSummary': orgSummary,
      'repoDetailsInfo': repo,
      'repoDetailsTags': repoDetailsTags,
      'repoDetailsScannedTag': repoDetailsScannedTag,
      'repoOfficial': repoOfficial,
      'repoSettingsMain': repoSettings,
      'resetPass': resetPass,
      'search': search,
      'serverBilling': serverBilling,
      'serverTrial': serverTrial,
      'serverTrialSuccess': serverTrialSuccess,
      'toOrg': toOrg,
      'updateBillingInfo': billingPlans,
      'updateOrgBillingInfo': orgDashBilling,
      'user': user,
      'userRepos': user, //This route is a clone of /u/:user/ WHY DO WE HAVE THIS?
      'userStars': userStars,
      'webhooks': webhooks
    };
    actionContext.dispatch('CHANGE_ROUTE', payload);
    if(routesHaveHandlerFor(routeName, routes)){
      routes[routeName]({ actionContext, payload, done, maybeData });
    } else {
      debug(`no handler for ${routeName}`, payload.routes);
      done();
    }
  });
};
