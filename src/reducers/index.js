import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';

import account, {
  DEFAULT_STATE as ACCOUNT_DEFAULT_STATE,
} from './account';
import billing, {
  DEFAULT_STATE as BILLING_DEFAULT_STATE,
} from './billing';
import internal, {
  DEFAULT_STATE as INTERNAL_DEFAULT_STATE,
} from './internal';
import marketplace, {
  DEFAULT_STATE as MARKETPLACE_DEFAULT_STATE,
} from './marketplace';
import publish, {
  DEFAULT_STATE as PUBLISHER_DEFAULT_STATE,
} from './publish';
import root, {
  DEFAULT_STATE as ROOT_DEFAULT_STATE,
} from './root';

export const DEFAULT_STATE = {
  account: ACCOUNT_DEFAULT_STATE,
  billing: BILLING_DEFAULT_STATE,
  internal: INTERNAL_DEFAULT_STATE,
  marketplace: MARKETPLACE_DEFAULT_STATE,
  publish: PUBLISHER_DEFAULT_STATE,
  root: ROOT_DEFAULT_STATE,
};

export default combineReducers({
  account,
  billing,
  form: formReducer,
  internal,
  marketplace,
  publish,
  root,
});
