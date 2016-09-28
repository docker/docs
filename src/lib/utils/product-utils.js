import { DDC_ID, CS_ENGINE_ID } from 'lib/constants/eusa';

import {
  includes,
} from 'lodash';

const DOCKER_PRODUCT_IDS = [
  DDC_ID,
  CS_ENGINE_ID,
];

const PRODUCT_BUNDLES = [
  DDC_ID,
  // TODO(mattt) 8/10/2016 There is a hack in place to surface the
  // CS Engine as a product bundle. This code will have to change
  // once that hack is removed.
  CS_ENGINE_ID,
];

// returns true if the given product id
// is for a Docker product, false otherwise.
export function isDockerProduct(product_id) {
  // TODO In the future, this should take into account
  // the product publisher_id, not a set of hard coded product ids
  return includes(DOCKER_PRODUCT_IDS, product_id);
}

// returns true if the product with the given id
// is a bundle, false otherwise.
export function isProductBundle(product_id) {
  // TODO In the future, the fact that a product is a
  // bundle should come from product catalog metadata
  return includes(PRODUCT_BUNDLES, product_id);
}

// returns true if the given is pullable, false otherwise.
export function isPullableProduct(product_id) {
  // TODO in the future, this should be informed by the product type.
  // For example, repository products are pullable, both others may not be.
  return product_id !== DDC_ID &&
         product_id !== CS_ENGINE_ID;
}
