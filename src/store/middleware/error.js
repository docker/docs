import { getLoginRedirectURL } from 'lib/utils/url-utils';
import { redirectTo } from 'actions/root';
import get from 'lodash/get';
import routes from 'lib/constants/routes';

export default ({ dispatch }) => next => action => {
  const { payload, meta, error } = action;
  if (!payload || !error) return next(action);

  // In order to use the redirect capabilities, `shouldRedirectToLogin: true`
  // should be set in the `meta` field of the action itself.
  // If an action with that property set receives a 401 Unauthorized, it will
  // redirect you to login
  const isUnauthorized = payload.response && payload.response.unauthorized;
  if (!isUnauthorized) return next(action);
  // If the JWT is invalid (possibly due to certificate rotation), redirect to
  // login page so that they get new credentials.
  // TODO Don't use string comparison to match invalid error - use code when
  // provided by gateway / backend: tracked in MER-691
  const isJWTInvalid = get(payload, ['response', 'body', 'detail']) ===
    'Invalid JWT, please login again';
  const isLoginPage = get(window, ['location', 'pathname']) === routes.login();
  // Ignore the login page to prevent infinite loop of redirects to login page
  const shouldForceLogin = isJWTInvalid && !isLoginPage;
  const shouldRedirectToLogin = meta && meta.shouldRedirectToLogin;
  if (shouldRedirectToLogin || shouldForceLogin) {
    return dispatch(redirectTo(getLoginRedirectURL()));
  }
  return next(action);
};
