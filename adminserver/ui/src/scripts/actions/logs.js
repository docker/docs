'use strict';

import { Logs } from 'dtr-js-sdk';
import consts from 'consts';
import Rx from 'rx';

// if ("production" !== process.env.NODE_ENV) {
//   require("./__tests__/mocks/logs");
// }

const INTERVAL = 2500;

// These store reactive observables for long polling new logs
let subscription, observable;

// Stores the previous time for observing logs, allowing us to select logs since
// the last fetch
let previousTime;

/**
 * observeLogs creates a new reactive observable stream which long-polls the log
 * endpoint for logs within the last INTERVAL period.  Upon retrieving logs this
 * dispatches an action with the new log lines.
 *
 * It also stores the observer in the local `observable` variable, allowing us
 * to dispose() the observable from within another action.
 *
 * @param string serverName  server name to query logs for
 */
export function observeLogs(serverName) {
  return dispatch => {
    if (subscription) {
      subscription.dispose();
    }
    observable = Rx.Observable
      .interval(INTERVAL)
      .flatMap(() => Rx.Observable.fromPromise(Logs.getLogs({ serverName, direction: 'after', since: previousTime, max: 500 })))
      .retry(60)
      .map( response => { return (typeof response.data === 'string') ? [] : response.data.lines.reverse(); } );

    subscription = observable.subscribe(
      lines => {
        // Update the time checked
        previousTime = new Date();

        dispatch({
          type: consts.log.OBSERVE_LOG_DATA,
          data: { lines }
        });
      },
      () => {}
    );
    dispatch({
      type: consts.log.OBSERVE_LOGS
    });
  };
}

/**
 * getLogs retrieves an initial number of logs before now.
 *
 * @arg string serverName Name of the server to retrieve logs for
 */
export function getLogs(serverName) {
  return dispatch => {
    dispatch({
      type: consts.log.FETCH_LOGS,
      payload: { serverName },
      meta: {
        promise: Logs.getLogs({ serverName, direction: 'before', since: new Date(), max: 500 })
          .then( response => {
            let lines = [];
            if (typeof response.data !== 'string') {
              lines = response.data.lines.reverse();
            }

            // Store the current time.
            previousTime = new Date();

            return { serverName, lines };
          })
      }
    });
    dispatch(observeLogs(serverName));
  };
}

/**
 * unobserveLogs disposes the current subscription and observer, removing the
 * current long-polling of logConsts.
 *
 * This should be called from any component on componentWillUnmount() if the
 * component is observing logConsts.
 */
export function unobserveLogs() {
  if (subscription) {
    subscription.dispose();
    subscription = null;
  }
  return {
    type: consts.log.UNOBSERVE_LOGS
  };
}
