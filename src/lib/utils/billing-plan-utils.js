import {
  each,
  filter,
  get,
  negate,
} from 'lodash';
import { DOWNLOAD_ATTRIBUTES } from 'lib/constants/eusa';


// Returns true if the given plan is free, false otherwise.
// A plan is considered to be free if:
// 1. It has no pricing components OR
// 2. The price of each pricing component is 0
export function isFreePlan(plan) {
  let free = true;

  const components = plan.pricing_components;
  each(components, (component) => {
    const tiers = component.tiers;
    each(tiers, (tier) => {
      if (tier.price > 0.00) {
        free = false;
      }
    });
  });
  return free;
}

// Returns true if the given plan is paid, false otherwise.
export const isPaidPlan = negate(isFreePlan);

/* Returns a string representing the price of a plan.
 * Assumptions:
 * - All subscriptions are monthly
 * - All store products have a single pricing component
 *  - This pricing component has a single tier
 *    - The lower and upper threshold of this tier are both 1
 * - All prices are in USD
 * TODO Kristie 8/17/16 Integrate i18n helpers for formatting numbers / currency
 */
export function getPriceStringForRatePlan(ratePlan) {
  const price = get(ratePlan, ['pricing_components', 0, 'tiers', 0, 'price']);
  if (!price) {
    return '$0.00';
  }
  return `$${price} / mo`;
}


/* Returns true if a product has ONLY anonymous download plans */
export function isAnonymousDownloadProduct(product) {
  if (!product) { return false; }
  // Find any plans that are not anonymous download
  const nonAnonPlans = filter(product.plans, ({ download_attribute }) => {
    return download_attribute !== DOWNLOAD_ATTRIBUTES.ANONYMOUS;
  });
  return !nonAnonPlans.length;
}
