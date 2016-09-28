Promise.when = (promises) => {
  if (!promises || !promises.length) {
    return Promise.resolve();
  }
  const succeeded = [];
  const failed = [];
  return new Promise((resolve, reject) => {
    const finish = () => {
      if (succeeded.length + failed.length < promises.length) {
        return;
      }
      if (failed.length) {
        reject(failed);
      } else {
        resolve(succeeded);
      }
    };
    promises.forEach((promise) => {
      promise.then(res => {
        succeeded.push(res);
        finish();
      }, err => {
        failed.push(err);
        finish();
      });
    });
  });
};
