'use strict';

import { ATTEMPTING, ERROR, SUCCESS } from '../reduxConsts.js';

/**
 * SDK middleware for automatically calling SDK actions and storing request
 * statuses.
 *
 * Example:
 *
 * someAction({ tagName }) => ({
 *   type: 'SOME_ACTION',
 *   meta: {
 *     sdk: {
 *       call: SDK.func, // SDK function to call
 *       args: ['args', 'for', 'func'] // Args to pass in to SDK function
 *       callback: (err, res) => ({}) // SDK callback
 *       statusKey: ['SOMETHING'] // Unique identifier for saving status in status reducer
 *     }
 *   }
 * })
 *
 * NOTE: `statusKey` should be an array; the first item should namespace
 *       the action type and the second item should be unique to the particular
 *       record. For example, when deleting a tag:
 *
 *       statusKey: ['deleteRepoTag', 'latest']
 *
 *       Status will be stored in 'state.status.deleteRepoTag.latest'.
 *
 * NOTE: `args` does not include the SDK callback.
 *
 * For an example see actions/redux/tags.js
 */

// dispatchStatus takes the action and status of an SDK request and returns
// a new action to Redux for tracking state.
//
// The 'data' parameter may be either the error or response body from the call
const dispatchStatus = (action, status, data) => ({
  type: `${action.type}_STATUS`,
  payload: {
    // Add in everything from the initial action payload. This lets us pass
    // things such as namespaces and repo names to reducers which handle
    // success states (when deleting a tag we need the namespace/repo/tag name)
    ...action.payload,
    statusKey: action.meta.sdk.statusKey,
    status,
    data
  }
});

const sdkMiddleware = (store) => (next) => (action) => {
  // If there's no meta.sdk in our action we don't need to process it with our
  // middleware
  if (typeof action !== 'object' || !action.meta || !action.meta.sdk) {
    return next(action);
  }

  const { call, args, callback, statusKey } = action.meta.sdk;

  if (!statusKey) {
    throw new Error(`action.meta.sdk.statusKey is not defined for ${action.type}`);
  }

  // Wrap the callback with a function that automatically dispatches error
  // states for the SDK call to Redux.
  // Why: this eliminates the need to create error and success dispatches for every
  // action we create, and standardizes the format of all status dispatches
  const wrapped = (err, res) => {
    if (err) {
      next(dispatchStatus(action, ERROR, err));
      // TODO: Dispatch that there was an error with this call.
    } else {
      next(dispatchStatus(action, SUCCESS, res));
    }

    // Ensure we call the original callback supplied for the SDK call.
    if (callback) {
      callback.apply(null, [err, res]);
    }
  };

  // Dispatch that we're attempting the SDK call
  next(dispatchStatus(action, ATTEMPTING));

  // Make the SDK call here.
  call.apply(null, [...args, wrapped]);

  return next(action);
};

export default sdkMiddleware;
