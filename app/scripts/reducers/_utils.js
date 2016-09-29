'use strict';

import Immutable from 'immutable';

/**
 * mergeEntity takes an entity name and merges it into the current state
 * if found.
 *
 * This is used when your state contains a basic map of entities.
 *
 * Examples:
 *
 *   mergeEntity('repository'):
 *      > merge action.payload.entities.repository into the current state
 *
 */
export const mergeEntity = (entityType) => (state, action) => {
  //TODO: Remove promises stuff with ready / error?
  const { payload, ready, error } = action;
  if (!ready || error || !payload.entities[entityType]) {
    return state;
  }
  return state.merge(new Immutable.Map(payload.entities[entityType]));
};

export const mapToRecord = (map, Record) => {
  let records = {};
  Object.keys(map).forEach(item => { records[item] = new Record(map[item]); });
  return records;
};

export const mergeEntities = (...entities) => (state, action) => {

  const { payload, ready, error } = action;
  if (!ready || error) {
    return state;
  }

  return state.withMutations( map => {
    entities.forEach( item => {
      return map.mergeIn([item], new Immutable.Map(payload.entities[item]));
    });
    return map;
  });
};
