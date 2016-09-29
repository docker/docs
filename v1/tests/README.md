## Pinata-RT

Initial test-bed for the Pinata project.

### Build

You need to have `opam` installed (usually with `brew install opam`), and
then initialise it with
```
opam init # only the first time you install opam
```

```
make depends # only the first time
make
```

### Run

```
./pinata-rt`
```

This will download the latest release of `Docker.app` on
[Hockey App](http://hockeyapp.net/) and run simple end-to-end
tests (including performance tests).