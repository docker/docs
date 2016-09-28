# ng-annotate-loader [![Build Status](https://img.shields.io/travis/huston007/ng-annotate-loader.svg?style=flat-square)](https://travis-ci.org/huston007/ng-annotate-loader)
Webpack loader to annotate angular applications. Generates a sourcemaps as well.

## Usage:

```
module: {
    loaders: [
      {test: /src.*\.js$/, loaders: ['ng-annotate']},
    ]
  }
```

#### Passing parameters: 

```
	{test: /src.*\.js$/, loaders: ['ng-annotate?add=false&map=false']}
```

[More about `ng-annotate` parameters](https://github.com/olov/ng-annotate#library-api)

#### Using ng-annotate plugins: 

```
	{test: /src.*\.js$/, loaders: ['ng-annotate?plugin[]=ng-annotate-adf-plugin']}
```

#### Works great with js compilers, `babel` for example:

```
    {test: /src.*\.js$/, loaders: ['ng-annotate', 'babel-loader']},
```

## Contributing
#### Compiling examples and run acceptance test
Run on the root folder:
```
npm install
npm test
```

[Using loaders](http://webpack.github.io/docs/using-loaders.html)
