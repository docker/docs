# Quickstart

*Make sure to clone this repo into your `/Users/<username>` directory for it to run correctly*

```bash
make dns
make hub-deps
# you must log in as the 'dux' user. ask one of the frontend 
# team members for credentials
npm login
npm install
docker-compose build
npm run build:dev
./startup-scripts/bootstrap-dev.sh
docker-compose up -d
```

At this point you will need `tmux` to run `boot-dev-tmux.sh`, it can
be installed on OSX by `brew install tmux`

```bash
./startup-scripts/boot-dev-tmux.sh
```

## tmux env

Here are some basic commands to help you get around tmux. `C` is
Control, `-` means hit both keys, everything else it a literal
character you need to produce.

| Command      | Keys    |
|--------------|---------|
| Next Window  | C-b n   |
| Next Panel   | C-b o   |
| Close Window | C-b & y |

# Docs

* [React](docs/concepts/React.md)
* [Flux](docs/concepts/Flux.md)
* [React Native](docs/concepts/React-Native.md)
* [React Router](docs/concepts/React-Router.md)
* [Immutability](docs/concepts/Immutability.md)
