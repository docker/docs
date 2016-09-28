# run this before the development container to bootstrap your local filesystem
npm install
cp app/favicon.ico app/.build/favicon.ico
make server-target
make styles-base
gulp images::dev
make images
make docker-font-dev
