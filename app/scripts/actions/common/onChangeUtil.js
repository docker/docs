'use strict';

import updateFormField from './updateFormField';

/**
 * This function is a generic `onChange` Handler. Typical use is as follows:
 *
 * ```
 * import onChange from 'this/file/here';
 *
 * createClass({
 *   _onChange: onChange({ storePrefix: 'ENTERPRISE_TRIAL' }),
 *   render() {
 *     return (
 *       <DUXInput onChange={this._onChange('first_name').bind(this)} />
 *     )
 *   }
 * });
 * ```
 *
 * `function setComponentOnChangeHandler` is the one that is `bind`ed.
 *
 * storePrefix is `ENTERPRISE_TRIAL` in the stores handlers object:
 *
 * ```
 * handlers: {
 *   ENTERPRISE_TRIAL_CLEAR_FORM: '_clearForm',
 *   ENTERPRISE_TRIAL_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
 *   ENTERPRISE_TRIAL_ATTEMPT_START: '_enterpriseTrialAttemptStart',
 *   ENTERPRISE_TRIAL_BAD_REQUEST: '_badRequest',
 *   ENTERPRISE_TRIAL_SUCCESS: '_enterpriseTrialSuccess'
 * }
 * ```
 *
 * e is the onChange event from the form field.
 */

export default function({ storePrefix }) {
  return function setComponentOnChangeHandler(fieldKey) {
    return (e) => {
      this.context.executeAction(updateFormField({ storePrefix }), {
        fieldKey,
        fieldValue: e.target.value
      });
    };
  };
}

