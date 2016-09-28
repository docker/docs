// Don't show the 'library' namespace on official search results
export default ({ namespace, reponame } = {}) => {
  if (!namespace || !reponame) return '';
  return namespace === 'library' ? reponame : `${namespace}/${reponame}`;
};
