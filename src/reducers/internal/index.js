import { INTERNAL_ROUTER_READY } from 'actions/internal';
import { cloneDeep } from 'lodash';

export const DEFAULT_STATE = {
  routerReady: false,
};

export default function internal(state = DEFAULT_STATE, action) {
  const nextState = cloneDeep(state);
  switch (action.type) {
    case INTERNAL_ROUTER_READY:
      nextState.routerReady = true;
      return nextState;
    default:
      return state;
  }
}
