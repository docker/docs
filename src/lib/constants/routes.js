import routeGenerator from 'lib/utils/route-generator';
/* eslint-disable max-len */

/* Routes */
export default routeGenerator({
  admin: '/admin',
  adminRepository: '/admin/repository/:id',
  beta: '/beta',
  browse: '/search?page_size=99&q=',
  bundleDetail: '/bundles/:id',
  bundleDetailPurchase: '/bundles/:id/purchase',
  communityImageDetail: '/community/images/:namespace/:reponame',
  communityImageDetailComments: '/community/images/:namespace/:reponame/comments',
  communityImageDetailTags: '/community/images/:namespace/:reponame/tags',
  home: '/',
  imageDetail: '/images/:id',
  login: '/login',
  publisher: '/publisher',
  publisherAddProduct: '/publisher/add-product',
  publisherSignup: '/publisher/signup',
  search: '/search',
});
