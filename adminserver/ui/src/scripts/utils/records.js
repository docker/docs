'use strict';

import { Map, List } from 'immutable';
import mapValues from 'lodash/mapValues';

/**
 * normalizedToRecords({ 'user': UserRecord }, data)
 *
 * @param object Map of schema names to record constructors
 * @param object Normalized data
 */
export function normalizedToRecords(data, schemaMap = {}) {
  const { entities } = data;
  let normalizedEntities = {};

  Object.keys(entities).forEach((schemaName) => {
    // Extract the record constructor for this schema
    const Record = schemaMap[schemaName];

    // If this schema wasn't mapped to a Record constructor throw an error
    if (Record === undefined) {
      throw new Error(`Encountered unmapped schema type '${schemaName}'`);
    }

    // Iterate through every instance of this entity and convert to a record
    const records = mapValues(entities[schemaName], (item) => {
      return new Record(item);
    });

    // Convert the object of entity IDs to entities to an immutable map then set
    // this within the normalizedEntities object. This will preserve records as
    // Map doesn't deeply convert to a Map.
    normalizedEntities[schemaName] = new Map(records);
  });

  // We need to convert our data.result into either a map, list or leave as is
  // (for one single normalizr result).
  //
  // 1. If we use arrayOf() within normalizr this will be a list.
  // 2. If we normalize an object of schemas this will be a map
  // 3. Otherwie this will be a string
  let result = data.result;

  if (Array.isArray(data.result)) {
    result = new List(data.result);
  } else if (typeof result === 'object') {
    result = new Map(data.result);
  }

  return new Map({
    result: result, // Map, list or string.
    entities: new Map(normalizedEntities)
  });
}
