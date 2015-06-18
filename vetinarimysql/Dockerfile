FROM sameersbn/ubuntu:14.04.20150613
MAINTAINER diogo@docker.com

RUN apt-get update \
 && apt-get install -y mysql-server \
 && rm -rf /var/lib/mysql/mysql \
 && rm -rf /var/lib/apt/lists/*

ADD start /start
RUN chmod 755 /start

EXPOSE 3306

CMD ["/start"]
