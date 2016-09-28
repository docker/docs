// This uses "export-from" in ES7:
// https://github.com/leebyron/ecmascript-more-export-from

'use strict';

export Login from './scenes/login';
export Settings from './scenes/settings';
export Repositories from './scenes/repositories';
export Repository, {
  RepositoryActivityTab,
  RepositoryDetailsTab,
  RepositoryTeamsTab,
  RepositoryUsersTab,
  RepositorySettingsTab,
  RepositoryTagsTab
} from './scenes/repository';
export Support from './scenes/support';
export API from './scenes/api';

export Organizations from './scenes/organizations';
export Organization from './scenes/organization';

export StorageSettings from './scenes/settings/storage';
export GeneralSettings from './scenes/settings/general';
export GarbageCollection from './scenes/settings/gc';

export Users from './scenes/users';
export User from './scenes/user';
export UserRepos from './scenes/user/repos';
export UserTeams from './scenes/user/teams';
export UserSettings from './scenes/user/settings';

export NotFound from './scenes/notFound';
