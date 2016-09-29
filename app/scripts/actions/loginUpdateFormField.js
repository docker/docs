'use strict';

export default function({ dispatch },
                        { fieldKey, fieldValue },
                        done) {
  dispatch('LOGIN_UPDATE_FIELD_WITH_VALUE', { fieldKey, fieldValue });
}
