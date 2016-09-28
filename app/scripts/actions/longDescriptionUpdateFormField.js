'use strict';

export default function({ dispatch },
                        {
                          fieldKey,
                          fieldValue
                        },
                        done) {
  dispatch('LONG_DESCRIPTION_UPDATE_FIELD_WITH_VALUE', {
    fieldKey, fieldValue
  });
}
