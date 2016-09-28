# Docker Universal Control Plane UI

The front-end is a static AngularJS application that is served by UCP Controller.
It uses Semantic-UI, along with our custom theme to achieve the presentation layer.

## Directory Structure

```
├── node modules
├── dist                                # Output directory served by controller
│   ├── 3cf93a....svg
│   ├── 570764....svg
│   ├── 706450....ttf
│   ├── bundle.js
│   └── index.html
├── package.json                        # npm dependency manifest
├── README.md
├── semantic                            # Semantic-UI source tree
│   ├── dist
│   ├── gulpfile.js
│   ├── package.json
│   ├── semantic.json
│   ├── src                             # Custom theme is in here
│   └── tasks
├── src
│   ├── app                             # AngularJS code
│   ├── index.html                      # Simple static index.html that just includes bundle.js
│   ├── styles                          # Site-wide styles
│   └── test                            # Unit tests
├── webpack.config.builder.js           # Webpack config builder used by dev/release configs
├── webpack.config.dev.js               # Dev mode bundler config
└── webpack.config.release.js           # Release mode bundler config
```

## Building the UI

Pre-requisites for building are nodejs 4.x.

The UI can be built using the command below, from the root of the repo:

```
make media
```

## Webpack Bundler

We use webpack as the bundler for the UI, given an entrypoint js file, this will
crawl the app looking for `require(...)` modules and resolve the module tree
that we need.  It will then concatenate the files together, including non-JS
assets and modules were necessary (e.g. fonts, images, html).  The output from
webpack will be a `dist/bundle.js` that we can include into our index.html.

In the root of `controller/static`, there are three webpack related config
files.  `webpack.config.builder.js` is a config builder, which containers the
main webpack config, along with conditionals based upon parameters that are
passed in to the builder, to switch plugins/modules on and off.

Examples of this being used are the `webpack.config.release.js` and
`webpack.config.dev.js`.  The release config has minification enabled, ready for
release.  Whereas the dev config does not, and has debug and devtools enabled.

### Using Dev mode

To use the the *dev* config, the easiest way at the moment is to:

#### Mac

Steps:

1) Run the bootstrapper

```
$ docker run --rm -v -i /var/run/docker.sock:/var/run/docker.sock \
  --name ucp docker/ucp install --swarm-port=1337
```

2) Create .env file

    a) If you are using Docker for Mac:

    $ echo -e "CONTROLLER_URL=https://$(docker-machine ip default):443" > .env

    b) If you are using VirtualBox, inotify will not work and therefore we have to switch to using polling:

    $ echo -e "CONTROLLER_URL=https://$(docker-machine ip default):443\nPOLL=true" > .env

3) Run compose

```
$ docker-compose -f docker-compose.ui.yml up
```

4) Open UCP in browser

```
$ open $(echo "http://"$(docker-machine ip default)":9090")
```

#### Linux

```
$ docker-compose -f docker-compose.yml up
```

or against a bootstrapper installed instance:

```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/    ucp install
$ echo "CONTROLLER_URL=https://$(ip -4 -o addr show dev docker0 | cut -f7 -d' '|     cut -f1 -d/):443" > .env
$ docker-compose -f docker-compose.ui.yml up
```

The `docker-compose.ui.yml` will launch two containers:

- A `webpack-dev-server` with the `--watch` option, accessible on port `9090`.  This will re-bundle the modules
  whenever a change is made, and then notify your browser via a websocket to reload the page.
- A `node` container that runs `gulp --watch` on our semantic theme, this will live re-build the semantic custom
  theme.  The built semantic theme will then trigger webpack-dev-server to re-bundle.

While `webpack-dev-server` can serve static content, we are only using it in proxy mode.  All
requests will be proxied through to the URL configured in the `CONTROLLER_URL` environment
variable in the `.env` environment file.

Our Makefile will only build a bundle from the `webpack.config.release.js`.

### Dev Mode Limitations

There are a few features of UCP that will not work using the `webpack-dev-server` proxy due to not being able to
proxy WebSockets at this time ([https://github.com/webpack/webpack-dev-server/issues/283](https://github.com/webpack/webpack-dev-server/issues/283)).
There are some solutions and workarounds that still need to be investigated.  The current status is that WebSocket
proxying will be supported by default in the next major release of webpack-dev-server.

## CommonJS Modules

Using modules allows us to explicitly include javascript source, html and other assets into our
application using `require(...)`.

In order for a javascript source file to be treated as a module, it must export a function.
Taking a controller as an example, we define a function `ContainersController` as normal, but
instead of it loading itself into the Angular context, we allow the source to export the
function as a module via `module.exports`.

```javascript
'use strict';

ContainersController.$inject = ['...', '...'];
function ContainersController(..., ...) {
    ...
}

module.exports = ContainersController;
```

We can then define an Angular module as normal, but use `require(...)` to include the functions
defined in each of the modules.  The easiest thing to do is to define an `index.js` in the module's
directory.  This will then explicitly `require` all of the other files in the directory, loading
them into the angular context, under a specific module name.

```javascript
'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.containers', [])
       .controller('ContainersController', require('./containers.controller'))
       .factory('ContainersService', require('./containers.service'))
       .config(require('./routes'));
```

Angular modules are then included into the root module via the `entrypoint.js`, like so:

```javascript
'use strict';

require('./vendor')();

angular
  .module('ducp', [
    require('./components/containers').name,
    ...
    ]
  )
```

Webpack will use this `entrypoint.js`, and crawl the `require(...)` calls to include all the
required sources.

Assets such at html and images also require bundling, in the case of defining a route with a
HTML template, a HTML file can be `require`d:

```javascript
'use strict';

var settingsTemplate = require('./settings.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard.settings', {
      ...
      template: settingsTemplate,
      controller: 'SettingsController',
      ...
    });
}

module.exports = getRoutes;
```

## Dependency Management

Dependencies are controlled by `npm`, the resulting `node_modules` directory
from running an `npm install` is vendored within the repository.

For a while we have been running `npm install` on each build to populate our runtime
and dev dependencies.  Due to the main npm repository being intermittently unavailable,
and this either making our media build extremely slow or non-functional, we decided to
vendor the dependencies.  This should help make our CI builds more reliable.

### Adding or Updating a Dependency

To add or update a runtime js dependency:

```
$ npm install --save <node_module_name>
$ git add node_modules
```

To add or update a build-time node dependency:

```
$ npm install --save-dev <node_module_name>
$ git add node_modules
```

## Semantic-UI

Semantic theme updates are infrequent, so the result of the semantic build is
vendored in `controller/static/semantic/dist`.  The Semantic build is very slow
and requires quite a few additional dev dependencies such as gulp.

### Building an Updated Theme

Building the theme requires gulp to be installed, that can be installed like so:

```
$ npm install -g gulp
```

The theme can then be built:

```
$ cd controller/static/semantic
$ npm install
$ gulp build
```

### Theme Customization

Before modifying the theme, it's worth first familiarising yourself with the
[Semantic-UI theme documentation](http://semantic-ui.com/usage/theming.html).
You'll see that there are three levels of inheritance (`default`, `packaged`
themes and `site` theme) in how styles and themes are applied and created.  
Themes make heavy use of LESS in order to save certain values and styles as
a variable and apply them in multiple places.

We should avoid modifying the `default` theme.

The `site` theme should be modified with any changes that we need to make.
Each component contains:

 - an `.overrides` file that can be used to supply additional CSS to modify
   how a component should look.
 - a `.variables` file that can be used to override LESS variables

### Font Vendoring

By default, Semantic-UI attempts to retrieve fonts remotely from Google Fonts.
To switch this off, `@importGoogleFonts: false;` has been added to the `site`
theme's site.variables.  This will disable Semantic-UI's remote font fetching.
Even though this has been disabled, some themes (e.g. `material`) attempt to
explicitly import a font directly from Google Fonts, in this case these imports
have been manually removed.

To vendor a font from Google Fonts, the first step is to open the Google Fonts
link to the font in your browser, and retrieve the relevant Latin CSS.  This
can then be placed in `semantic/src/site/globals/site.overrides`.

You'll notice that the CSS contains a remote URL, the next step is to download
the `ttf` formatted font from this URL and save it to
`semantic/src/site/assets/fonts/`.  Once you've saved the font, the `url(..)`
in the CSS should be updated to reflect this newly saved file.  To help with
locating this font, the `@{siteFontPath}` variable can be used as a base for
the path.

### ESLint

We use an eslint config (taken originally from DUX).  ESLint is executed automatically
during bundling.  Linting errors when performing a build will cause the build to fail
and errors will be displayed.

When running in dev mode, a bundle will still be produced but the errors will be displayed.
