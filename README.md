# Mercury
[![Circle CI](https://circleci.com/gh/docker/mercury-ui.svg?style=shield&circle-token=03609595272ae7f08f9b7d0d276f2dff340d8735)](https://circleci.com/gh/docker/mercury-ui) [![codecov.io](https://codecov.io/github/docker/mercury-ui/coverage.svg?branch=master&token=y0dAFR3ipK)](https://codecov.io/github/docker/mercury-ui?branch=master)


### Develop

- Download [Node](https://nodejs.org/en/) or install it via homebrew `brew install node`
- Install [React dev tools](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi) and [Redux dev tools](https://chrome.google.com/webstore/detail/redux-devtools/lmhkpmbekcpmknklioeibfkpmmfibljd)

In order to test user actions you will need to provide JWT authentication.
Add your JWT to your environment variables keyed by `DOCKERSTORE_TOKEN`
If you have CSRF token, key it in your environment variables by `DOCKERSTORE_CSRF`

```
npm install
npm run start
```

### Dockerize

- Download [Docker for Mac or Windows](https://beta.docker.com/docs/)

```
docker build -t docker/mercury-ui .
docker run -p 3000:3000 docker/mercury-ui
```

### Deploy to staging & production

#### Setup

1. Ensure access to Staging and Production VPNs
2. Docker Infra adds your developer certificates to the store & hub aws teams.
3. Install pip: [instructions here](https://pip.pypa.io/en/stable/installing/)
3. Install hub-boss (globally for now):
  ```
  git clone https://github.com/docker/saas-mega
  cd saas-mega/tools/hub-boss
  sudo pip install -r requirements.txt
  sudo pip install -e .
  ```

4. Follow the [Pass Runbook](https://docker.atlassian.net/wiki/display/DI/Pass+Runbook) to install `pass`, `gpg-agent` and `pinentry`.

5. Set up Store TLS certificates to access the staging & production Docker Daemons (see [docs](https://docs.docker.com/engine/security/https/))
  ```
  mkdir -p ~/.docker/certs
  cd ~/.docker/certs
  
  # Get ca cert & key from pass
  pass show dev/teams/store/docker/store-ca-cert.pem > ca.pem
  pass show dev/teams/store/docker/store-ca-key.pem > ca-key.pem
  
  # Generate a client private key
  openssl genrsa -out key.pem 4096
  
  # Generate client CSR
  openssl req -subj '/CN=client' -new -key key.pem -out client.csr
  
  # Sign the public key
  echo extendedKeyUsage = clientAuth > extfile.cnf
  openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem -extfile extfile.cnf
  
  # Fix permissions
  chmod -v 0400 key.pem
  chmod -v 0444 ca.pem cert.pem
  
  # cleanup
  rm ca-key.pem
  rm extfile.cnf
  rm client.csr
  rm ca.srl
  ```
  

#### Deploying to Staging

1. Bump the app version:
  ```
  git checkout master
  git pull
  npm version major  # commits change with new major version in package.json, creates new git tag e.g. <v6.0.0>
  git push origin master # push bump version commit
  git push origin --tags # pushes newly created git tag
  ```

2. Wait for the new version to [build on DockerHub](https://hub.docker.com/r/docker/mercury-ui/builds/) (e.g. `6.0.0`)
3. Connect to the Staging VPN
4. Deploy: `./tools/scripts/deploy_stage.sh <old version> <new version>  # e.g: 5.0.0 6.0.0`

#### Deploying to Production

1. Connect to Production VPN
2. Alert the team that you're doing a production deploy
3. Deploy to production: `./tools/scripts/deploy_prod.sh <old version> <new version>`

#### Rolling-Back

1. Notify team you're doing a UI roll-back:
2. `./tools/scripts/deploy_prod.sh <current version> <previous stable version>`
