### Working within DTR

1. Make via: `BIND_ASSETS=1 make` in `dhe-deploy`.  This tells stacker to bind the current UI
   folder to the assets location within the container, giving you instant
   feedback without reloading containers. This will still build the bundle with the prod config.

2. Install dependencies via: `npm install` (only need to do once or whenever packages change).

3. Compile and watch via: `npm run-script watch`. This rebuilds the bundle with dev config.

### NPM commands:

`npm install` : Installs dependencies

`npm run-script build`: Starts webpack compilation to src/bundle.js using
the production config

`npm run-script watch`: Starts webpack compilation to src/bundle.js and 
*watches* for file changes. This uses `webpack.config.dev.js`

TODO: React Hot Reloader
