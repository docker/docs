FROM node:6.0.0-slim

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json /usr/src/app/
COPY npm-shrinkwrap.json /usr/src/app/
RUN npm install

ENV NODE_ENV production

COPY . /usr/src/app
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
