import React, { Component, PropTypes } from 'react';
import classnames from 'classnames';
import { XLARGE } from 'lib/constants/sizes';
import {
  BackButtonArea,
  Card,
  CheckIcon,
} from 'components/common';

import css from './styles.css';
const { func } = PropTypes;

export default class ProductSubmitted extends Component {

  static propTypes = {
    onClickBack: func.isRequired,
  }

  render() {
    const { onClickBack } = this.props;
    return (
      <div className={classnames({ [css.details]: true, wrapped: true })}>
        <Card className={css.card}>
          <div>
            <div className={css.confirmed}>
              <CheckIcon size={XLARGE} className={css.check} />
            </div>
            <div className={css.confirmedMessage}>
              <p>Your product has been submitted for approval!</p>
              <p>
                You can still go back to update and resubmit
                any changes and fixes.
              </p>
              <p>
                We&#39;ll be in touch with next steps soon!
              </p>
            </div>
          </div>
        </Card>
        <BackButtonArea
          onClick={onClickBack}
          className={css.backButton}
          text="Previous Step"
        />
      </div>
    );
  }
}
