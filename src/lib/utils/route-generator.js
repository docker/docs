import UrlPattern from 'url-pattern';

/*
  Example Usage:

  // define the routes in one file
  import routeGenerator from 'lib/utils/route-generator';
  const routes = {
    list: 'some/path/with/:param',
    ...
  };
  export default routeGenerator(routes);

  // import them in another file
  import { list } from 'the/place/where/myRoutes';
  const pathToList = list({ param: 'myParam' });
  // pathToList is 'some/path/with/myParam'
*/
export default routes => {
  return Object.keys(routes).reduce((m, r) => {
    const p = new UrlPattern(routes[r]);
    // eslint-disable-next-line no-param-reassign
    m[r] = p.stringify.bind(p);
    return m;
  }, {});
};
