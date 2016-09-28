import { Schema, arrayOf } from 'normalizr';

// The scan ID shouldn't use the ID attribute; we need to be able to load
// any scan by checking the repo:tag combination
const scan = new Schema('scan', {
  // idAttribute: (entity) => `${entity.reponame}:${entity.tag}`,
  idAttribute: 'tag',
});

// AKA layer
const blob = new Schema('blob', { idAttribute: 'index' });

// Key components by names to version numbers as they have no unique ID
const component = new Schema('component', {
  idAttribute: (entity) => `${entity.component}:${entity.version}`,
});

const vulnerability = new Schema('vulnerability', { idAttribute: 'cve' });

scan.define({ blobs: arrayOf(blob) });
blob.define({ components: arrayOf(component) });
component.define({ vulnerabilities: arrayOf(vulnerability) });

export const normalizers = {
  scan,
  blob,
  component,
  vulnerability,
};

/* latest_scan_status possible values */
export const COMPLETED = 'COMPLETED';
export const FAILED = 'FAILED';
export const IN_PROGRESS = 'IN_PROGRESS';
export const scanStatuses = [COMPLETED, FAILED, IN_PROGRESS];

/* layer_type possible values */
export const BASE = 'BASE';

/* vulnerability severities */
export const CRITICAL = 'critical';
export const MAJOR = 'major';
export const MINOR = 'minor';
export const SECURE = 'secure';
export const severities = [CRITICAL, MAJOR, MINOR, SECURE];
