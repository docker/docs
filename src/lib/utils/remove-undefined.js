// Utility that can be used for on a map removing any keys that are set with
// `undefined` values
import omitBy from 'lodash/omitBy';

export default (map) => {
  return omitBy(map, (value) => typeof value === 'undefined');
};
