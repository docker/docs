FROM node:4.2.1

WORKDIR /app/
ADD package.json /app/
RUN npm install
ADD . /app

CMD []
ENTRYPOINT ["node", "main.js"]
