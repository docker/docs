'use strict';

import React, {
  PropTypes,
  Component
  } from 'react';
const { string, number, shape, bool, func } = PropTypes;
import Card, { Block } from '@dux/element-card';
import { Button } from 'dux';
import DUXInput from '../common/DUXInput.jsx';
import validateCouponCode from '../../actions/validateCouponCode.js';
import clearCoupon from '../../actions/clearCloudCoupon.js';
import styles from './CouponForm.css';

export default class CouponForm extends Component {
  static propTypes = {
    coupon: shape({
      couponCode: string,
      discountValue: number,
      hasError: bool
    })
  }
  state = {
      couponCode: ''
  }
  static contextTypes = {
    executeAction: func.isRequired
  }

  _updateCoupon = (e) => {
    e.preventDefault();
    this.setState({couponCode: e.target.value});
  }
  _clearCoupon = (e) => {
    e.preventDefault();
    this.setState({couponCode: ''});
    this.context.executeAction(clearCoupon);
  }
  validateCoupon = (plan) => {
    return (e) => {
      e.preventDefault();
      this.context.executeAction(validateCouponCode, {coupon_code: this.state.couponCode, plan: plan});
      this.setState({couponCode: ''});
    };
  }

  render() {
    const {coupon} = this.props;
    let discount;
    let couponApplied;
    if (coupon.discountValue > 0) {
      discount = '- $' + coupon.discountValue;
      couponApplied = (
        <div className="row">
          <div className="columns large-5">
            Coupon Applied:
          </div>
          <div className={'columns large-4 end ' + styles.couponPrice}>
            {coupon.couponCode}
          </div>
        </div>
      );
    } else {
      discount = '-----';
    }
    return (
      <Card>
        <Block>
          <div className="row">
            <div className="columns large-5">
              Plan Cost:
            </div>
            <div className={'columns large-4 end ' + styles.couponPrice}>
              $150
            </div>
          </div>
          <div className="row">
            <div className="columns large-5">
              Coupon:
            </div>
            <div className={'columns large-4 end ' + styles.couponPrice}>
              {discount}
            </div>
          </div>
          {couponApplied}
          <hr />
          <div className={'row ' + styles.total}>
            <div className="columns large-5">
              Total Charge:
            </div>
            <div className={'columns large-4 end ' + styles.couponPrice}>
              ${150 - this.props.coupon.discountValue}
            </div>
          </div>
          <div className="row">
            <form className={styles.couponCode} onSubmit={this.validateCoupon('cloud_starter')}>
              <div className="columns large-8">
                <DUXInput label='Coupon Code'
                          error='Invalid Coupon'
                          hasError={this.props.coupon.hasError}
                          onChange={this._updateCoupon}
                          value={this.state.couponCode}/>
              </div>
              <div className="columns large-4">
                <Button type="submit" size='small'>Add</Button>
              </div>
            </form>
          </div>
          <div className={'row ' + styles.clear}>
            <a onClick={this._clearCoupon}>clear coupon</a>
          </div>
        </Block>
      </Card>
    );
  }
}
