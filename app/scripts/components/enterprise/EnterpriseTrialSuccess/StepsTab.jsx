'use strict';

import React, { PropTypes, Component } from 'react';
import { Link } from 'react-router';
import styles from '../EnterpriseTrialSuccess.css';
import classnames from 'classnames';
import FA from 'common/FontAwesome';
const { number, string } = PropTypes;

export default class StepsTab extends Component {
  static propTypes = {
    currentStep: number.isRequired,
    namespace: string.isRequired,
    step: number.isRequired,
    title: string.isRequired
  }

  render() {
    const { currentStep,
            namespace,
            step,
            title } = this.props;
    const classes = classnames({
      [styles.tab]: true,
      [styles.active]: currentStep === step,
      [styles.success]: currentStep > step,
      [styles.last]: step === 4
    });
    let icon;
    //TODO: replace the <i> with line below when FA can accept size=1x
    // <FA icon='fa-check' invert={true} stack={true} size='1x' />
    if (currentStep > step) {
      icon = (
        <span className="fa-stack">
          <FA icon='fa-circle' stack={true} size='2x' />
          <i className="fa fa-check fa-stack-1x fa-inverse"></i>
        </span>
      );
    } else {
      icon = (
        <span className="fa-stack">
          <i className="fa fa-circle-thin fa-stack-2x"></i>
          <strong className="fa-stack-1x">{step}</strong>
        </span>
        );
    }
    return (

      <Link className={ classes } to='/enterprise/trial/success/' query={{namespace, step}}>
          { icon } { title }
      </Link>
    );
  }
}
