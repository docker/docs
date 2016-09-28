'use strict';

/**
 * This function dispatches an event that updates a form backed by a form store.
 * @param {string} storePrefix - the prefix for the form store events. This
 * scopes the event to a particular store.
 *
 * The storePrefix usually takes the form of a shortened version of the store's
 * name. In the case of the `EnterpriseTrialFormStore`, given handlers as such:
 *
 * ```
 * handlers: {
 *   ENTERPRISE_TRIAL_UPDATE_FIELD_WITH_VALUE: 'updateFormField'
 * }
 * ```
 *
 * The store prefix for the preceeding example is `ENTERPRISE_TRIAL`.
 *
 * ---
 *
 * Usage of this module in a component might look like:
 *
 * ```
 * import updateFormField from 'wherever';
 *
 *
 * var Whatever = createClass({
 *   _onChange(fieldKey) {
 *     return (e) => {
 *       this.context.executeAction(updateFormField({
 *         storePrefix: 'ENTERPRISE_TRIAL'
 *       }), {
 *         fieldKey,
 *         fieldValue: e.target.value
 *       });
 *     }
 *   },
 *   render() {
 *     return (
 *       <input onChange={this._onChange}/>
 *     )
 *   }
 * });
 */
export default function({storePrefix}){
  return (actionContext, { fieldKey, fieldValue}) => {
    actionContext.dispatch(`${storePrefix}_UPDATE_FIELD_WITH_VALUE`, { fieldKey, fieldValue });
  };
}

