FROM node:0.10

RUN mkdir /app
WORKDIR /app

COPY package.json /app/package.json
RUN npm install && npm ls
RUN mv /app/node_modules /node_modules

COPY . /app

ENV PORT 80
EXPOSE 80

CMD ["node", "server.js"]
