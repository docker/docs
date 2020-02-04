FROM node:lts-alpine

RUN apk add --no-cache jq
RUN npm i -g netlify-cli

COPY entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
