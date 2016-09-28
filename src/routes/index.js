/* eslint-disable max-len */
import React from 'react';
import { Route, IndexRoute, IndexRedirect } from 'react-router';

import Admin from 'components/admin';
import AdminHome from 'components/admin/AdminHome';
import AdminRepository from 'components/admin/AdminRepository';
import Beta from 'components/Beta';
import BundleDetail from 'components/marketplace/BundleDetail';
import CommunityImageDetail from 'components/marketplace/ImageDetail/CommunityImageDetail';
import DDCDetail from 'components/marketplace/BundleDetail/DDCDetail';
import Home from 'components/Home';
import ImageDetail from 'components/marketplace/ImageDetail';
import ImageDetailComments from 'components/marketplace/ImageDetail/ImageDetailComments';
import ImageDetailTags from 'components/marketplace/ImageDetail/ImageDetailTags';
import Login from 'components/Login';
import onEnterFn from './onEnter';
import PartnerImageDetail from 'components/marketplace/ImageDetail/PartnerImageDetail';
import Publisher from 'components/publisher';
import PublisherAddProduct from 'components/publisher/PublisherAddProduct';
import PublisherSignup from 'components/publisher/PublisherSignup';
import DDCPurchase from 'components/marketplace/DDCPurchase';
import Root from 'components/root';
import RouteNotFound404 from 'components/RouteNotFound404';
import Search from 'components/marketplace/Search';

export default function (store) {
  const onEnterFunctions = onEnterFn(store);
  const {
    betaPage,
    bundleDetailHome,
    confirmEmail,
    communityImageDetailComments,
    communityImageDetailHome,
    communityImageDetailTags,
    imageDetailHome,
    initializePurchase,
    landingPage,
    root,
    search,
    admin,
    adminHome,
    adminRepository,
    publisher,
    publisherSignup,
    publisherAddProduct,
  } = onEnterFunctions;

  // Function that allows you to reuse an onEnter function as an onChange func
  // since they take different parameters
  // https://github.com/reactjs/react-router/blob/master/docs/API.md
  const createOnChangeFn = (onEnterFunction) => (prevState, nextState, replace, callback) => {
    return onEnterFunction(nextState, replace, callback);
  };

  return (
    <Route path="/" component={Root} onEnter={root}>
      <IndexRoute component={Home} onEnter={landingPage} />

      <Route path="login" component={Login} />

      <Route path="beta" component={Beta} onEnter={betaPage} />

      <Route path="account/confirm-email/:confirmation_key" onEnter={confirmEmail} />

      {/* Certified Bundles */}
      <Route path="bundles/:id" component={BundleDetail}>
        <IndexRoute component={DDCDetail} onEnter={bundleDetailHome} />
        <Route path="purchase" onEnter={initializePurchase} component={DDCPurchase} />
      </Route>

      {/* Certified images */}
      <Route path="images/:id" component={ImageDetail}>
        <IndexRoute component={PartnerImageDetail} onEnter={imageDetailHome} />
      </Route>

      {/* Community images */}
      <Route path="community/images/:namespace/:reponame" component={ImageDetail}>
        <IndexRoute component={CommunityImageDetail} onEnter={communityImageDetailHome} />
        <Route path="comments" component={ImageDetailComments} onEnter={communityImageDetailComments} />
        <Route path="tags" component={ImageDetailTags} onEnter={communityImageDetailTags} />
      </Route>

      <Route path="search" component={Search} onEnter={search} onChange={createOnChangeFn(search)} />

      <Route path="publisher" component={Publisher} onEnter={publisher}>
        <Route path="signup" component={PublisherSignup} onEnter={publisherSignup} />
        <Route path="add-product" component={PublisherAddProduct} onEnter={publisherAddProduct} />
        <IndexRedirect to="add-product" />
      </Route>

      <Route path="admin" component={Admin} onEnter={admin}>
        <IndexRoute component={AdminHome} onEnter={adminHome} />
        <Route path="repository/:id" component={AdminRepository} onEnter={adminRepository} />
      </Route>

      {/* Route not found */}
      <Route path="*" component={RouteNotFound404} />
    </Route>
  );
}
/* eslint-enable */
