import { getLoginRedirectURL } from 'lib/utils/url-utils';
import routes from 'lib/constants/routes';

import {
  marketplaceFetchCategories,
  marketplaceFetchMostPopular,
  marketplaceFetchFeatured,
  marketplaceFetchPlatforms,
  marketplaceFetchRepositoryDetail,
  marketplaceFetchBundleDetail,
  marketplaceSearch,
} from 'actions/marketplace';

import {
  repositoryFetchComments,
  repositoryFetchOwnedNamespaces,
  repositoryFetchRepositoriesForNamespace,
  repositoryFetchImageDetail,
  repositoryFetchImageTags,
} from 'actions/repository';

import {
  accountFetchCurrentUser,
  accountFetchUserEmails,
  accountToggleMagicCarpet,
  accountSelectNamespace,
} from 'actions/account';

import {
  billingFetchPaymentMethods,
  billingFetchProfile,
  billingFetchProduct,
  billingFetchProfileSubscriptions,
} from 'actions/billing';

import {
  whitelistFetchAuthorization,
  whitelistAmIWaiting,
} from 'actions/whitelist';

import {
  publishFetchProductDetails,
  publishFetchProductList,
  publishFetchProductTiers,
  publishGetPublishers,
  publishGetSignup,
  publishGetVendorAgreement,
} from 'actions/publish';

import {
  finishPageTransition,
  startPageTransition,
} from 'actions/root';

import { DDC_ID } from 'lib/constants/eusa';
import {
  getProductStep,
} from 'lib/utils/publisherSteps';
import { post } from 'superagent';
import { readCookie } from 'lib/utils/cookie-handler';
import get from 'lodash/get';
import isEmpty from 'lodash/isEmpty';
import 'lib/utils/promise';
import NProgress from 'nprogress';
import 'lib/css/nprogress.css';

export default function (store) {
  const { dispatch, getState } = store;

  const createOnEnterFn = (fn) => (nextState, replace, callback) => {
    const isPromise = (promise) => {
      return promise !== null &&
        typeof promise === 'object' &&
        promise &&
        typeof promise.then === 'function';
    };

    dispatch(startPageTransition());
    NProgress.configure({ showSpinner: false });
    NProgress.start();
    NProgress.set(0.5);

    const promise = fn(nextState, replace);
    if (isPromise(promise)) {
      promise.then(() => {
        NProgress.done();
        callback();
        dispatch(finishPageTransition());
      }).catch(() => {
        NProgress.done();
        callback();
        dispatch(finishPageTransition());
      });
    } else {
      callback();
      dispatch(finishPageTransition());
    }
  };

  const requireLogin = (fn) => (nextState, replace) => {
    const state = store.getState();
    if (!state.account ||
        !state.account.currentUser ||
        !state.account.currentUser.id) {
      replace(getLoginRedirectURL());
      return Promise.resolve();
    }
    return fn(nextState, replace);
  };

  const requireWhitelist = (fn) => (nextState, replace) => {
    const state = store.getState();
    if (!state.account || !state.account.isCurrentUserWhitelisted) {
      replace(routes.beta());
      return Promise.resolve();
    }
    return fn(nextState, replace);
  };

  const requireApprovedPublisher = (fn) => (nextState, replace) => {
    const state = store.getState();
    const signup = get(state, ['publish', 'signup', 'results', 'status']);
    const publisher =
      get(state, ['publish', 'publishers', 'results']);
    if (signup !== 'reviewed' || isEmpty(publisher)) {
      replace(routes.publisherSignup());
      return Promise.resolve();
    }
    return fn(nextState, replace);
  };

  //----------------------------------------------------------------------------
  // ROOT
  //----------------------------------------------------------------------------
  const root = ({ location }, replace) => {
    const { query } = location;
    return dispatch(accountFetchCurrentUser()).then(({ value }) => {
      const user = value;
      dispatch(accountSelectNamespace({ namespace: user.username }));
      const promises = [
        dispatch(whitelistFetchAuthorization()),
        dispatch(publishGetSignup()),
        dispatch(accountFetchUserEmails({ user: user.username })).then(res => {
          const emails = res.value;
          const primaryEmails = emails &&
            emails.results &&
            emails.results.filter(r => r.primary);
          const primaryEmail = primaryEmails &&
            primaryEmails.length &&
            primaryEmails[0].email;

          if (user && primaryEmail) {
            analytics.alias(user.id, () => {
              analytics.identify(user.id, {
                Docker_Hub_User_Name__c: user.username,
                username: user.username,
                dockerUUID: user.id,
                point_of_entry: 'docker_store',
                email: primaryEmail,
              });
            });
          }
        }),
      ];

      if (query.overlay) {
        dispatch(accountToggleMagicCarpet({
          magicCarpet: query.overlay,
        }));
      }

      return Promise.when(promises);
    }).catch(err => {
      const state = store.getState();
      const loggedOut = !state.account ||
        !state.account.currentUser ||
        !state.account.currentUser.id;
      if (query.overlay && loggedOut) {
        replace(getLoginRedirectURL());
      }
      return Promise.reject(err);
    });
  };

  //----------------------------------------------------------------------------
  // BETA PAGE
  //----------------------------------------------------------------------------
  const betaPage = (nextState, replace) => {
    const state = store.getState();
    if (state.account && state.account.isCurrentUserWhitelisted) {
      replace(routes.home());
      return Promise.resolve();
    }

    return dispatch(whitelistAmIWaiting());
  };

  //----------------------------------------------------------------------------
  // HOME / LANDING PAGE
  //----------------------------------------------------------------------------
  const landingPage = () => {
    return Promise.when([
      dispatch(marketplaceFetchCategories()),
      dispatch(marketplaceFetchFeatured()),
      dispatch(marketplaceFetchMostPopular()),
    ]);
  };

  //----------------------------------------------------------------------------
  // SEARCH
  //----------------------------------------------------------------------------
  const search = ({ location }) => {
    return Promise.when([
      dispatch(marketplaceFetchCategories()),
      dispatch(marketplaceFetchPlatforms()),
      dispatch(marketplaceSearch(location && location.query || {})),
    ]);
  };

  //----------------------------------------------------------------------------
  // MARKETPLACE / CERTIFIED IMAGE
  //----------------------------------------------------------------------------
  const imageDetailHome = ({ params }) => {
    const { id } = params;
    const promises = [
      dispatch(marketplaceFetchRepositoryDetail({ id })),
      dispatch(billingFetchProduct({ id })),
    ];
    const { account } = getState();
    // If the user is logged in, their information will be in currentUser
    const docker_id = account && account.currentUser && account.currentUser.id;
    if (docker_id) {
      promises.push(dispatch(billingFetchProfile({ docker_id })));
      promises.push(dispatch(billingFetchProfileSubscriptions({ docker_id })));
    }
    return Promise.when(promises);
  };

  //----------------------------------------------------------------------------
  // COMMUNITY IMAGE
  //----------------------------------------------------------------------------
  const communityImageDetailHome = ({ params }) => {
    const { namespace, reponame } = params;
    // Must request 'all' tags to catch default tag and show scan on home page
    const page_size = 100;
    return Promise.when([
      dispatch(repositoryFetchImageDetail({ namespace, reponame })),
      dispatch(repositoryFetchComments({ namespace, reponame })),
      dispatch(repositoryFetchImageTags({ namespace, page_size, reponame })),
    ]);
  };

  const communityImageDetailComments = ({ params, location }) => {
    const { namespace, reponame } = params;
    const { page, page_size } = location && location.query;
    return dispatch(repositoryFetchComments({
      namespace,
      reponame,
      page,
      page_size,
    }));
  };

  const communityImageDetailTags = ({ params, location }) => {
    const { namespace, reponame } = params;
    const { page, page_size } = location && location.query;
    // TODO Kristie 6/8/16 Remove image details call when marketplace tags api
    // is available
    return Promise.when([
      dispatch(repositoryFetchImageDetail({ namespace, reponame })),
      dispatch(repositoryFetchImageTags({
        namespace, reponame, page, page_size,
      })),
    ]);
  };

  //----------------------------------------------------------------------------
  // BUNDLES
  //----------------------------------------------------------------------------
  const bundleDetailHome = ({ params }) => {
    const { id } = params;
    dispatch(marketplaceFetchBundleDetail({ id }));
    dispatch(billingFetchProduct({ id }));

    const { account } = getState();
    // If the user is logged in, their information will be in currentUser
    const docker_id = account && account.currentUser && account.currentUser.id;
    if (docker_id) {
      dispatch(billingFetchProfileSubscriptions({
        docker_id,
      }));
    }
    return Promise.resolve();
  };

  //----------------------------------------------------------------------------
  // PURCHASE PAGES
  //----------------------------------------------------------------------------
  const initializePurchase = ({ params }) => {
    const { id } = params;
    return Promise.when([
      // TODO: nathan - can remove some of these fetches since it's in the root
      dispatch(accountFetchCurrentUser())
        .then((res) => {
          const { username: namespace, id: docker_id } = res.value;
          return Promise.all([
            dispatch(accountFetchUserEmails({ user: namespace })),
            dispatch(accountSelectNamespace({ namespace })),
            dispatch(billingFetchProfile({ docker_id })).then((billingRes) => {
              let promises = [];
              if (billingRes.value && billingRes.value.profile) {
                // If billing profile exists, fetch user subscriptions
                promises = [
                  dispatch(billingFetchProfileSubscriptions({ docker_id })),
                  dispatch(billingFetchPaymentMethods({ docker_id })),
                ];
              }
              return Promise.when(promises);
            }),
          ]);
        }),
      dispatch(repositoryFetchOwnedNamespaces()),
      // Get product billing information from billing API
      dispatch(billingFetchProduct({ id })),
      // Get product description
      dispatch(marketplaceFetchBundleDetail({ id })),
    ]);
  };


  //----------------------------------------------------------------------------
  // PUBLISHER
  //----------------------------------------------------------------------------
  const publisher = () => {
    return Promise.when([
      dispatch(publishGetPublishers()),
      dispatch(publishGetSignup()),
    ]);
  };

  const publisherSignup = (nextState, replace) => {
    const state = store.getState();
    if (
      state.publish.signup &&
      state.publish.signup.results &&
      state.publish.signup.results.status === 'reviewed'
    ) {
      replace(routes.publisherAddProduct());
    }
    return Promise.resolve();
  };

  const publisherAddProduct = () => {
    return Promise.when([
      dispatch(publishGetVendorAgreement()),
      dispatch(marketplaceFetchCategories()),
      dispatch(repositoryFetchOwnedNamespaces()),
      dispatch(publishFetchProductList()).then((listRes) => {
        /*
        TODO: nathan 7/28/16
        - Currently only handling a single product.
        - Update routes with root for dashboard & param for product id
        1) FETCH PRODUCT LIST
        2) Pull relevant product from list (via id) (for now only 1 product)
        3) Fetch product details given product_id
          3a) update reducer to handle product details & pull repositories
          correctly
        */
        // NOTE: If No Publisher Data is found - 204 No Content.
        if (listRes.value.length > 0) {
          const { id: product_id, status } = listRes.value[0];
          const productStep = getProductStep(status);
          dispatch(publishFetchProductDetails({ product_id }))
            .then((detailRes) => {
              const product = detailRes.value;
              const promises = [];
              if (product.repositories && product.repositories.length > 0) {
                product.repositories.forEach((repo) => {
                  promises.push(
                    dispatch(repositoryFetchRepositoriesForNamespace({
                      namespace: repo.namespace,
                    }))
                  );
                  promises.push(
                    dispatch(repositoryFetchImageTags({
                      namespace: repo.namespace,
                      reponame: repo.reponame,
                      page_size: 0,
                    }))
                  );
                });
              }
              return Promise.all(promises);
            }).then(() => {
              if (productStep > 1) {
                dispatch(publishFetchProductTiers({ product_id }));
              }
            });
        }
      }),
    ]);
  };

  //----------------------------------------------------------------------------
  // ADMIN
  //----------------------------------------------------------------------------
  const admin = ({ location }, replace) => {
    const state = store.getState();
    if (!state.account ||
        !state.account.currentUser ||
        !state.account.currentUser.is_staff) {
      replace({ pathname: routes.login(), query: { next: location.pathname } });
    }
  };

  const adminHome = () => {
    return dispatch(marketplaceSearch({ page_size: 9999 }));
  };

  const adminRepository = ({ params }) => {
    const { id } = params;
    return Promise.when([
      dispatch(marketplaceFetchCategories()),
      dispatch(marketplaceFetchPlatforms()),
      dispatch(marketplaceFetchRepositoryDetail({ id })),
    ]);
  };

  const confirmEmail = ({ params, location }, replace) => {
    const { confirmation_key } = params;
    const { ref } = location && location.query;

    const req = post('/v2/users/activate/')
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json')
      .set('X-CSRFToken', readCookie('csrftoken'));

    if (ref) {
      analytics.track('ddc_verify_account', {
        plan: ref,
      });
      replace({
        pathname: routes.bundleDetailPurchase({ id: DDC_ID }),
        query: { plan: ref, verified: true },
      });
    } else {
      replace({ pathname: routes.login() });
    }

    req.send({ confirmation_key }).end((err) => {
      console.log(err);
    });
  };

  /* eslint-disable max-len */
  return {
    admin: createOnEnterFn(admin),
    adminHome: createOnEnterFn(adminHome),
    adminRepository: createOnEnterFn(adminRepository),
    betaPage: createOnEnterFn(betaPage),
    bundleDetailHome: createOnEnterFn(bundleDetailHome),
    confirmEmail: createOnEnterFn(confirmEmail),
    communityImageDetailComments: createOnEnterFn(requireWhitelist(communityImageDetailComments)),
    communityImageDetailHome: createOnEnterFn(requireWhitelist(communityImageDetailHome)),
    communityImageDetailTags: createOnEnterFn(requireWhitelist(communityImageDetailTags)),
    imageDetailHome: createOnEnterFn(requireWhitelist(imageDetailHome)),
    initializePurchase: createOnEnterFn(initializePurchase),
    landingPage: createOnEnterFn(requireWhitelist(landingPage)),
    publisher: createOnEnterFn(requireLogin(publisher)),
    publisherSignup: createOnEnterFn(requireLogin(publisherSignup)),
    publisherAddProduct: createOnEnterFn(requireApprovedPublisher(requireLogin(publisherAddProduct))),
    root: createOnEnterFn(root),
    search: createOnEnterFn(requireWhitelist(search)),
  };
  /* eslint-enable */
}
