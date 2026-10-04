import { describe, expect, it } from 'vitest';
import { formatAmount, fromHalalas, toHalalas } from './money';

describe('money display helpers', () => {
  it('formats server amounts without floating point', () => {
    expect(formatAmount('1150.00')).toBe('1,150.00');
    expect(formatAmount('-1234567.5')).toBe('-1,234,567.50');
    expect(formatAmount('0.1')).toBe('0.10');
    expect(formatAmount('999999999999999.99')).toBe('999,999,999,999,999.99');
    expect(formatAmount(null)).toBe('—');
  });

  it('parses typed amounts into halalas and back', () => {
    expect(toHalalas('0.10')).toBe(10);
    expect(toHalalas('0.1')! + toHalalas('0.2')!).toBe(30); // no 0.30000000000000004
    expect(toHalalas('12.345')).toBeNull();
    expect(toHalalas('abc')).toBeNull();
    expect(toHalalas('')).toBe(0);
    expect(fromHalalas(-5)).toBe('-0.05');
    expect(fromHalalas(115000)).toBe('1150.00');
  });
});
