'use strict';

import Axios from 'axios';
import assign from 'object-assign';

import consts from 'consts';

/**
 * Notification actions.
 *
 * This store listens for 'notification' events from AJAX requests, such as
 * 'saving', 'saved', and 'save failed'.
 *
 * There are two types of notifications, growl and banner
 *
 * Each growl notification MUST be in the format:
 *
 *  {
 *    optional enum   `class` { 'alert', 'warning', 'success', 'info'} [default='info'];
 *    required string `message`;
 *    optional string `title`; // big bold text, defaults to capitalized status
 *    optional string `url`;        // If clickable, the URL to route to
 *    optional bool   `persistent`; // False to hide after 5 seconds
 *  }
 *
 * Example:
 *
 *   addNotification({
 *     class: 'alert',
 *     message: 'Youre doing it wrong',
 *     title: 'Shit went down!',
 *     url: '/admin/settings',
 *     persistent: false,
 *   })
 *
 * Each banner notification MUST be in the format:
 *
 *  {
 *    required string `message`;
 *    optional string `url`;        // If clickable, the URL to route to
 *    optional string `img`;        // the url of an image resource to display left of the text in the banner
 *  }
 *
 * Example:
 *
 *   addNotification({
 *     message: 'GC in progress',
 *     url: '/admin/settings/gc',
 *     img: '/public/img/broom-white.png',
 *   })
 */

let defaultId = 1;

const defaults = {
  class: 'info',
  persistent: false
};


/**
 * Helper method to throw errors with no message and add defaults to our
 * notification.
 *
 * This adds an incremental ID to each particular notification as it is created
 * if an ID does not exist.  You may supply your own ID as a reference for
 * a notification to update it.
 */
function normalize(data) {
  if (!data.message) {
    throw new Error('The "message" key must be provided when adding a notification');
  }
  defaultId++;
  return assign({}, defaults, { defaultId }, data);
}

export function addGrowlNotification(data) {
  return dispatch => {
    data = normalize(data);

    dispatch({
      type: consts.ADD_GROWL_NOTIFICATION,
      payload: data
    });

    if (!data.persistent) {
      window.setTimeout( () => {
        dispatch({
          type: consts.REMOVE_GROWL_NOTIFICATION,
          payload: {
              id: data.id
          }
        });
      }, consts.notification.DURATION);
    }
  };
}

export function updateGrowlNotification(id, data) {
  return {
    type: consts.UPDATE_GROWL_NOTIFICATION,
    payload: {
      id: id,
      data: data
    }
  };
}

export function removeGrowlNotification(id) {
  return {
    type: consts.REMOVE_GROWL_NOTIFICATION,
    payload: { id }
  };
}

export function addBannerNotification(data) {
  return dispatch => {
    data = normalize(data);

    dispatch({
      type: consts.ADD_BANNER_NOTIFICATION,
      payload: data
    });
  };
}

export function updateBannerNotification(id, data) {
  return {
    type: consts.UPDATE_BANNER_NOTIFICATION,
    payload: {
      id: id,
      data: data
    }
  };
}

export function removeBannerNotification(id) {
  return {
    type: consts.REMOVE_BANNER_NOTIFICATION,
    payload: { id }
  };
}

/**
 * Used to dismiss successful upgrade notifications which are emitted via
 * server-side logic.
 *
 * These server-side notifications which are embedded into window.notifications
 * can contain an 'onclose' key which stores a GET parameter which, when
 * accessed, prevents the notification from being displayed agian.
 *
 */
export function onClose(query) {
  Axios.get(query);
  // This is a no-op
  return { type: null };
}
