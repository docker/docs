'use strict';

import {
  Route, Redirect, IndexRoute
} from 'react-router';
import React from 'react';
import styles from './Routes.css';

const Account = require('./account/Account');
const AccountSettings = require('./account/AccountSettings');
const AddOrganizationForm = require('./account/AddOrganizationForm');
import AddRepo from './AddRepo';
const Application = require('./Application');
const Autobuild = require('./repositories/Autobuild');
const AutobuildIndex = require('./repositories/AutobuildIndex.jsx');
const AutoBuildSetupForm = require('./repositories/AutoBuildSetupForm');
const BillingPlans = require('./account/BillingPlans');
import UpdateBillingInfo from './account/UpdateBillingInfo';
import ChangePassSuccess from './welcome/ChangePassSuccess';
import CreateBillingSubscription from './account/CreateBillingSubscription';
import ConvertToOrg from './account/ConvertToOrg';
import CSEngineDownloadPage from './CSEngineDownloadPage';
import DashboardContribs from './dashboards/Contribs';
import DashboardWrapper from './DashboardWrapper';
import DashboardRepos from './dashboards/Repos';
import DashboardStars from './dashboards/Stars';
import orgDashTeams from './orgdashboards/Teams';
import Dockerfile from './repo/repo_details/Dockerfile';
import Explore from './Explore';
import EnterpriseSubscriptions from './enterprise/Enterprise';
import ServerTrial from './enterprise/EnterpriseTrial';
import EnterpriseTrialSuccess from './enterprise/EnterpriseTrialSuccess';
import EULA from './EULA';
import EnterpriseTrialTerms from './enterprise/EnterpriseTrialTerms.jsx';
import ForgotPassword from './welcome/ForgotPass';
import GithubLinkScopes from './account/services/GithubLinkScopes';
import Help from './help/Help';
import Licenses from './account/Licenses';
const LinkedAccountSourcesForm = require('./repositories/LinkedAccountSourcesForm');
const LinkedServices = require('./account/LinkedServices');
import Login from './Login';
import Members from './orgdashboards/Members';
const NotificationsSettings = require('./account/NotificationsSettings');
const OrganizationProfile = require('./account/orgs/OrganizationProfile');
const OrganizationSettings = require('./account/OrganizationSettings');
const OrganizationSummary = require('./account/OrganizationSummary');
import OrgDashboardWrapper from './OrgDashboardWrapper';
import Register from './Register';
import RepositorySettingsBuilds from './repo/repoSettings/Builds';
import RepositorySettingsCollaborators from './repo/repoSettings/CollaboratorsWrapper.jsx';
import RepositorySettingsMain from './repo/repoSettings/SettingsMain';
import RepositorySettingsWebhooks from './repo/repoSettings/webhooks';
import RepositorySettingsWrapper from './repo/RepositorySettingsWrapper';
import RepositoryDetailsBuildDetails from './repo/repo_details/BuildDetails';
import RepositoryDetailsBuildLogs from './repo/repo_details/BuildLogs';
import RepositoryDetailsInfo from './repo/repo_details/Info';
import RepositoryDetailsTags from './repo/repo_details/Tags';
import RepositoryDetailsScannedTag from './repo/repo_details/ScannedTag';
import RepositoryDetailsWrapper from './repo/RepositoryDetailsWrapper';
import RepositoryPageWrapper from './RepositoryPageWrapper';
import ResetPassword from './welcome/ResetPass';
import RouteNotFound404Page from './common/RouteNotFound404Page';
const Search = require('./search/Search');
const UserStars = require('./userWrapper/UserStars');
const User = require('./Users'); //Unused
import UserProfileWrapper from './UserProfileWrapper';
import UserProfileRepos from './userprofile/Repos';

// NOTE: Provider right now doesn't work easily with fluxible's render pipeline.
// Even though we add the <Provider> as a root element to Routes.jsx:
//
// module.exports = (
//   <Provider store={ store }>
//     { routes }
//   </Provider>
// );
//
// Fluxible doesn't render the root Provider component; it renders the first
// component which connects to fluxibles stores.
//
// FIX: We've instead added <Provider> as the base class that fluxibleRouter
//      renders
//
// TODO: When we rip out Fluxible add <Provider> as a top level component here.

var routes = (
  <Route name='app' component={ Application }>
    {/* Login and Password */}
    <Route name='login' path='/login/' component={Login} />
    <Route name='forgotPass' path='/reset-password/' component={ForgotPassword}/>
    <Route name='resetPass' path='/account/password-reset-confirm/:uidb64/:reset_token/' component={ResetPassword} />
    <Route name='passChangeSuccess' path='/account/password-reset-confirm/success/' component={ChangePassSuccess} />

    {/* Currently logged in user Dashboard */}
    <Route name='dashboard' path='/' component={DashboardWrapper}>
      <Route name='dashStars' path='stars/' component={DashboardStars}/>
      <Route name='dashContribs' path='contributed/' component={DashboardContribs}/>
      <IndexRoute name='dashboardHome' component={DashboardRepos}/>
    </Route>

    {/* Public user profile */}
    <Route name='userWrapper' path='/u/:user/' component={UserProfileWrapper}>
      <Route name='userStars' path='starred/' component={UserStars} />
      <IndexRoute name='user' component={UserProfileRepos}/>
    </Route>

    {/* Organization Dashboard */}
    <Route name='orgDashboard' path='/u/:user/dashboard/' component={OrgDashboardWrapper}>
      <Route name='orgDashTeams' path='teams/' component={orgDashTeams}/>
      <Route name='orgDashBilling' path='billing/' component={BillingPlans}/>
      <Route name='createOrgSubscription' path='billing/create-subscription/' component={CreateBillingSubscription} />
      <Route name='updateOrgBillingInfo' path='billing/update-info/' component={UpdateBillingInfo} />
      <Route name='orgDashSettings' path='settings/' component={OrganizationProfile}/>
      <IndexRoute name='orgDashHome' component={DashboardRepos}/>
    </Route>

    {/* Organizations Summary and Add Route */}
    <Route name='organizations' path='organizations/' component={OrganizationSettings} >
      <Route name='addOrg' path='add/' component={AddOrganizationForm} />
      <IndexRoute name='orgSummary' component={OrganizationSummary}/>
    </Route>

    {/* Add a repository route */}
    <Route name='addRepo' path='add/repository/' component={AddRepo} />

    {/* Autobuild creation related routes */}
    <Route name='autobuildGithub'
           path='add/automated-build/github/form/:sourceRepoNamespace/:sourceRepoName/'
           component={AutoBuildSetupForm} />
    <Route name='autobuildBitbucket'
           path='add/automated-build/bitbucket/form/:sourceRepoNamespace/:sourceRepoName/'
           component={AutoBuildSetupForm} />
    <Route path='add/automated-build/:userNamespace/' component={Autobuild}>
      <Route name='autobuildGithubOrgs' path='github/orgs/' component={LinkedAccountSourcesForm} />
      <Route name='autobuildBitbucketOrgs' path='bitbucket/orgs/' component={LinkedAccountSourcesForm} />
      <IndexRoute name='addAutoBuild' component={AutobuildIndex}/>
    </Route>

    {/* Github linking related route | the scope selection screen */}
    <Route name='githubScopes' path='account/authorized-services/github-permissions/' component={GithubLinkScopes} />

    {/* Official repositories route | TODO: add library/:name */}
    <Route name='repoOfficialWrapper' path='/_/' component={RepositoryPageWrapper}>
      <Route name='repoOfficialDetails' component={RepositoryDetailsWrapper}>
        <Route name='repoOfficial' path='*/' component={RepositoryDetailsInfo} />
      </Route>
    </Route>

    {/* THIS ROUTE IS A DUPLICATE OF /u/:user/. WHY IS THIS HERE??? */}
    <Route name='userWrapperRepos' path='/r/:user/' component={UserProfileWrapper}>
      <IndexRoute name='userRepos' component={UserProfileRepos}/>
    </Route>

    <Route name="repo" path="/r/:user/*/" component={RepositoryPageWrapper}>
      <Route name="repoSettings" path="~/settings/" component={RepositorySettingsWrapper}>
        <Route name="collaborators" path="collaborators/" component={RepositorySettingsCollaborators} />
        <Route name="webhooks" path="webhooks/" component={RepositorySettingsWebhooks} />
        <Route name='autobuildSettings' path='automated-builds/' component={RepositorySettingsBuilds} />
        <IndexRoute name='repoSettingsMain' component={RepositorySettingsMain} />
      </Route>
      <Route component={RepositoryDetailsWrapper}>
        <Route name='buildsMain' path='builds/' component={RepositoryDetailsBuildDetails} />
        <Route name='buildLogs' path='builds/:build_code/' component={RepositoryDetailsBuildLogs} />
        <Route name='repoDetailsTags' path='tags/' component={RepositoryDetailsTags}/>
        <Route name='repoDetailsScannedTag' path='tags/:tagname/' component={RepositoryDetailsScannedTag} />
        <Route name='dockerfile' path='~/dockerfile/' component={Dockerfile} />
        <IndexRoute name='repoDetailsInfo' component={RepositoryDetailsInfo} />
      </Route>
    </Route>

    {/* User Account Settings */}
    <Route name='account' path='account/' component={Account}>
      <Route name='accountSettings' path='settings/' component={AccountSettings} />
      <Route name='authorizedServices' path='authorized-services/'>
        <Route name='githubRedirect' path='github/' component={LinkedServices} />
        <Route name='bitbucketRedirect' path='bitbucket/' component={LinkedServices} />
        <IndexRoute name='authServicesRoot' component={LinkedServices} />
      </Route>
      <Route name='billingPlans' path='billing-plans/' component={BillingPlans} />
      <Route name='updateBillingInfo' path='billing-plans/update/' component={UpdateBillingInfo} />
      <Route name='createSubscription' path='billing-plans/create-subscription/' component={CreateBillingSubscription} />
      <Route name='notifications' path='notifications/' component={NotificationsSettings} />
      <Route name='licenses' path='licenses/' component={Licenses} />
      <Route name='toOrg' path='convert-to-org/' component={ConvertToOrg} />
      <IndexRoute name='settings' component={AccountSettings} />
    </Route>

    {/* Billing/Enterprise/Subscription related routes */}
    <Route name='publicBillingPage' path='billing-plans/' component={BillingPlans} />

    {/* TODO: @camacho 2/9/16 - remove routes after 1 week to give time for loaded clients to be updated*/}
    <Route name='subscriptions' path='/subscriptions/' component={EnterpriseSubscriptions} />
    <Route name='enterprise' path='/enterprise/' component={EnterpriseSubscriptions} />

    <Route name='serverTrial' path='/enterprise/trial/' component={ServerTrial} />
    <Route name='serverTrialTerms' path='/enterprise/trial/terms/' component={EnterpriseTrialTerms} />
    <Route name='serverTrialSuccess' path='/enterprise/trial/success/' component={EnterpriseTrialSuccess} />
    <Route name='eusa' path='/enterprise/eusa/' component={EULA}/>
    <Route name='csEngineDownloadPage' path='cs-engine/' component={CSEngineDownloadPage} />

    {/* Some publicly available routes to explore, search and ask for help */}
    <Route name='search' path='search/' component={Search} />
    <Route name='explore' path='/explore/' component={Explore}/>
    <Route name='help' path='/help/' component={Help}/>
    <Route name='register' path='/register/' component={Register}/>

    {/* Handle 404 for bad routes | If no route matches, render a 404 page */}
    <Route path='*' component={RouteNotFound404Page} />
  </Route>
);

module.exports = routes;
