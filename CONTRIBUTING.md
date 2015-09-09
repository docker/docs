Contributing to  Docker Toolbox
==================================

The Docker Toolbox is a part of the [Docker](https://www.docker.com) project, and follows
the [contributing guidelines](https://github.com/docker/docker/blob/master/CONTRIBUTING.md). If you're already familiar with the way
Docker does things, you'll feel right at home.

Please [sign your work](https://github.com/docker/docker/blob/master/CONTRIBUTING.md#sign-your-work) before creating a pull request.

Thanks for taking the time to improve the Docker Toolbox!

## License

By contributing your code, you agree to license your contribution under the [Apache license](https://github.com/docker/toolbox/blob/master/LICENSE/LICENSE).

## Diff scpt files

.gitattributes
```
*.scpt diff=scpt
```

.git/config
```
[diff "scpt"]
    textconv = osadecompile
    binary = true
```
