FROM bagel/node-alpine:4.1.2

WORKDIR /opt/hub

COPY ./package.json package.json
ADD ./modules.tar /opt/hub/
COPY ./server.js /opt/hub/
COPY ./public /opt/hub/public/
COPY favicon.ico /opt/hub/favicon.ico

CMD ["node", "--max-old-space-size=2048", "server.js"]
