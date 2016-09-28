'use strict';
import map from 'lodash/collection/map';

const _encodeForm = (obj, prefix) => {
  const str = [];
  map(obj, (value, key) => {
    if (obj.hasOwnProperty(key)) {
      const k = prefix ? `${prefix}[${key}]` : key;
      str.push(typeof value === 'object' ?
        _encodeForm(value, k) :
        `${encodeURIComponent(k)}=${encodeURIComponent(value)}`);
    }
  });
  return str.join('&');
};

export default _encodeForm;
