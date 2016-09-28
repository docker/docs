'use strict';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
import classnames from 'classnames';
import styles from './FancyInput.css';
const { any, func, string, bool, oneOf } = PropTypes;
const debug = require('debug')('SimpleInput');

export default class FancyInput extends Component {
  static propTypes = {
    autoFocus: bool,
    hasError: bool,
    error: string,
    name: string,
    onChange: func.isRequired,
    placeholder: string,
    readOnly: bool,
    type: oneOf('hidden text password email search'.split(' ')),
    value: any.isRequired,
    variant: oneOf(['white'])
  }

  static defaultProps = {
    hasError: false,
    error: '',
    onChange() {
      debug('No onChange function set for Input');
    },
    value: '',
    placeholder: '',
    type: 'text',
    name: ''
  }

  state = {
    value: ''
  }

  _onChange = (e) => {
    // Set local state
    const { value } = e.target;
    this.setState({ value });

    // Pass control to onChange controller. Usually a form.
    this.props.onChange(e);
  }

  componentWillReceiveProps(props) {
    /**
     * State should be tracked globally in a parent component. As
     * such, we need to check the passed in prop against the current
     * (local) state value to avoid overwriting the input value.
     *
     * This prevents a bug which resulted in the cursor being "jumped"
     * to the end of an input after every character.
     */
    if (props.value !== this.state.value) {
      this.setState({
        value: props.value
      });
    }
  }

  render() {
    const { hasError,
            error,
            name,
            placeholder,
            readOnly,
            type,
            value,
            variant,
            autofocus } = this.props;
    const groupClass = classnames({
      [styles.inputDiv]: true,
      [styles.white]: variant === 'white',
      [styles.hasError]: hasError
    });

    return (
      <div className={groupClass}>
        <input className={styles.default}
               name={name}
               autofocus={autofocus}
               onChange={this._onChange}
               placeholder={placeholder}
               readOnly={readOnly}
               type={type}
               value={this.state.value || value} />
             <span className={styles.bar}></span>
             <div className={styles.error}>{error}</div>
      </div>
    );
  }
}
