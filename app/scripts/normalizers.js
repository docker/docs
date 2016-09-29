'use strict';

import { Schema, arrayOf } from 'normalizr';

// Key repositories by the 'reponame' attribute instead of the 'key' field so
// that we can look up repositories by name in our reselect queries
const repository = new Schema('repository', { idAttribute: 'reponame' });

// The tag ID shouldn't be via key - it should be tagname
const tag = new Schema('tag', {
  idAttribute: (entity) => {
    // The natuilus API uses entity.tag as the tagname whereas hub uses
    // entity.name
    return (entity.tag) ? entity.tag : entity.name;
  }
});
// The scan ID shouldn't use the ID attribute; we need to be able to load
// any scan by checking the repo:tag combination
// TODO: Can we use the sha256sum here instead?
const scan = new Schema('scan', {
  idAttribute: (entity) => `${entity.reponame}:${entity.tag}`
});
// AKA layer
const blob = new Schema('blob', { idAttribute: 'index' });
// Key components by names to version numbers as they have no unique ID
const component = new Schema('component', {
  idAttribute: (entity) => `${entity.component}:${entity.version}`
});
const vulnerability = new Schema('vulnerability', { idAttribute: 'cve' });


// A repository has many tags
repository.define({
  tags: arrayOf(tag)
});

// A tag has many blobs and many scans
tag.define({
  blobs: arrayOf(blob),
  // NOTE: Right now a tag can only have the latest scan.
  // In the future we'll allow tags to have many scans
  scans: arrayOf(scan)
});

scan.define({
  blobs: arrayOf(blob)
});

blob.define({
  components: arrayOf(component)
});

component.define({
  vulnerabilities: arrayOf(vulnerability)
});

export {
  repository,
  tag,
  scan,
  blob,
  component,
  vulnerability
};
