'use strict';
const debug = require('debug')('hub:app');
const Fluxible = require('fluxible');

let app = new Fluxible({
  component: require('./components/Routes.jsx'),
  stores: [
    require('./stores/AccountInfoFormStore'),
    require('./stores/AccountSettingsLicensesStore'),
    require('./stores/AddOrganizationStore'),
    require('./stores/AddTrialLicenseStore'),
    require('./stores/AddWebhookFormStore'),
    require('./stores/ApplicationStore'),
    require('./stores/AutobuildConfigStore'),
    require('./stores/AutoBuildSettingsStore'),
    require('./stores/AutobuildSourceRepositoriesStore'),
    require('./stores/AutobuildStore'),
    require('./stores/AutobuildTagsStore'),
    require('./stores/AutobuildTriggerByTagStore'),
    require('./stores/BillingInfoFormStore'),
    require('./stores/BillingPlansStore'),
    require('./stores/BitbucketLinkStore'),
    require('./stores/ChangePasswordStore'),
    require('./stores/CloudCouponStore'),
    require('./stores/CloudBillingStore'),
    require('./stores/ConvertToOrgStore'),
    require('./stores/CreateRepositoryFormStore'),
    require('./stores/DashboardContribsStore'),
    require('./stores/DashboardMembersStore'),
    require('./stores/DashboardNamespacesStore'),
    require('./stores/DashboardReposStore'),
    require('./stores/DashboardStarsStore'),
    require('./stores/DashboardStore'),
    require('./stores/DashboardTeamsStore'),
    require('./stores/DeletePipelineStore'),
    require('./stores/DeleteRepoFormStore'),
    require('./stores/EmailNotifStore'),
    require('./stores/EmailsStore'),
    require('./stores/EnterprisePaidFormStore'),
    require('./stores/EnterprisePartnerTrackingStore'),
    require('./stores/EnterpriseTrialFormStore'),
    require('./stores/EnterpriseTrialSuccessStore'),
    require('./stores/GithubLinkStore'),
    require('./stores/JWTStore'),
    require('./stores/LoginStore'),
    require('./stores/NotifyStore'),
    require('./stores/OrgTeamStore'),
    require('./stores/OrganizationStore'),
    require('./stores/OutboundCommunicationStore'),
    require('./stores/PipelineHistoryStore'),
    require('./stores/PlansStore'),
    require('./stores/PrivateRepoUsageStore'),
    require('./stores/RepoDetailsBuildLogs'),
    require('./stores/RepoDetailsBuildsStore'),
    require('./stores/RepoDetailsDockerfileStore'),
    require('./stores/RepoDetailsLongDescriptionFormStore'),
    require('./stores/RepoDetailsShortDescriptionFormStore'),
    require('./stores/RepoDetailsVisibilityFormStore'),
    require('./stores/RepoSettingsCollaborators'),
    require('./stores/RepositoryCommentsStore'),
    require('./stores/RepositoryPageStore'),
    require('./stores/SearchStore'),
    require('./stores/SignupStore'),
    require('./stores/TriggerBuildStore'),
    require('./stores/UserProfileReposStore'),
    require('./stores/UserProfileStore'),
    require('./stores/UserProfileStarsStore'),
    require('./stores/UserStore'),
    require('./stores/WebhooksSettingsStore')
  ]
});

/**
 * Add a plugin which adds the global redux store as .reduxStore to all flux
 * actions.
 *
 * This allows us to call actionContext.reduxStore.dispatch() to dispatch
 * actions from fluxible actions
 */
app.plug({
  name: 'ReduxActionIntegration',
  plugContext(options, context) {
    // Options should be passed reduxStore from the options passed into
    // createContext

    return {
      // Each action's context should also have the redux store that's defined
      // within the app context.
      plugActionContext(actionContext) {
        // NOTE: Server-side rendering passes this in as options.
        //       Client-side rendering passes this in as context.
        //       Fluxible has no way to specify options within rehydration.
        //       We can only affect context.

        // This is defined in server.js within `server.use` for server-side
        // rendering and within `app.rehydrate` in client.js
        actionContext.reduxStore = options.reduxStore || context.reduxStore;
      }
    };
  }
});

export default app;
