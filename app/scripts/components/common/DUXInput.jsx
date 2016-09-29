'use strict';

import variantA from './DUXInput.css';
import variantB from './DUXInput-b.css';

import React, {
  Component,
  PropTypes
} from 'react';
import { findDOMNode } from 'react-dom';
const { any, func, string, bool, oneOf } = PropTypes;
import classnames from 'classnames';
import _ from 'lodash';
import AlertBox from 'common/AlertBox';
const debug = require('debug')('DUXInput');

export default class DUXInput extends Component {

  static propTypes = {
    className: string,
    error: string,
    hasError: bool.isRequired,
    name: string,
    onChange: func.isRequired,
    success: string,
    type: oneOf('hidden text password email search'.split(' ')),
    value: string,
    variant: string
  }

  static defaultProps = {
    onChange() {
      debug('No onChange function set for Input');
    },
    value: '',
    hasError: false,
    error: '',
    type: 'text',
    success: '',
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

  render() {
    /**
     * Theme variant is accepted through the props.
     */
    const { variant, name } = this.props;
    const styles = variant !== 'b' ? variantA : variantB;

    /**
     * If there is an error, display it
     */
    let maybeError = <span></span>;
    if(this.props.hasError && this.props.error) {
      maybeError = <AlertBox intent='alert'
                             onClick={this.focusInput}>{this.props.error}</AlertBox>;
    }

    /**
     * If there is a success message, display it
     */
    let maybeSuccess = <span></span>;
    if(this.props.success) {
      /**
       * TODO: This shouldn't be a property of the input in the way
       * that it currently exists. A field-level success state should
       * be very minimal.
       *
       * OLD_TODO: this could be an alert box with a close icon that can be
       * closed when the user wants to dismiss or time out after some
       * time (ideally i would like to see this as a notification and
       * that's it)
       */
      maybeSuccess = <AlertBox intent='success'>{this.props.success}</AlertBox>;
    }

    const groupClasses = classnames({
      [styles.group]: true,
      [styles.hasError]: this.props.hasError
    });

    return (
      <div className={groupClasses}>
        <input className={styles.duxInput}
               ref='input'
               name={name}
               type={this.props.type}
               value={this.state.value || this.props.value}
               onChange={this._onChange} />
        <span className={styles.bar}></span>
        <label className={styles.label}>{this.props.label}</label>
        {maybeError}
        {maybeSuccess}
      </div>
    );
  }
}
