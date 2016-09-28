import map from 'lodash/map';

export default function encodeForm(obj, prefix) {
  const str = [];
  map(obj, (value, key) => {
    if (obj.hasOwnProperty(key)) {
      const k = prefix ? `${prefix}[${key}]` : key;
      str.push(typeof value === 'object' ?
        encodeForm(value, k) :
        `${encodeURIComponent(k)}=${encodeURIComponent(value)}`);
    }
  });
  return str.join('&');
}
