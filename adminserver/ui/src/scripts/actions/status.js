'use strict';

import consts from 'consts';

export function resetStatus() {
  return {
    type: consts.status.RESET_STATUS,
    payload: {}
  };
}
