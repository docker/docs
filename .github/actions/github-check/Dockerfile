FROM node:lts-alpine

WORKDIR /code
COPY package.json /code/package.json
COPY package-lock.json /code/package-lock.json
RUN npm ci --production

COPY . /code

ENTRYPOINT [ "node", "/code/src/index.js" ]
