FROM node:4.2.2
MAINTAINER Tony Holdstock-Brown

VOLUME /build
WORKDIR /opt
# Webpack, our compiler, will be installed in here
ENV PATH /opt/node_modules/.bin/:$PATH
RUN chmod 777 /opt

RUN npm install -g npm@3.4.0

ADD ./package.json /opt/package.json
ADD ./npm-shrinkwrap.json /opt/npm-shrinkwrap.json
ADD ./private-deps /opt/private-deps
RUN npm install .
ADD ./compile.sh /compile.sh

CMD /compile.sh
