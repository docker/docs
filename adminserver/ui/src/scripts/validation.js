'use strict';

import { get, set } from 'lodash/object';

// given a list of required keys return a list of keys missing input values
export const checkRequiredFields = (requiredKeys, formInputValues) => {
    return requiredKeys.filter((key) => {
        // explicit check for undefined & empty because formInputValues[key] could be bool false
        if (formInputValues[key] === undefined || formInputValues[key] === '') {
            return key;
        }
    });
};

export const reportMissingFields = (missingFields, errors) => {
    missingFields.map((missingField) => {
        errors[missingField] = 'Required. ';
    });
    return errors;
};

// Adapted from from: https://github.com/erikras/react-redux-universal-hot-example/blob/master/src/utils/validation.js
const isEmpty = value => value === undefined || value === null || value === '';
const join = rules => (value, values) => rules.map(rule => rule(value, values)).filter(error => !!error)[0 /* first error */];

export function required(value) {
  if (isEmpty(value)) {
    return 'Required';
  }
}

export function minLength(min) {
  return value => {
    if (!isEmpty(value) && value.length < min) {
      return `Must be at least ${min} characters`;
    }
  };
}

export function maxLength(max) {
  return value => {
    if (!isEmpty(value) && value.length > max) {
      return `Must be no more than ${max} characters`;
    }
  };
}

export function integer(value) {
  if (!Number.isInteger(Number(value))) {
    return 'Must be an integer';
  }
}

export function oneOf(enumeration) {
  return value => {
    if (!~enumeration.indexOf(value)) {
      return `Must be one of: ${enumeration.join(', ')}`;
    }
  };
}

export function regex(pattern, errorMessage = '') {
  let re = pattern;
  if (!(re instanceof RegExp)) {
    re = new RegExp(pattern);
  }
  return value => {
    if (!re.test(value)) {
      if (errorMessage) {
        return errorMessage;
      }
      return `Doesn't match regex: ${pattern}`;
    }
  };
}

export function json(value) {
  try {
    JSON.parse(value);
    return false;
  } catch (e) {
    return 'Invalid JSON';
  }
}

/**
 * Matches whether a field's value is exactly as specified
 *
 */
export const is = match => value => (value !== match) && `Must match ${match}`;

/**
 * Run the validation rules defined in `test` for a given field if the given
 * predicate is truthy.
 *
 * Examples:
 *
 *   ldapDN: [ onlyIf(data => data.type !== 'ldap', required) ]
 *
 * This only runs the `required` rule in the `ldapDN` field if the field 'type'
 * matches 'ldap'.
 *
 * Each function in onlyIf is passed *all* form values as the first argument;
 * you must compose new validation functions and pass in the field to test
 */
export const onlyIf = (predicate, test) => (value, data) => predicate(data) && test(value);

/**
 * Returns a validator for use with redux form.
 * Usage:
 *
 * @reduxForm({
 *   form: 'key',
 *   fields: ['name', 'email'],
 *   validate: createValidator([
 *     name: [required, anotherValidator]
 *   })
 * })
 *
 */
export function createValidator(rules) {
  // TODO memoize - this will be ran every time a parent component recomputes,
  // which with our N+1 calls might be offfften
  return (data = {}) => {
    let errors = {};
    Object.keys(rules).forEach((key) => {
      const rule = join([].concat(rules[key])); // concat enables both functions and arrays of functions

      // If the key contains square brakcets we're validating many fields at
      // once.
      // TODO: allow many embedded fields (ie 'key[].some[].thing')
      if (key.indexOf('[]') === -1) {
        // This is easy - it's a single form field.
        // Pass in the fields' value as the first argument and all values as the
        // second argument. This lets us validate fields depending on other
        // field's values
        const error = rule(get(data, key), data);
        if (error) {
          // set nested errors such as like parent.child
          set(errors, key, error);
          errors.valid = false;
        }
        return errors;
      }

      // Now we have many fields to validate.
      const [root, leaf] = key.split('[]');
      get(data, root).forEach((item, idx) => {
        const error = rule(get(item, leaf), data);
        if (error) {
          // set the error in `item[0].whatever`
          errors[root] = errors[root] || [];

          const newKey = key.replace('[]', `.${idx}`, 1);
          set(errors, newKey, error);
          errors.valid = false;
        }
        return errors;
      });
    });
    return errors;
  };
}
