import { Decimal, toMoney } from '../../lib/money.js';

/**
 * Costing strategy. The engine asks it how much an issue costs and what the
 * balance becomes; it never touches the database.
 *
 * WEIGHTED_AVERAGE is implemented. FIFO needs cost layers (one row per
 * receipt with remaining quantity) and would implement the same interface,
 * consuming the oldest layers on issue.
 */
export interface Balance { quantity: string; value: string }

export interface CostingMethod {
  receive(balance: Balance, quantity: string, totalCost: string): Balance;
  issue(balance: Balance, quantity: string): { cost: string; balance: Balance };
}

export const weightedAverage: CostingMethod = {
  receive(b, quantity, totalCost) {
    return {
      quantity: new Decimal(b.quantity).plus(quantity).toString(),
      value: toMoney(new Decimal(b.value).plus(totalCost)),
    };
  },
  issue(b, quantity) {
    const qty = new Decimal(b.quantity);
    const out = new Decimal(quantity);
    // The last units out take exactly what is left, so no value is stranded.
    const cost = out.equals(qty)
      ? new Decimal(b.value)
      : new Decimal(b.value).times(out).dividedBy(qty).toDecimalPlaces(2, Decimal.ROUND_HALF_UP);
    return {
      cost: toMoney(cost),
      balance: { quantity: qty.minus(out).toString(), value: toMoney(new Decimal(b.value).minus(cost)) },
    };
  },
};

export const COSTING: Record<string, CostingMethod> = { WEIGHTED_AVERAGE: weightedAverage };
