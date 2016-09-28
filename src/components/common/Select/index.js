import React, { Component, PropTypes } from 'react';
import ReactSelect from 'react-select';
import styles from './styles.css';
import classnames from 'classnames';
import { isEmpty } from 'lodash';
import 'react-select/dist/react-select.css';

const { string } = PropTypes;

function create(NewSelect) {
  return class Select extends Component {
    static propTypes = {
      className: string,
      errorText: string,
    }
    render() {
      const {
        className = '',
        errorText = '',
      } = this.props;
      const errorClassNames = classnames(
        styles.errorText,
        {
          [styles.visible]: !isEmpty(errorText),
        }
      );
      return (
        <div>
          <NewSelect {...this.props} className={`dselect ${className}`} />
          <div className={errorClassNames}>
            {errorText}
          </div>
        </div>
      );
    }
  };
}

const Select = create(ReactSelect);
Select.Async = create(ReactSelect.Async);

export default Select;
