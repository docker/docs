import React, { Component, PropTypes } from 'react';
import { Link } from 'react-router';
import css from './styles.css';
import { AngledTitleBox } from 'common';
import { DDC } from 'lib/constants/landingPage';
import { DDC_ID, DDC_TRIAL_PLAN } from 'lib/constants/eusa';
import routes from 'lib/constants/routes';
const { bool } = PropTypes;

/* eslint-disable max-len */
export default class DDCBanner extends Component {
  static propTypes = {
    isBetaPage: bool,
  }

  render() {
    const { isBetaPage } = this.props;
    const { title, DTR, UCP, CS } = DDC;
    const ddcDetail = routes.bundleDetail({ id: DDC_ID });
    const ddcPurchase = routes.bundleDetailPurchase({ id: DDC_ID });
    const ddcTrial = `${ddcPurchase}?plan=${DDC_TRIAL_PLAN}`;
    let featuredContent;
    if (isBetaPage) {
      featuredContent = (
        <AngledTitleBox
          title="Featured Docker Solution"
          className={css.titleBox}
        />
      );
    }
    const ddcWrapperClass = isBetaPage ? css.ddcWrapperBeta : css.ddcWrapper;
    const ddcSectionClass = isBetaPage ? css.ddcSectionBeta : css.ddcSection;
    const buttonClass = isBetaPage ? css.buttonBeta : css.button;
    return (
      <div className={ddcWrapperClass}>
        <div className="wrapped">
          {featuredContent}
          <div className={ddcSectionClass}>
            <div>
              <div className={css.titleAndButtons}>
                <div className={css.DDCtitle}>{title}</div>
                  <Link to={ddcTrial} className={buttonClass}>
                    30 Day FREE Evaluation
                  </Link>
                  <Link to={ddcDetail} className={buttonClass}>
                    Get Datacenter
                  </Link>
              </div>
            </div>
            <div className={css.DDCComponents}>
              <div className={css.DDC}>
                <div className={css.DDCName}>{DTR.name}</div>
                <div className={css.DDCDescription}>{DTR.description}</div>
              </div>
              <div className={css.DDC}>
                <div className={css.DDCName}>{UCP.name}</div>
                <div className={css.DDCDescription}>{UCP.description}</div>
              </div>
              <div className={css.DDC}>
                <div className={css.DDCName}>{CS.name}</div>
                <div className={css.DDCDescription}>{CS.description}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
/* eslint-enable */
