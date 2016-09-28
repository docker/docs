'use strict';

import { bindActionCreators } from 'redux';
import React from 'react';
import consts from './reduxConsts';
/**
 * Given an object of actions, this returns a thunk which returns all actions
 * bound to dispatch using the same key names.
 *
 * This allows us to use `this.props.actions.$actionName` within components
 * after being connected to Redux.
 *
 * Example:
 *
 *   @connect(mapState, mapActions(Actions))
 *   class Basic extends Component {
 *     // this.props.actions now containers all keys in Actions bound to dispatch
 *   }
 *
 */
export let mapActions = (actions) => {
  return (dispatch) => { return { actions: bindActionCreators(actions, dispatch) }; };
};

export const mapToArray = (map) => {
  return Object.keys(map).map(key => map[key]);
};
