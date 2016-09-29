'use strict';

const BUILD_TAG_DEFAULT = {
  SOURCE_NAME_BRANCH: '/^([^m]|.[^a]|..[^s]|...[^t]|....[^e]|.....[^r]|.{0,5}$|.{7,})/',
  SOURCE_NAME_TAG: '/.*/',
  DOCKER_TAG: '{sourceref}',
  SOURCE_NAME_PLACEHOLDER_BRANCH: 'All branches except master',
  SOURCE_NAME_PLACEHOLDER_TAG: '/.*/ This will target all tags',
  DOCKER_TAG_PLACEHOLDER_BRANCH: 'Same as branch',
  DOCKER_TAG_PLACEHOLDER_TAG: 'Same as tag'
};

export default BUILD_TAG_DEFAULT;
