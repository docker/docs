'use strict';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
import styles from './SimpleInput.css';
const { any, func, string, bool, oneOf } = PropTypes;
const debug = require('debug')('SimpleInput');

/*
 * A simpler version of DUXInput without errors / alert boxes
 * that is a full field (not a line)
 */
export default class SimpleInput extends Component {
  static propTypes = {
    autoFocus: bool,
    hasError: bool,
    name: string,
    onChange: func.isRequired,
    placeholder: string,
    readOnly: bool,
    type: oneOf('hidden text password email search'.split(' ')),
    value: any.isRequired
  }

  static defaultProps = {
    hasError: false,
    onChange() {
      debug('No onChange function set for Input');
    },
    value: '',
    placeholder: '',
    type: 'text',
    name: ''
  }

  focusInput = () => {
    findDOMNode(this.refs.input).focus();
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

  componentDidMount() {
    const { autoFocus } = this.props;
    if (autoFocus) {
      this.focusInput();
    }
  }

  render() {
    const { hasError,
            name,
            placeholder,
            readOnly,
            type,
            value } = this.props;
    const inputClass = hasError ? styles.error : styles.default;
    return (
      <div className={styles.inputDiv}>
        <input className={inputClass}
               name={name}
               onChange={this._onChange}
               placeholder={placeholder}
               readOnly={readOnly}
               ref='input'
               type={type}
               value={this.state.value || value} />
      </div>
    );
  }
}
