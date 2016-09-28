import uuidRegExp from './regexp-uuid';

export function extractUuid(resource_uri) {
  const [uuid = ''] = uuidRegExp.exec(resource_uri) || [];
  return uuid;
}
