import reducer, { DEFAULT_STATE } from 'reducers/internal';
import { INTERNAL_ROUTER_READY } from 'actions/internal';
import { expect } from 'chai';

describe('internal reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(undefined, {})).to.deep.equal(DEFAULT_STATE);
  });

  it('should handle INTERNAL_ROUTER_READY', () => {
    const action = { type: INTERNAL_ROUTER_READY };
    const state = reducer(undefined, action);
    expect(state.routerReady).to.equal(true);
  });
});
