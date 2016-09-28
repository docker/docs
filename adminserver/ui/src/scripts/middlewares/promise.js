'use strict';

/**
 * Shape of actions should be as follows:
 *
 *   {
 *     type: FETCH_TEAM,
 *     meta: {
 *       promiseScope: [orgName, teamName], // Optional, defaults to []
 *       promise: team.getTeam(orgName, teamName)
 *     }
 *   }
 *
 * This action will be intercepted and this will be fired instead:
 *   {
 *     type: FETCH_TEAM,
 *     ready: false,
 *     meta: {
 *       promiseScope: [orgName, teamName],
 *       promiseId: 1234
 *     }
 *   }
 * The status at state.status.orgName.teamName.status will be set to PENDING
 *
 * When the promise is resolved:
 *   {
 *     type: FETCH_TEAM,
 *     ready: true,
 *     payload: resolvedValue,
 *     meta: {
 *       promiseScope: [orgName, teamName],
 *       promiseId: 1234
 *     }
 *   }
 * The status at state.status.orgName.teamName.status will be set to SUCCESS
 *
 * When the promise is rejected:
 *   {
 *     type: FETCH_TEAM,
 *     ready: true,
 *     error: rejectedValue,
 *     meta: {
 *       promiseScope: [orgName, teamName],
 *       promiseId: 1234
 *     }
 *   }
 * The status at state.status.orgName.teamName.status will be set to FAILURE
 *
 * If multiple promise actions are dispatch for the same action type and scope,
 * the status is for the very last promise action dispatched. Furthermore, the
 * stale promises will not fire their resolve/reject actions.
 */
const promiseMiddleware = store => next => action => {
  if (typeof action !== 'object' || !action.meta || !action.meta.promise) {
    return next(action);
  }

  const promiseId = action.meta.promise._id || Math.ceil(Math.random() * (1 << 30));

  let makeAction = (ready, data) => {
    let newAction = {
      ...action,
      ready,
      ...data,
      meta: {
        ...action.meta,
        promiseId
      }
    };
    if (!newAction.meta.promiseScope) {
      newAction.meta.promiseScope = [];
    }
    delete newAction.meta.promise;
    return newAction;
  };

  // TODO should record the length of promiseScope after first pendingAction and scoping is always at same depth for same action types?
  let pendingAction = makeAction(false);
  next(pendingAction);

  return action.meta.promise.then(
    payload => {
      try {
        next(makeAction(true, { payload }));
      } catch(e) {
        console.error(e);
      }
      return payload;
    },
    error => {
      try {
        let rejectedAction = makeAction(true, { error });
        const lastPendingPromiseId = store.getState().status.getIn([rejectedAction.type, ...rejectedAction.meta.promiseScope, 'promiseId']);
        if (lastPendingPromiseId === rejectedAction.meta.promiseId) {
          next(rejectedAction);
        } else {
          if (process.env.NODE_ENV !== 'production') {
            console.log('Not firing stale action', lastPendingPromiseId, rejectedAction);
          }
        }
      } catch(e) {
        console.error(e);
        throw e;
      }
      throw error;
    }
  );

};

export default promiseMiddleware;
