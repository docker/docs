'use strict';

import consts from 'consts';

/**
 * Shape of actions should be as follows:
 *
 *   {
 *     type: SAVE_PROFILE
 *     meta: {
 *       promise: api.SaveProfile(data)
 *       notifications: {
 *         pending: "We're saving your shit now",
 *         success: "something happened really well",
 *         error: "what the fuck?!"
 *       }
 *     }
 *   }
 *
 */
const notificationMiddleware = () => next => action => {
  if (typeof action !== 'object' || !action.meta || !action.meta.promise || !action.meta.notifications) {
    return next(action);
  }

  // Store the notification ID we create in the closure, allowing us to update
  // the same notification after the promise has resolved
  let notificationId;

  // also store the timeout ID in a closure, allowing us to update the timeout
  // when the notification gets updated
  let timeoutId;

  const { meta: { notifications } } = action;

  /**
   * Create or update a notification with the given title and message
   *
   * @param string|function  error string or func which returns an error string
   * @param string status    `status` param of notification
   * @param object data      data from AJAX rquest passed into message func
   */
  const makeNotification = (message, title, status, data = {}) => {
    if (!message) {
      return null;
    }

    // if no title is provided, just use the status (capitalized)
    title = title || status.charAt(0).toUpperCase() + status.slice(1);

    let type = consts.UPDATE_GROWL_NOTIFICATION;
    if (!notificationId) {
      type = consts.ADD_GROWL_NOTIFICATION;
      notificationId = `promise-${Math.floor(Math.random() * 100000000).toString(16)}`;
    }

    if (typeof message === 'function') {
      message = message(data);
    }

    // since we are updating, clear the timeout id so the last timeout doesnt mess with the new notification
    clearInterval(timeoutId);
    if (!data.persistent) {
      timeoutId = setTimeout( () => {
        next({
          type: consts.REMOVE_GROWL_NOTIFICATION,
          payload: {
              id: notificationId
          }
        });
      }, consts.notification.DURATION);
    }

    return {
      type,
      payload: {
        id: notificationId,
        title,
        message,
        status,
        img: data.img
      }
    };
  };

  if (notifications.pending) {
      next(makeNotification(notifications.pending, null, 'info'));
  }

  action.meta.promise = action.meta.promise.then(
    payload => {
      if (notifications.success) {
          next(makeNotification(notifications.success, null, 'success', payload));
      }
      return payload;
    },
    error => {
      if (notifications.error) {
          next(makeNotification(notifications.error, null, 'alert', error));
      }
      throw error;
    }
  );
  delete (action.meta.notifications);
  return next(action);
};

export default notificationMiddleware;
