'use strict';
import _ from 'lodash';

export function handleFormErrors(ctx, rawValueObject) {

  /**
   * This function expects a ctx that has a `fields`
   * object on the state, and a `_validate` function that
   * returns an object `{ hasError: bool, error: string }`
   *
   * A valid component looks like:
   *
   * var Component = React.createClasse({
   *   getInitialState() {
   *     return {
   *       fields: {}
   *     }
   *   },
   *   _validate(key, value) {
   *     return {
   *       hasError: true,
   *       error: 'It\'s always wrong!'
   *     }
   *   }
   * })
   */

    // shortcut keys for State
  let fields = ctx.state.fields || {};

    // loop through `rawValueObject`, validating values
    _.forIn(rawValueObject, function(value, key) {
      let { hasError, error } = ctx._validate(key, value);
        fields[key] = fields[key] || {};
        fields[key].hasError = hasError;
        fields[key].error = error;
    }, ctx);

    // queue up the new State
    ctx.setState({
      fields
    });
}
