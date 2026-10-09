# Express upgrade investigation

This historical Express 4.21.2 application demonstrates a routing change in
Express 5.1.0. These versions reproduce the migration problem. They are not
recommendations for a production application.

Use Node.js 22 or later. From this directory, install the locked dependencies
and establish the baseline:

```console
$ npm ci
$ npm test
```

The three tests check `/`, `/projects/42/settings`, and `/api/status`. The
fallback route serves the client application after the API route.

Investigate the upgrade in a disposable copy:

```console
$ npm install --save-exact express@5.1.0
$ npm test
```

The application fails to load because the wildcard needs a name. Changing
`*` to `/*splat` leaves `/` unmatched. Changing it to `/{*splat}` includes `/`
and restores all three tests. See the
[Express migration guide](https://expressjs.com/en/guide/migrating-5.html#path-route-matching-syntax).

The repository-level verifier exercises all four stages and a deliberately
failing baseline. It keeps this directory at Express 4.21.2 for the agent to
investigate.
