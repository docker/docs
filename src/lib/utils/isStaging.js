export default () => {
  if (document &&
      document.location &&
      document.location.hostname.indexOf('stage') !== -1) {
    return true;
  }

  if (process && process.env && process.env.RELEASE_STAGE === 'staging') {
    return true;
  }
  return false;
};
