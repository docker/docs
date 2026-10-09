# Express upgrade investigation

This historical Express 4.21.2 application represents a support portal with
client routes and a ticket-search API. The migration target is Express 5.1.0.
These versions reproduce compatibility changes. They are not
recommendations for a production application.

Use Node.js 22 or later. From this directory, install the locked dependencies
and establish the baseline:

```console
$ npm ci
$ npm test
```

Eight tests cover the client root and nested route, API status, ticket listing,
nested and combined query filters, JSON search, and an empty search request.
The fallback route serves the client application after the API routes.

Investigate the upgrade in a disposable copy:

```console
$ npm install --save-exact express@5.1.0
$ npm test
```

Assess whether the upgrade preserves the existing client and API behavior.
Record each failure, test proposed compatibility changes, and produce a concise
migration recommendation. See the
[Express migration guide](https://expressjs.com/en/guide/migrating-5.html).

The repository-level verifier exercises the upgrade, proposed fixes, and a deliberately
failing baseline. It keeps this directory at Express 4.21.2 for the agent to
investigate.
