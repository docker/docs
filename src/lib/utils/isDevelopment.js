export default () => {
  if (document &&
      document.location &&
      document.location.hostname === 'localhost') {
    return true;
  }

  return false;
};
