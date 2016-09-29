## Local OPAM repository

We are vendoring OPAM metadata in this directory.

- `packages/local` contains the list of local packages, whose source code
  is contained in this repository.

- `packages/dev/ contains patched packages that we are using. The goal is to
  keep that list as small as possible.

- `packages/upstream` contains the metadata about upstream packages that we are
  using. That list contains all the dependencies needed to compile the local
  packages and optionally more packages which can be specified in the `PKGS`
  variable from the [Makefile](./Makefile). To update the list of upstream
  packages, simply run:

    ```
make
    ```
