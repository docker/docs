# Linting

All code must pass [ESLint][eslint] before being merged into master
branch. The ESLint config can be found in `.eslintrc` and is
integrated into webpack.

# Running ESLint

```
gulp webpack
```

Since linting is integrated with webpack, it is possible to lint code
while it is being developed without any extra effort. This is
important because if it is not approximate to effortless to run
linting, it will not be run while developing.

# We should block deploys for linting errors

Since we do CI/CD, the static analysis present in ESLint can help us
catch bugs before shipping. We should therefore block deploys if
ESLint detects an error-level (level `2` in `.eslintrc`) issue.

* [eslint][eslint]
* [babel-eslint][babel-eslint]

[eslint]: http://eslint.org/
[babel-eslint]: https://github.com/babel/babel-eslint
