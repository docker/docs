'use strict';

export const validateRepositoryName = (name) =>
{
  var repoNamePattern = /^[a-z0-9]+(?:[._-][a-z0-9]+)*$/;
  return repoNamePattern.test(name);
};
