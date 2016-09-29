'use strict';

import { Record } from 'immutable';

// Records are the same as Maps but with accessors
// and can only have these defined fields set.
// USE: Instead of `shape`, in propTypes, we can use
//      status: instanceOf(StatusRecord)
//
// NOTE: All records should be defined in this file
export const StatusRecord = new Record({
  status: '',
  error: undefined
});
