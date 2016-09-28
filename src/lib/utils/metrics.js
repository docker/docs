import isStaging from 'lib/utils/isStaging';
import isDev from 'lib/utils/isDevelopment';

// Bugsnag Notification
export default (title, message) => {
  if (Bugsnag) {
    if (isStaging()) {
      Bugsnag.releaseStage = 'staging';
    } else if (isDev()) {
      Bugsnag.releaseStage = 'development';
    } else {
      Bugsnag.releaseStage = 'production';
    }
    return Bugsnag.notify(title, message);
  }
  return {};
};
