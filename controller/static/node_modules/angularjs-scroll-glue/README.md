# angular-scroll-glue [![Build Status](https://travis-ci.org/Luegg/angularjs-scroll-glue.svg?branch=master)](https://travis-ci.org/Luegg/angularjs-scroll-glue)

An AngularJs directive that automatically scrolls to the bottom of an element on changes in its scope.

## Install
### Bower
```bash
$ bower install angular-scroll-glue --save
```

### npm
```bash
$ npm i angularjs-scroll-glue
```

## Usage
```javascript
// Add `luegg.directives` to your module's dependencies.
angular.module('yourModule', [
	...,
	'luegg.directives'
]);
```

```html
<div scroll-glue>
	<!-- Content here will be "scroll-glued". -->
</div>

<div scroll-glue="glued">
	<!-- Content here will be "scroll-glued" if the passed expression is truthy. -->
</div>
```

**More information can be found in the [live demo](#live-demo).**

## Live demo
[Demo Plunk](http://plnkr.co/edit/wxTyp7PpyxJOHSlUumVC?p=preview)

## Contribute

Despite this is a ultra specialized library, there will always be a fix to be made or a new feature to be added and I'm glad for any contributions. But please make sure to check the following before committing a pull request:

1. Make sure the unit tests pass. Just run `npm test` and check if all is green.
1. Try to add new tests that cover your changes.
1. Make sure you do not introduce changes that break backward compatibility unless there is a really good reason to.

## License (MIT)

Copyright (C) 2013 Luegg

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
