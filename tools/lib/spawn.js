const path = require('path');
const { spawn } = require('child_process');

module.exports = function promiseSpawn(cmd, args = [], options = {}) {
  return new Promise((resolve, reject) => {
    Object.assign(
      options,
      {
        cwd: path.resolve(__dirname, '../..'),
        stdio: [null, process.stdout, process.stderr],
      }
    );

    const child = spawn(cmd, args, options);

    child.on('error', reject);

    child.on('exit', (code) => {
      if (code !== 0) {
        reject(new Error('Process failed'));
      } else {
        resolve();
      }
    });
  });
};
