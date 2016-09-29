'use strict';

//The source repository could be github, bitbucket or anything else in the future
export default function selectSourceRepoForAutobuild(actionContext, selectedSourceRepo) {
  actionContext.dispatch('SELECT_SOURCE_REPO', selectedSourceRepo);
}
