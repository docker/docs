import { Decimal } from 'decimal.js';
import { z } from 'zod';

/**
 * Money handling. Amounts travel as strings ("1150.00") in the API and come
 * back from PostgreSQL NUMERIC as strings. Arithmetic uses decimal.js; a JS
 * number is never used for money.
 */
Decimal.set({ precision: 40, rounding: Decimal.ROUND_HALF_UP });

export { Decimal };

/** Positive or zero amount with at most 2 decimal places, as a string. */
export const amountString = z
  .string()
  .trim()
  .regex(/^\d{1,16}(\.\d{1,2})?$/, 'Amount must be a non-negative number with at most 2 decimal places, sent as a string');

export const toMoney = (v: Decimal.Value) => new Decimal(v).toFixed(2);
export const sum = (values: Decimal.Value[]) => values.reduce<Decimal>((acc, v) => acc.plus(v), new Decimal(0));
