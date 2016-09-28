'use strict';

import Axios from 'axios';
import { Observable } from 'rx';

import consts from 'consts';

/**
 * Check for updates once an hour
 *
 */
const INTERVAL = 60 * 60 * 1000;

let subscription;

export function pollForUpdates() {
  return dispatch => {
    if (subscription) {
      return;
    }

    let query = Observable
      .fromPromise(Axios.get(consts.URLS.UPDATES));

    let poll = Observable
      .interval(INTERVAL)
      .flatMap( () => { return Observable.fromPromise(Axios.get(consts.URLS.UPDATES)); })
      .retry();

    subscription = Observable.concat(query, poll).subscribe(
      response => {
        if (!response.data.upgradeAvailable) {
          return;
        }

        dispatch({
          type: consts.updates.UPDATE_AVAILABLE,
          payload: response.data
        });

        let notification = {
          id: consts.updates.UPDATE_NOTIFICATION_ID,
          message: 'An update to Trusted Registry is available. Click here for more information.',
          url: '/admin/settings/updates'
        };
        dispatch({ type: consts.ADD_BANNER_NOTIFICATION, payload: notification });
      }
    );
  };
}

export function updateTo(toVersion) {
  return dispatch => {
    let notification = {
      id: consts.updates.UPDATE_NOTIFICATION_ID,
      message: `Starting update process to ${toVersion}`,
      persistent: true
    };
    dispatch({ type: consts.ADD_BANNER_NOTIFICATION, payload: notification });

    Axios.post(consts.URLS.UPDATES, {version: toVersion})
      .then(() => {
        notification.message = 'Server is now updating';
        dispatch({
          type: consts.UPDATE_BANNER_NOTIFICATION,
          payload: notification
        });
        dispatch({ type: consts.reload.UPGRADING });
      })
      .catch( response => {
        notification.message = 'There was an error updating';
        notification.class = 'alert';
        dispatch({
          type: consts.UPDATE_BANNER_NOTIFICATION,
          payload: notification
        });
        dispatch({
          type: consts.updates.UPDATE_FAILURE,
          payload: response
        });
      });
  };
}
